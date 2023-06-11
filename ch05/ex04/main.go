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
	for _, link := range visit(doc) {
		fmt.Println(link)
	}
}

func visit(doc *html.Node) []string {
	var links []string
	var fn func(*html.Node)
	fn = func(n *html.Node) {
		if n == nil {
			return
		}
		if n.Type == html.ElementNode && (n.Data == "img" || n.Data == "script") {
			for _, a := range n.Attr {
				if a.Key == "src" {
					links = append(links, a.Val)
				}
			}
		}
		fn(n.FirstChild)
		fn(n.NextSibling)
	}
	fn(doc)
	return links
}
