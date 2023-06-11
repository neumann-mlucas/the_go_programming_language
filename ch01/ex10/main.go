package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

func main() {
	fmt.Println("First Run")
	fetchall(os.Args[1:], "run1.out")
	fmt.Println()

	fmt.Println("Second Run")
	fetchall(os.Args[1:], "run2.out")
}

func fetchall(urls []string, filename string) {
	f, err := os.Create(filename)
	if err != nil {
		fmt.Printf("while creating %s: %v", filename, err)
		return
	}
	defer f.Close()

	start := time.Now()
	ch := make(chan string)
	for _, url := range urls {
		go fetch(url, ch) // start a goroutine
	}

	for range os.Args[1:] {
		fetchTime := []byte(<-ch)
		_, err := f.Write(fetchTime)
		if err != nil {
			fmt.Printf("while writing %s: %v", filename, err)
			return
		}
	}
	fmt.Printf("%.2fs elapsed\n", time.Since(start).Seconds())
}

func fetch(url string, ch chan<- string) {
	start := time.Now()
	resp, err := http.Get(url)
	if err != nil {
		ch <- fmt.Sprint(err) // send to channel ch
		return
	}
	nbytes, err := io.Copy(ioutil.Discard, resp.Body)
	resp.Body.Close() // don't leak resources
	if err != nil {
		ch <- fmt.Sprintf("while reading %s: %v", url, err)
		return
	}
	secs := time.Since(start).Seconds()
	ch <- fmt.Sprintf("%.2fs\t%7d\t%s\n", secs, nbytes, url)
}
