package main

import (
	"fmt"
	"golang.org/x/net/html"
	"net/http"
	"os"
	"strings"
)

var depth int

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "argument number error: %d\n", len(os.Args)-1)
		return
	}
	url := os.Args[1]

	resp, err := http.Get(url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "HTTP request error: %s\n", err)
		return
	}

	doc, err := html.Parse(resp.Body)
	resp.Body.Close()
	if err != nil {
		fmt.Fprintf(os.Stderr, "parsing HTML error: %s\n", err)
		return
	}

	forEachNode(doc, startNode, endNode)
}

func forEachNode(n *html.Node, pre, post func(n *html.Node)) {
	if pre != nil {
		pre(n)
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		forEachNode(c, pre, post)
	}
	if post != nil {
		post(n)
	}
}

func startNode(n *html.Node) {
	switch n.Type {
	case html.ElementNode:
		startElement(n)
	case html.TextNode:
		startText(n)
	}
}

func startElement(n *html.Node) {
	printDepth()
	fmt.Printf("<%s", n.Data)
	for _, a := range n.Attr {
		fmt.Printf(" %s=\"%s\"", a.Key, a.Val)
	}
	if n.FirstChild == nil {
		fmt.Print("/")
	}
	fmt.Print(">\n")
	depth++
}

func startText(n *html.Node) {
	text := strings.TrimSpace(n.Data)
	if len(text) > 0 {
		printDepth()
		fmt.Printf("\"%s\"\n", text)
	}
}

func endNode(n *html.Node) {
	switch n.Type {
	case html.ElementNode:
		endElement(n)
	}
}

func endElement(n *html.Node) {
	depth--
	if n.FirstChild == nil {
		return
	}
	printDepth()
	fmt.Printf("</%s>\n", n.Data)
}

func printDepth() {
	for i := 0; i < depth; i++ {
		if i%2 == 0 {
			fmt.Print("┊ ")
		} else {
			fmt.Print("  ")
		}

	}
}
