package main

import (
	"bufio"
	"fmt"
	"golang.org/x/net/html"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func main() {
	var urls []*url.URL
	for _, rawURL := range os.Args[1:] {
		URL, err := url.Parse(rawURL)
		if err == nil {
			urls = append(urls, URL)
		}
	}
	breadthFirst(crawl, urls)
}

func crawl(URL *url.URL) []*url.URL {
	fmt.Println(URL)
	_, err := saveHTML(URL)
	if err != nil {
		log.Print(err)
		return nil
	}

	list, err := Extract(URL)
	if err != nil {
		log.Print(err)
		return nil
	}
	return list
}

func breadthFirst(f func(item *url.URL) []*url.URL, worklist []*url.URL) {
	seen := make(map[string]bool)
	for len(worklist) > 0 {
		items := worklist
		worklist = nil
		for _, item := range items {
			if !seen[item.String()] {
				seen[item.String()] = true
				worklist = append(worklist, f(item)...)
			}
		}
	}
}

func Extract(URL *url.URL) ([]*url.URL, error) {
	resp, err := http.Get(URL.String())
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("getting %s: %s", URL, resp.Status)
	}

	doc, err := html.Parse(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("parsing %s as HTML: %v", URL, err)
	}

	var links []*url.URL
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
				if link.Host == URL.Host {
					links = append(links, link)
				}
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

func saveHTML(URL *url.URL) (int64, error) {
	file, err := createFile(URL)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	resp, err := http.Get(URL.String())
	if err != nil {
		return 0, err
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return 0, fmt.Errorf("getting %s: %s", URL, resp.Status)
	}

	writer := bufio.NewWriter(file)
	n, err := writer.ReadFrom(resp.Body)
	if err != nil {
		return 0, err
	}
	return n, nil
}

func createFile(URL *url.URL) (*os.File, error) {
	filename := "html/" + URL.Host + URL.Path + ".html"
	parentDir := filename[:strings.LastIndex(filename, "/")]
	os.MkdirAll(parentDir, 0777)
	return os.Create(filename)
}
