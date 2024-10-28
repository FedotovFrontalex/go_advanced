package main

import (
	"flag"
	"net/http"
	"os"
	"strings"
	"urlPingChecker/logger"
)

func main() {
	path := flag.String("file", "url.txt", "Path to urls file")
	flag.Parse()

	file, err := os.ReadFile(*path)

	if err != nil {
		panic(err.Error())
	}

	urlSlice := strings.Split(string(file), "\n")

	ch := make(chan int)
	errCh := make(chan error)

logger.Message(urlSlice)

	for _, url := range urlSlice {
		if url == "" {
				continue
		}
		go ping(url, ch, errCh)
	}

	for range urlSlice {
		select {
		case err := <-errCh:
			logger.Error(err)
		case status := <-ch:
			logger.Success(status)
		}
	}
}

func ping(url string, ch chan int, errCh chan error) {
	resp, err := http.Get(url)

	if err != nil {
		errCh <- err
		return
	}

	ch <- resp.StatusCode
}
