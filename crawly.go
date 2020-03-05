package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
)

// Diccionario donde almacenamos las urls visitadas
var visited = make(map[string]bool)
var not_vis = make([]string)

func main() {
	// Argumentos iniciales
	flag.Parse()
	args := flag.Args()
	fmt.Println(args)
	if len(args) < 1 {
		fmt.Printf("Por favor, especifique la url\n")
		fmt.Printf("Uso: ./crawly http://google.com\n")
		os.Exit(1)
	}

	/* Creamos el canal de comunicación con las gorutinas
	   que va a ser una cola donde especificamos las urls */
	cola := make(chan string)

	// Introducimos el elemento en la cola
	go func() { cola <- args[0] }()

	/* Recorremos la cola para ver sus elementos
	   y los incluimos para que sean encolados */
	for uri := range cola {
		fetch(uri, cola)
	}

}

func fetch(u string, cola chan string) {
	fmt.Println("Fetching: ", u)
	// Indicamos que la URL ha sido visitada
	visited[u] = true
	// Deshabilitamos la validación SSL
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
	}
	transport := &http.Transport{
		TLSClientConfig: tlsConfig,
	}
	// Creamos el client HTTP
	client := &http.Client{Transport: transport}
	// Creamos la petición GET
	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		fmt.Sprint(err)
	}

	// Lanzamos la peticion
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Se produjo un error leyendo la respuesta\n")
		not_vis = append(not_vis, u)
		return
	}

	// Con defer prorrogamos el cierre de la conexion hasta que la funcion acaba
	defer resp.Body.Close()

	// Leemos la respuesta de la peticion
	body, _ := ioutil.ReadAll(resp.Body)
	//fmt.Printf(string(body))

	// Parseamos los resultados para volver a revisarlos a incluirlos en el listado
	links := parse(body)
	//Recorremos el listado de resultados y encolamos las nuevas urls
	for _, l := range links {
		abs := fixUrl(l, u)
		if u != "" {
			// Si no ha sido visitada
			if !visited[abs] {
				// Metemos la uri en la cola
				go func() { cola <- abs }()
			}
		}
	}
}
func parse(body []byte) (res_u []string) {
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
		res_u = append(res_u, res)
	}
	return
}

// Metodo para resolver las ulrs relativas que se encuentran
func fixUrl(href, base string) string {
	uri, err := url.Parse(href)
	if err != nil {
		// Si no se consigue parsear se devuelve vacío
		return ""
	}
	baseUrl, err := url.Parse(base)
	if err != nil {
		return ""
	}
	uri = baseUrl.ResolveReference(uri)
	return uri.String()
}
