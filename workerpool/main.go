package main

import (
	"fmt"
	"time"
)

func worker(jobs <-chan int, res chan<- int) {
	for i := range jobs {
		time.Sleep(time.Second)
		res <- i * 5
	}
}

func main() {
	numWorkers := 5
	numJobs := 10

	ch := make(chan int)
	res := make(chan int)

	for i := 0; i < numWorkers; i++ {
		go worker(ch, res)
	}

	go func() {
		defer close(ch)
		for i := 1; i <= numJobs; i++ {
			ch <- i
		}
	}()

	for i := 0; i < numJobs; i++ {
		fmt.Println(<-res)
	}
}
