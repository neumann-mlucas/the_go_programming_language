package main

import (
	"fmt"
	"golang.org/x/net/html"
	"os"
	"strings"
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
	var texts []string
	var fn func(*html.Node)
	fn = func(n *html.Node) {
		if n == nil {
			return
		}
		if n.Type == html.TextNode && n.Parent.Data != "script" && n.Parent.Data != "style" {
			for _, line := range strings.Split(n.Data, "\n") {
				line = strings.Trim(line, "\n\t ")
				if line != "" {
					texts = append(texts, line)
				}
			}
		}
		fn(n.FirstChild)
		fn(n.NextSibling)
	}
	fn(doc)
	return texts
}
