package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	//	"net/url"
	"os"
	"regexp"
	//  "runtime"
	"log"
	"strings"
)

// Tipo de datos para nuestras urls
type host_data struct {
	dir    string
	status bool
}

// Diccionario donde almacenamos las urls visitadas
var dicc []host_data
var results []string
var lnks map[string]bool

func main() {
	// Almacenamos la direccion inicial pasada como parametro
	elem := os.Args[1]
	// Inicializamos el mapa para almacenar los resultados
	lnks = make(map[string]bool)
	// Lo almacenamos en nuestro fichero de resultados
	lnks[elem] = false
	// String para almacenar los resultados de la primera peticion
	fmt.Printf("CRAWLING URL --> %s\n", elem)
	//fmt.Printf("%v\n", lnks)
	// Lanzamos la peticion GET a la 1ª URL
	fetch(elem)
	// Recorremos el listado de elementos encontrados
	for k, _ := range lnks {
		// Si no ha sido visitado lanzamos la peticion
		//fmt.Printf("%d\n", len(lnks))
		if lnks[k] == false {
			// HAY QUE CONTROLAR LOS ACCESOS AL MAPA PARA IR ACTUALIZANDO LAS DIRECCIONES
			// Se lanza una peticion concurrente para cada url que se ha indexado
			go fetch(k)
			//fmt.Printf("%v\n",lnks)
		} else {
			fmt.Printf("URL: %s, Visted?: %t\n", k, v)
		}
	}
	// Listado final de resultados
	//fmt.Printf("%v\n", lnks)
}

func fetch(u string) {
	// Creamos el client HTTP
	client := &http.Client{}
	// Lanzamos la peticion y si se produce un error lo capturamos
	resp, err := client.Get(u)
	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		fmt.Sprint(err)
	}

	// Cabeceras
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Host", u)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 6.1; Win64; x64; rv:73.0) Gecko/20100101 Firefox/73.0")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")
	// Hay que comentar la siguiente linea para no obtener la respuesta comprimida (gzip)
	//req.Header.Set("Accept-Encoding", "gzip, deflate")
	req.Header.Set("Connection", "close")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	//fmt.Println(req.Header)

	// Cookies de sesion
	php_sess := http.Cookie{Name: "PHPSESSID", Value: "5v3ja6sr5v76et3r77f376s0vm"}
	wp_sec := http.Cookie{Name: "wordpress_sec_d40d69d37c2b367862b80f337e8e8778", Value: "ignacio.riveracorullon%40telefonica.com%7C1583419088%7CmeI7EmqCz4KszVV4XKsKzn1vIWMuRkDbLmAJyUdQgaC%7Cf4bf55a0831dcea3c81ee19ce0881f6c1b6eb5709824c0e3801da30389fae1e5"}
	wp_login := http.Cookie{Name: "wordpress_logged_in_d40d69d37c2b367862b80f337e8e8778", Value: "ignacio.riveracorullon%40telefonica.com%7C1583419088%7CmeI7EmqCz4KszVV4XKsKzn1vIWMuRkDbLmAJyUdQgaC%7C8f7175c753c982599c1b2672bcfc66d19382abcbf2e731cefe18a4babb852d42"}
	wp_lan := http.Cookie{Name: "wp-wpml_current_language", Value: "es"}
	req.AddCookie(&php_sess)
	req.AddCookie(&wp_sec)
	req.AddCookie(&wp_login)
	req.AddCookie(&wp_lan)
	//fmt.Println(req.Cookies())

	// Lanzamos la peticion
	resp, err = client.Do(req)
	if err != nil {
		log.Fatal("Error reading response. ", err)
	} else {
		// Indicamos que esa URL ya ha sido visitada
		lnks[u] = true
		//fmt.Printf("%v\n", lnks)
		status = true
	}
	defer resp.Body.Close()

	// Leemos la respuesta de la peticion
	body, _ := ioutil.ReadAll(resp.Body)
	//fmt.Printf(string(body))

	// Parseamos los resultados para volver a revisarlos a incluirlos en el listado
	parse(body)
}

func parse(body []byte) {
	// Expresion regular para sacar las etiquetas href
	re := regexp.MustCompile(`href="http[^ ]*"`)
	// Ejecutamos la expresion regular
	match := re.FindAll(body, -1)
	// Recorremos los resultados
	for i := 0; i < len(match); i++ {
		str := string(match[i])
		// Eliminamos el contenido que nos sobra (-1 para todas las ocurrencias)
		res := strings.Replace(str, "href=\"", "", -1)
		res = strings.Replace(res, "\"", "", -1)
		//fmt.Printf("%q\n", match[i])
		//fmt.Println(res)
		// Si el mapa contiene el elemento no lo incluimos de nuevo
		if _, ok := lnks[res]; !ok {
			lnks[res] = false
		}
	}
}
