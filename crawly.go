package main

/* TO-DO
- Exportar los resultados navegados a un fichero
- Exportar los resultados erroneos a un fichero
- Herramienta para extraer solo las etiquetas href
- Incluir opcion para cabeceras por si se necesita un "Authentication Basic"
*/
import (
	"crypto/tls"
	"flag"
	"fmt"
	//	"io/ioutil"
	"github.com/jackdanger/collectlinks"
	"net/http"
	nu "net/url"
	"os"
	"regexp"
	"strings"
	"time"
)

// Flags globales
var url = flag.String("u", "", "URL to crawl")

// Diccionario donde almacenamos las urls visitadas
var visited = make(map[string]bool)

// Diccionario donde almacenamos las urls erroneas
var not_visited = make([]string, 0)

func main() {
	// Banner
	banner()

	// Indicamos que con -h se muestre nuestro metodo de ayuda
	flag.Usage = usage
	// Argumentos iniciales
	flag.Parse()
	if *url == "" {
		//flag.PrintDefaults()
		usage()
	}
	timestamp()

	// Diccionario de cookies
	var cookie_jar = flag.Args()
	// Comprobamos si hay cookies
	if len(cookie_jar) != 0 {
		fmt.Println("Cookies:")
		for i := range cookie_jar {
			fmt.Printf("- %s\n", cookie_jar[i])
		}
	}

	/* Creamos el canal de comunicación con las gorutinas
	   que va a ser una cola donde especificamos las urls */
	queue := make(chan string)
	queue_fil := make(chan string)

	// Introducimos el elemento en la cola
	go func() { queue <- *url }()
	// Revisamos los elementos de la cola
	go filter(queue, queue_fil)

	// Canal bool para sincronizar la ejecución de N crawlers concurrentes
	done := make(chan bool)

	// Sacamos los elementos a revisar de la cola filtrada y los metemos en la cola
	for i := 0; i < 5; i++ {
		go func() {
			/* Recorremos la cola para ver sus elementos
			   y los incluimos para que sean encolados */
			for uri := range queue_fil {
				fetch(uri, queue, cookie_jar)
			}
			done <- true
		}()
	}
	<-done
}

func timestamp() {
	fmt.Printf("Date: %s", time.Now().Format("02.01.2006 15:04:05\n"))
}

func banner() {
	fmt.Printf("#####################################\n")
	fmt.Printf("######### CRAWLY  -- v0.0.1 #########\n")
	fmt.Printf("#####################################\n")
}

func usage() {
	fmt.Printf("\nERROR - Must complete all input params\n")
	fmt.Printf("\nUsage mode:\n")
	fmt.Printf("%s -u <URL> <Cookie1=Value1> <Cookie2=Value2> ... \n", os.Args[0])
	fmt.Println("Info: Cookie must be set in 'Name=Value' format")
	os.Exit(1)
}

func filter(in chan string, out chan string) {
	for v := range in {
		if !visited[v] {
			visited[v] = true
			out <- v
		}
	}
}

func fetch(u string, queue chan string, cookies []string) {
	// Deshabilitamos la validación SSL
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}

	// Creamos el client HTTP
	client := &http.Client{Transport: transport, Timeout: time.Second * 10}
	// Creamos la petición GET
	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		fmt.Sprint(err)
	}

	// Comprobamos si hay cookies
	if len(cookies) != 0 {
		for i := range cookies {
			s := strings.Split(cookies[i], "=")
			c_i := http.Cookie{Name: s[0], Value: s[1]}
			req.AddCookie(&c_i)
		}
	}

	// Lanzamos la peticion
	resp, err := client.Do(req)
	fmt.Println("Fetching: ", u)

	if err != nil {
		fmt.Printf("There was an error reading the answer\n")
		not_visited = append(not_visited, u)
		//fmt.Printf("%v\n", not_visited)
		return
	}

	// Con defer prorrogamos el cierre de la conexion hasta que la funcion acaba
	defer resp.Body.Close()

	// Leemos la respuesta de la peticion
	// body, _ := ioutil.ReadAll(resp.Body)
	//fmt.Printf(string(body))

	// Parseamos los resultados para volver a revisarlos a incluirlos en el listado
	//links := parse(body)
	links := collectlinks.All(resp.Body)
	//Recorremos el listado de resultados y encolamos las nuevas urls
	for _, l := range links {
		abs := fixUrl(l, u)
		if u != "" {
			// Metemos la uri en la cola
			go func() { queue <- abs }()
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
	uri, err := nu.Parse(href)
	if err != nil {
		// Si no se consigue parsear se devuelve vacío
		return ""
	}
	baseUrl, err := nu.Parse(base)
	if err != nil {
		return ""
	}
	uri = baseUrl.ResolveReference(uri)
	return uri.String()
}
