package main

import (
	"fmt"
	"sync"
	"time"
)

func generator(done <-chan interface{}, ints ...int) <-chan int {
	res := make(chan int)
	go func() {
		defer close(res)
		for _, i := range ints {
			select {
			case <-done:
				return
			case res <- i:
			}
		}
	}()
	return res
}

func fanIn(done <-chan interface{}, chans ...<-chan int) <-chan int {
	wg := sync.WaitGroup{}
	joined := make(chan int)

	join := func(ch <-chan int) {
		defer wg.Done()
		for i := range ch {
			select {
			case <-done:
				return
			case joined <- i:
			}
		}
	}

	wg.Add(len(chans))
	for _, ch := range chans {
		go join(ch)
	}

	go func() {
		wg.Wait()
		close(joined)
	}()

	return joined
}

func slowFunc(ch <-chan interface{}) <-chan interface{} {
	time.Sleep(3)
	return ch
}

func main() {
	ints := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	numChans := 5

	done := make(chan interface{})
	defer close(done)

	chans := make([]<-chan int, numChans)
	for i := 0; i < numChans; i++ {
		chans[i] = generator(done, ints[i*2:(i+1)*2]...)
	}

	for i := range fanIn(done, chans...) {
		fmt.Println(i)
	}
}
