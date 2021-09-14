package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// URL status
type URLVisited struct {
	sync.RWMutex
	urls map[string]bool
}

var (
	// Global flags
	initURL    = flag.String("u", "", "URL to crawl (mandatory)")
	authHeader = flag.String("a", "", "Authorization Basic Header (optional)")
	threads    = flag.Int("t", 10, "Number of threads (optional)")
	logName    = flag.String("l", "", "Log file name (optional)")
	depth      = flag.Int("d", 1, "Crawling depth (0 -> only input url, 1 -> infinite)")

	// Crawling results
	fromCrawl = make(chan string)
	// To crawl
	toCrawl = make(chan string)
	// URL status
	urlVisited URLVisited = URLVisited{
		urls: map[string]bool{},
	}
	// Concurrent wait
	wg sync.WaitGroup
)

func main() {
	// Set custom usage method
	flag.Usage = usage
	// Parse args
	flag.Parse()
	if *initURL == "" {
		fmt.Printf("ERROR - Must input <URL> \n")
		usage()
		os.Exit(1)
	}
	timestamp()

	var fileName string
	if *logName != "" {
		fileName = *logName
	} else {
		fileName = "crawl.log"
	}
	logName, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	multiOut := io.MultiWriter(logName, os.Stdout)
	if err != nil {
		log.Error().Err(err).Msg("")
	}
	defer func() {
		err := logName.Close()
		if err != nil {
			log.Error().Err(err).Msg("")
		}
	}()

	// New logger
	logger := zerolog.New(multiOut).With().Timestamp().Logger()

	// Set workers
	setWorkers(threads)

	// Insert the first URL to crawl
	go addToChannel(fromCrawl, *initURL)

	// Checks fromCrawl results to enqueue
	go checkVisited(fromCrawl, toCrawl, &urlVisited)

	// Launch workers and crawl URLs
	for i := 0; i <= *threads; i++ {
		wg.Add(1)
		go crawler(i, &wg, &logger)
	}
	// Wait until all workers finish
	wg.Wait()
}
