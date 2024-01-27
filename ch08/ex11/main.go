package main

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

func main() {
	cancel := make(chan struct{})
	responses := make(chan string, 3)
	for _, url := range []string{"https://gnu.org", "https://kernal.org", "https://google.com"} {
		localUrl := url
		go func() {
			resp, err := requestWithCancel(localUrl, cancel)
			if err == nil {
				responses <- resp
			}
		}()
	}
	url := <-responses // return the quickest response
	fmt.Println("Fastest: ", url)
	time.Sleep(1 * time.Second)
}

func requestWithCancel(url string, abort chan struct{}) (string, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		<-abort
		cancel()
	}()

	fmt.Println("Requesting url:", url)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	req = req.WithContext(ctx)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	defer resp.Body.Close()

	close(abort)
	return url, nil
}
