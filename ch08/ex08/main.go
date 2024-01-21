package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
	"time"
)

func echo(c net.Conn, shout string, delay time.Duration) {
	fmt.Fprintln(c, "\t", strings.ToUpper(shout))
	time.Sleep(delay)
	fmt.Fprintln(c, "\t", shout)
	time.Sleep(delay)
	fmt.Fprintln(c, "\t", strings.ToLower(shout))
}

func handleConn(c net.Conn) {
	defer c.Close()
	input := bufio.NewScanner(c)

	ch := make(chan string, 10)
	go func() {
		for input.Scan() {
			ch <- input.Text()
			fmt.Println(input.Text())
		}
	}()

	timeout := time.NewTimer(10 * time.Second)
	for {
		select {
		case <-timeout.C:
			return
		case t := <-ch:
			go echo(c, t, 1*time.Second)
			timeout.Reset(10 * time.Second)
		default:
			time.Sleep(100 * time.Millisecond)
		}
	}
}

func main() {
	l, err := net.Listen("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Print(err) // e.g., connection aborted
			continue
		}
		go handleConn(conn)
	}
}
