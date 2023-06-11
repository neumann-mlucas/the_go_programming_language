package main

import (
	"fmt"
	"golang.org/x/net/html"
	"net/http"
	"os"
	"strings"
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

	words, images := countWordsAndImages(doc)
	fmt.Printf("url:\t%s\nwords:\t%06d\nimages:\t%06d\n", url, words, images)
}

func countWordsAndImages(doc *html.Node) (words, images int) {
	var fn func(*html.Node)
	fn = func(n *html.Node) {
		if n == nil {
			return
		}
		if n.Type == html.ElementNode && n.Data == "img" {
			images++
		} else if n.Type == html.TextNode && n.Parent.Data != "script" && n.Parent.Data != "style" {
			words += len(strings.Split(n.Data, " "))
		}
		fn(n.FirstChild)
		fn(n.NextSibling)
	}
	fn(doc)
	return
}
