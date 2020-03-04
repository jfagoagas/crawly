package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	//	"net/url"
	"os"
    "regexp"
	//  "runtime"
    "strings"
)

// Tipo de datos para nuestras urls
type host_data struct {
    dir string
    status bool
}
// Diccionario donde almacenamos las urls visitadas
var dicc []host_data

func main() {
    // Almacenamos la direccion inicial pasada como parametro
    var elem host_data
    elem.dir = os.Args[1]
    elem.status = false
    dicc = append(dicc, elem)
    fmt.Printf("%v\n", dicc)
    // Lanzamos la peticion GET a la URL si no se ha visitado ya
    if elem.status == false {
	    body := fetch(elem.dir)
        //elem.status = true
        // Parseamos la respuesta para sacar las etiquetas href
        parse(body)
    }
    //fmt.Printf("%v\n", dicc)
}

func fetch(u string) []byte {
    // Lanzamos la peticion y si se produce un error lo capturamos
	resp, err := http.Get(u)
	if err != nil {
		fmt.Sprint(err)
	}
	defer resp.Body.Close()
    // Leemos la respuesta de la peticion
	body, err := ioutil.ReadAll(resp.Body)

    return body
}

func parse(body []byte) {
    // Expresion regular para sacar las etiquetas href
    re := regexp.MustCompile(`href="http[^ ]*"`)
    // Ejecutamos la expresion regular
    match := re.FindAll(body, -1)
    // Recorremos los resultados
    for i := 0; i < len(match); i++ {
        str := string(match[i])
        // Eliminamos el contenido que nos sobra
        //match[i] = strings.Replace(strings(match[i]), "href", "", -1)
        res := strings.Replace(str, "href=\"", "", -1)
        res = strings.Replace(res, "\"", "", -1)
        //fmt.Printf("%q\n", match[i])
        fmt.Println(res)
    }

}
