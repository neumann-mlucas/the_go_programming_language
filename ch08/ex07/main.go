package main

import (
	"fmt"
	"golang.org/x/net/html"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const maxDepth = 10

func main() {
	worklist := make(chan []string)  // lists of URLs, may have duplicates
	unseenLinks := make(chan string) // de-duplicated URLs

	baseUrl := os.Args[1]

	// err := os.Mkdir(cleanUrl(os.Args[1]), 0777)
	// if err != nil {
	// 	panic(err)
	// }

	go func() {
		worklist <- []string{os.Args[1]}
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
			if strings.HasPrefix(link, baseUrl) && !seen[link] {
				seen[link] = true
				unseenLinks <- link
			}
		}
	}
}

func crawl(link string) []string {
	fmt.Println(link)
	list, err := Extract(link)
	if err != nil {
		log.Print(err)
	}
	return list
}

func save(url string, body io.Reader) error {
	// parse Path
	path := toPath(url)
	err := os.MkdirAll(filepath.Dir(path), 0777)
	fmt.Println(path, filepath.Dir(path), filepath.Ext(path))
	if err != nil {
		return err
	}

	// create a file
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	resp, err := http.Get(url)
	if err != nil {
		return err
	}

	// write the response body to the file
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func toPath(s string) string {
	s, _ = strings.CutPrefix(s, "https://")
	s, _ = strings.CutPrefix(s, "https://")
	if filepath.Ext(s) == "" || filepath.Ext(s) == ".dev" {
		s = filepath.Join(s, "index.html")
	}
	return s
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

	save(url, resp.Body)

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
