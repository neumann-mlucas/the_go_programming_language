package main

import (
	"fmt"
	"time"
)

var ping = make(chan int)
var pong = make(chan int)

func main() {
	go func() {
		for {
			value := <-ping
			pong <- (value + 1)
		}

	}()

	go func() {
		for {
			value := <-pong
			ping <- (value + 1)
		}
	}()

	ping <- 0
	time.Sleep(1 * time.Second)
	fmt.Println("Pings: ", <-ping)
}
