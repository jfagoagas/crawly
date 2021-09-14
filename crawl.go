package main

import (
	"crypto/tls"
	"encoding/base64"
	"flag"
	"io/ioutil"
	"regexp"
	"sync"

	"net/http"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func fetch(url string, fromCrawl chan string, logger *zerolog.Logger) {
	// Disable SSL verification
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}

	client := &http.Client{Transport: transport, Timeout: time.Second * 10}
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Error().Err(err).Msg("")
	}

	checkAuthHeader(request)
	addCookies(request)

	logger.Info().Msg("Fetching: " + url)
	response, err := client.Do(request)
	if err != nil {
		log.Error().Err(err).Msg("There was an error reading the answer")
		return
	}
	defer func() {
		err := response.Body.Close()
		if err != nil {
			log.Error().Err(err).Msg("")
		}
	}()

	body, _ := ioutil.ReadAll(response.Body)
	links := parse(body)
	for _, link := range links {
		absolute := fixURL(link, url)
		if url != "" {
			// Insert the new URL
			go func() { fromCrawl <- absolute }()
		}
	}
}

func addToChannel(ch chan<- string, url string) {
	ch <- url
}

func checkVisited(fromCrawl chan string, toCrawl chan string, urlVisited *URLVisited) {
	r := regexp.MustCompile(`^(?:https?:\/\/)?(?:[^@\/\n]+@)?(?:www\.)?([^:\/?\n]+)`)
	for url := range fromCrawl {
		urlVisited.Lock()
		if ok := urlVisited.urls[url]; !ok {
			if *depth == 0 && strings.Contains(url, r.FindString(*initURL)) {
				go addToChannel(toCrawl, url)
			} else if *depth == 1 {
				go addToChannel(toCrawl, url)
			}
			urlVisited.urls[url] = true
		}
		urlVisited.Unlock()
	}
}

func crawler(id int, wg *sync.WaitGroup, logger *zerolog.Logger) {
	defer wg.Done()
	for uri := range toCrawl {
		fetch(uri, fromCrawl, logger)
	}

}

func checkAuthHeader(request *http.Request) {
	if *authHeader != "" {
		authEnc := base64.StdEncoding.EncodeToString([]byte(*authHeader))
		request.Header.Add("Authorization:", "Basic "+authEnc)
	}
}

func addCookies(request *http.Request) {
	cookieJar := flag.Args()
	if len(cookieJar) != 0 {
		for _, value := range cookieJar {
			cookie := http.Cookie{Name: strings.Split(value, "=")[0], Value: strings.Split(value, "=")[1]}
			request.AddCookie(&cookie)
		}
	}
}
