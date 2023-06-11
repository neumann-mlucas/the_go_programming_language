package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	// "strings"
	// "time"
)

const (
	apikey  = "e96228fc"
	baseurl = "http://www.omdbapi.com/"
)

func invalidArgument() {
	var usage = "%s\n usage:\t[MOVIE NAME]\n"
	log.Fatalf(usage, os.Args)
}

type MovieInfo struct {
	Title   string
	Year    string
	Runtime string
	Genre   string
	Poster  string
}

func main() {
	if len(os.Args) != 2 {
		invalidArgument()
	}
	title := os.Args[1]

	URL := buildURL(title)
	movieInfo, err := getInfo(URL)
	if err != nil {
		log.Fatal(err)
	}
	getPoster(movieInfo)

}

func buildURL(movieTitle string) string {
	query := url.Values{}
	query.Set("t", url.QueryEscape(movieTitle))
	query.Set("apikey", apikey)

	movieURL, _ := url.Parse(baseurl)
	movieURL.RawQuery = query.Encode()
	return movieURL.String()
}

func getInfo(url string) (*MovieInfo, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("can't get %s: %s", url, resp.Status)
	}

	var info MovieInfo
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		return nil, err
	}
	return &info, nil
}

func getPoster(info *MovieInfo) {
	resp, err := http.Get(info.Poster)
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		log.Fatalf("can't get %s: %s", info.Poster, resp.Status)
	}

	writer := bufio.NewWriter(os.Stdout)
	_, err = writer.ReadFrom(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
}
