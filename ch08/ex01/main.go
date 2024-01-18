package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

func main() {
	// parse args
	var clocks []*Clock
	for _, arg := range os.Args[1:] {
		clocks = append(clocks, newClock(arg))
	}

	// make a connection for each clock
	ch := make(chan string)
	for _, c := range clocks {
		go c.Conn(ch)

	}

	// print clocks
	for s := range ch {
		fmt.Print(s)
	}
}

type Clock struct {
	TZ   string
	port string
	Time string
}

func (c *Clock) String() string {
	return fmt.Sprintf("> %s: %s", c.TZ, c.Time)

}

func newClock(inp string) *Clock {
	// format: US/Eastern=localhost:8000
	return &Clock{strings.Split(inp, "=")[0], strings.Split(inp, "=")[1], "00:00:00"}

}

func (c *Clock) Conn(ch chan string) {
	conn, err := net.Dial("tcp", c.port)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	for {
		message, err := bufio.NewReader(conn).ReadString('\n')
		c.Time = message
		if err != nil {
			log.Fatal(err)
			break
		}
		ch <- c.String()
	}
	close(ch)
}
