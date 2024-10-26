package main

import (
	"fmt"
	"net/http"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	code := make(chan int)

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			getRequest(code)
			wg.Done()
		}()
	}

	go func() {
		wg.Wait()
		close(code)
	}()

	for res := range code {
		fmt.Println(res)
	}
}

func getRequest(codeCh chan int) {
	res, err := http.Get("https://google.ru")
	if err != nil {
		fmt.Println(err)
		codeCh <- 0
		return
	}

	codeCh <- res.StatusCode
}
