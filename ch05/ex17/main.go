package main

import (
	"fmt"
	"golang.org/x/net/html"
	"net/http"
	"os"
)

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

	tags := elementsByTagName(doc, "h1", "h2", "h3", "h4")
	for _, tag := range tags {
		fmt.Printf("<%s> %s\n", tag.Data, tag.FirstChild.Data)
	}
}

var nodes []*html.Node

func elementsByTagName(doc *html.Node, names ...string) []*html.Node {
	if checkTagName(doc, names...) {
		nodes = append(nodes, doc)
	}

	for c := doc.FirstChild; c != nil; c = c.NextSibling {
		elementsByTagName(c, names...)
	}
	return nodes
}

func checkTagName(n *html.Node, names ...string) bool {
	for _, name := range names {
		if n.Type == html.ElementNode && n.Data == name {
			return true
		}
	}
	return false
}
