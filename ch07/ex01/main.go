package main

import (
	"bufio"
	"bytes"
	"fmt"
)

func main() {
	var c WordCounter
	c.Write([]byte("hello world"))
	fmt.Println(c) // "2"

	c = 0 // reset the counter
	var name = "one\ntwo\nthree"
	fmt.Fprintf(&c, "%s", name)
	fmt.Println(c) // "3"
}

type WordCounter int

func (c *WordCounter) Write(p []byte) (int, error) {
	scanner := bufio.NewScanner(bytes.NewReader(p))
	scanner.Split(bufio.ScanWords)
	count := 0
	for scanner.Scan() {
		count++
	}

	*c += WordCounter(count)
	return count, nil
}

type LineCounter int

func (c *LineCounter) Write(p []byte) (int, error) {
	scanner := bufio.NewScanner(bytes.NewReader(p))
	scanner.Split(bufio.ScanLines)
	count := 0
	for scanner.Scan() {
		count++
	}

	*c += LineCounter(count)
	return count, nil
}
