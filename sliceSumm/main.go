package main

import (
	"sliceSumm/logger"
	"sync"
)

func main() {
	arr := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}
	numG := 3
	summ := 0
	summChan := make(chan int, numG)
	var wg sync.WaitGroup
	var countInG int = len(arr) / numG

	for i := 0; i < numG; i++ {
		wg.Add(1)
		go func() {
			summSlice(arr[i*countInG:i*countInG+countInG], summChan)
			wg.Done()
		}()
	}

	go func() {
		wg.Wait()
		close(summChan)
	}()

	for message := range summChan {
		summ += message
	}

	logger.Success(summ)
}

func summSlice(nums []int, ch chan int) {
	summ := 0
	for _, num := range nums {
		summ += num
	}

	ch <- summ
}
