package main

import "fmt"

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

func add(done <-chan interface{}, ch <-chan int, adder int) <-chan int {
	res := make(chan int)
	go func() {
		defer close(res)
		for i := range ch {
			select {
			case <-done:
				return
			case res <- i + adder:
			}
		}
	}()

	return res
}

func multiply(done <-chan interface{}, ch <-chan int, multiplier int) <-chan int {
	res := make(chan int)
	go func() {
		defer close(res)
		for i := range ch {
			select {
			case <-done:
				return
			case res <- i * multiplier:
			}
		}
	}()

	return res
}

func main() {
	ints := []int{1, 2, 3, 4, 5}

	done := make(chan interface{})
	defer close(done)

	ch := generator(done, ints...)

	for i := range add(done, multiply(done, ch, 5), 10) {
		fmt.Println(i)
	}
}
