package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

type client struct {
	name string
	ch   chan string // an outgoing message channel
}

func (c *client) formatMsg(text string) string {
	return fmt.Sprintf("[%s] :: %s", c.name, text)

}

var (
	entering = make(chan client)
	leaving  = make(chan client)
	messages = make(chan string) // all incoming client messages
)

func main() {
	listener, err := net.Listen("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}

	go broadcaster()
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		go handleConn(conn)
	}
}

func broadcaster() {
	clients := make(map[client]bool) // all connected clients
	for {
		select {
		case msg := <-messages:
			// Broadcast incoming message to all
			// clients' outgoing message channels.
			for cli := range clients {
				cli.ch <- msg
			}

		case cli := <-entering:
			clients[cli] = true

		case cli := <-leaving:
			delete(clients, cli)
			close(cli.ch)
		}
	}
}

func handleConn(conn net.Conn) {
	input := bufio.NewScanner(conn)
	var user string
	for input.Scan() {
		user = input.Text()
		break
	}
	client := client{user, make(chan string)} // outgoing client messages

	go clientWriter(conn, client)

	client.ch <- "You are " + user
	messages <- client.name + " has arrived"
	entering <- client

	for input.Scan() {
		messages <- client.formatMsg(input.Text())
	}

	leaving <- client
	messages <- client.name + " has left"
	conn.Close()
}

func clientWriter(conn net.Conn, c client) {
	for msg := range c.ch {
		fmt.Fprintln(conn, msg) // NOTE: ignoring network errors
	}
}
