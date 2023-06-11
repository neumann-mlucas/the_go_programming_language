package main

import (
	"fmt"
	"golang.org/x/net/html"
	"net/http"
	"os"
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

	tag := elementByID(doc, "quote_slide0")
	fmt.Println(tag)
}

func elementByID(doc *html.Node, id string) *html.Node {
	if checkID(doc, id) {
		return doc
	}
	for c := doc.FirstChild; c != nil; c = c.NextSibling {
		node := elementByID(c, id)
		if node != nil {
			return node
		}
	}
	return nil
}

func checkID(n *html.Node, targetId string) bool {
	for _, item := range n.Attr {
		if item.Key == "id" && item.Val == targetId {
			return true
		}
	}
	return false
}
