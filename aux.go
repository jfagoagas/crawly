package main

import (
	"flag"
	"fmt"
	nu "net/url"
	"regexp"
	"runtime"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
)

func timestamp() {
	fmt.Printf("Date: %s\n", time.Now().Format("02.01.2006 15:04:05\n"))
}

func usage() {
	fmt.Printf(`Usage mode:
-a string
	Authorization Basic Header (optional)
-d int
	Crawling depth (0 -> only input url, 1 -> infinite) (default 1)
-l string
	Log file name (optional)
-t int
	Number of threads (optional) (default 10)
-u string
	URL to crawl (mandatory)

You can set cookies as the last argument like Cookie1=Value1 Cookie2=Value2
`)
}

// Set number of concurrent threads
func setWorkers(threadNumber *int) {
	if *threads > 0 {
		runtime.GOMAXPROCS(*threads)
	}
}

// Extract hrefs
func parse(body []byte) (result []string) {
	re := regexp.MustCompile(`href="http[^ ]*"`)
	match := re.FindAll(body, -1)
	for _, value := range match {
		str := string(value)
		res := strings.Replace(str, "href=\"", "", -1)
		res = strings.Replace(res, "\"", "", -1)
		result = append(result, res)
	}
	return
}

// Resolve relative URIs
func fixURL(href, base string) string {
	uri, err := nu.Parse(href)
	if err != nil {
		log.Error().Err(err).Msg("")
		return ""
	}
	baseURL, err := nu.Parse(base)
	if err != nil {
		log.Error().Err(err).Msg("")
		return ""
	}
	uri = baseURL.ResolveReference(uri)
	return uri.String()
}

func checkCookies() {
	var cookieJar = flag.Args()
	if len(cookieJar) != 0 {
		fmt.Println("Cookies:")
		for i := range cookieJar {
			fmt.Printf("- %s\n", cookieJar[i])
		}
	}
}
