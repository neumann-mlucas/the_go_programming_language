package main

import (
	"fmt"
	"time"
)

func main() {
	in := make(chan int)
	in_ := in

	for i := 0; i < 1_000_000; i++ {
		in = create_pipe(in)
	}

	t0 := time.Now()
	in_ <- 0
	fmt.Println(<-in)
	fmt.Println(time.Since(t0))
}

func create_pipe(in chan int) chan int {
	out := make(chan int)
	go func() {
		value := <-in
		out <- value + 1
	}()
	return out
}
