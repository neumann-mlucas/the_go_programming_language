package main

import (
	"encoding/gob"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
)

const (
	filename = "ComicIndex.gob"
)

type ComicIndex map[int]ComicInfo

type ComicInfo struct {
	Day        string
	Month      string
	Year       string
	Num        int
	Link       string
	News       string
	Safe_title string
	Transcript string
	Alt        string
	Img        string
	Title      string
}

func (c ComicInfo) String() string {
	return fmt.Sprintf("XKCD Comic %04d TITLE: %s>\nURL: %s\nTRANSCRIPT:\n%s", c.Num, c.Title, c.Img, c.Transcript)
}

func getComicInfo(url string) (*ComicInfo, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("can't get %s: %s", url, resp.Status)
	}

	var info ComicInfo
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		return nil, err
	}
	return &info, nil
}

func getLastComicNum() (int, error) {
	info, err := getComicInfo("https://xkcd.com/info.0.json")
	if err != nil {
		return 0, err
	}
	return info.Num, nil
}

func fetchComic(num int, comics chan *ComicInfo) {
	url := fmt.Sprintf("https://xkcd.com/%d/info.0.json", num)
	comic, err := getComicInfo(url)
	if err != nil {
		log.Printf("Can't get comic %d: %s", num, err)
		comics <- &ComicInfo{Num: num}
	} else {
		comics <- comic
	}
}

func fetchAllComics(comics chan *ComicInfo, lastComic int) {
	for n := 1; n <= lastComic; n++ {
		go fetchComic(n, comics)
	}
}

func createIndex(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	lastComic, err := getLastComicNum()
	if err != nil {
		return err
	}
	ch := make(chan *ComicInfo)
	fetchAllComics(ch, lastComic)

	index := make(ComicIndex)
	for n := 1; n <= lastComic; n++ {
		comic := <-ch
		index[comic.Num] = *comic
	}

	enc := gob.NewEncoder(file)
	err = enc.Encode(index)
	if err != nil {
		return err
	}
	return nil
}

func readIndex(filename string) (ComicIndex, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	dec := gob.NewDecoder(file)

	var index ComicIndex
	err = dec.Decode(&index)
	if err != nil {
		return nil, err
	}
	return index, nil
}

func invalidArgument() {
	var usage = "\n%s\nusage:\t[fetch|INTEGER]\n"
	log.Fatalf(usage, os.Args)
}

func main() {
	if len(os.Args) != 2 {
		invalidArgument()
	}
	cmd := os.Args[1]

	// fetch comics from xkcd.com
	if cmd == "fetch" {
		err := createIndex(filename)
		if err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}

	// search comics by number in the index

	ncomic, err := strconv.Atoi(cmd)
	if err != nil {
		invalidArgument()
	}
	index, err := readIndex(filename)
	if err != nil {
		log.Fatal(err)
	}

	comic, ok := index[ncomic]
	if ok {
		fmt.Println(comic)
	} else {
		log.Fatalf("comic %04d not in index", ncomic)
	}

}
