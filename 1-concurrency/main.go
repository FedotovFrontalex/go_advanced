package main

import (
	"concurrency/logger"
	"math/rand"
	"sync"
)

func main() {
	numCount := 10
	initialCh := make(chan int, numCount)
	resCh := make(chan int, numCount)
	var wg sync.WaitGroup

	go createRandomSlice(numCount, initialCh)

	go func() {
		wg.Add(1)
		pow(numCount, initialCh, resCh)
		wg.Done()
	}()

	wg.Wait()

	for message := range resCh {
		logger.Success(message)
	}
}

func createRandomSlice(numCount int, numCh chan int) {
	numSlice := []int{}
	for i := 0; i < numCount; i++ {
		numSlice = append(numSlice, rand.Intn(101))
	}

	logger.Message(numSlice)

	var wg sync.WaitGroup

	for _, num := range numSlice {
		wg.Add(1)
		go func(val int) {
			defer wg.Done()
			numCh <- val
		}(num)
	}

	go func() {
		wg.Wait()
		close(numCh)
	}()
}

func pow(numCount int, initialCh chan int, resCh chan int) {
	var wg sync.WaitGroup

	for i := 0; i < numCount; i++ {
		num := <-initialCh
		wg.Add(1)
		go func(val int) {
			defer wg.Done()
			resCh <- val * val
		}(num)
	}

	go func() {
		wg.Wait()
		close(resCh)
	}()
}
