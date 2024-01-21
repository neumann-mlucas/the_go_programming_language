package main

import (
	"fmt"
	"golang.org/x/net/html"
	"log"
	"net/http"
	"os"
	"strings"
)

const maxDepth = 10

type Link struct {
	url   string
	depth int
}

func main() {
	worklist := make(chan []Link)  // lists of URLs, may have duplicates
	unseenLinks := make(chan Link) // de-duplicated URLs

	baseUrl := os.Args[1]
	// Add command-line arguments to worklist.
	go func() {
		link := Link{os.Args[1], 0}
		worklist <- []Link{link}
	}()

	// Create 20 crawler goroutines to fetch each unseen link.
	for i := 0; i < 2; i++ {
		go func() {
			for link := range unseenLinks {
				foundLinks := crawl(link)
				go func() { worklist <- foundLinks }()
			}
		}()
	}

	// The main goroutine de-duplicates worklist items
	// and sends the unseen ones to the crawlers.
	seen := make(map[string]bool)
	for list := range worklist {
		for _, link := range list {
			if strings.HasPrefix(link.url, baseUrl) && link.depth <= maxDepth && !seen[link.url] {
				seen[link.url] = true
				unseenLinks <- link
			}
		}
	}
}

func crawl(link Link) []Link {
	fmt.Println(link)
	list, err := Extract(link.url)
	if err != nil {
		log.Print(err)
	}

	links := []Link{}
	for _, i := range list {
		links = append(links, Link{i, link.depth + 1})
	}
	return links
}

func Extract(url string) ([]string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("getting %s: %s", url, resp.Status)
	}

	doc, err := html.Parse(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("parsing %s as HTML: %v", url, err)
	}

	var links []string
	visitNode := func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, a := range n.Attr {
				if a.Key != "href" {
					continue
				}
				link, err := resp.Request.URL.Parse(a.Val)
				if err != nil {
					continue // ignore bad URLs
				}
				links = append(links, link.String())
			}
		}
	}
	forEachNode(doc, visitNode, nil)
	return links, nil
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

//!-
