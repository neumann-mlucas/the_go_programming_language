package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"time"
)

type client struct {
	name string
	ch   chan string // an outgoing message channel
	idle *time.Timer
}

const timeout = 30 * time.Second

func NewClient(name string) client {
	return client{
		name,
		make(chan string),
		time.NewTimer(timeout),
	}
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
			for other := range clients {
				cli.ch <- other.name + " in the chat"
			}
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
	client := NewClient(user)

	go clientWriter(conn, client)

	client.ch <- "You are " + user
	messages <- client.name + " has arrived"
	entering <- client

	for loop := true; input.Scan() && loop; {
		select {
		case <-client.idle.C:
			messages <- client.name + " has timeout"
			loop = false
			break
		default:
			messages <- client.formatMsg(input.Text())
			client.idle.Reset(30 * time.Second)
		}
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
