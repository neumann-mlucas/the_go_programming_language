package main

import (
	"fmt"
	"golang.org/x/net/html"
	"os"
)

// example usage: curl -s https://go.dev | go run main.go
func main() {
	doc, err := html.Parse(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "findlinks: %v\n", err)
		os.Exit(1)
	}
	for tag, number := range visit(doc) {
		fmt.Printf("%-10s: %04d\n", tag, number)
	}
}

func visit(doc *html.Node) map[string]int {
	counts := make(map[string]int)
	var fn func(*html.Node)
	fn = func(n *html.Node) {
		if n == nil {
			return
		}
		if n.Type == html.ElementNode {
			counts[n.Data]++
		}
		fn(n.FirstChild)
		fn(n.NextSibling)
	}
	fn(doc)
	return counts
}
