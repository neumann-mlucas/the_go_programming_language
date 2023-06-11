package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	// count of word frequencies
	words := make(map[string]int)
	input := bufio.NewScanner(os.Stdin)
	input.Split(bufio.ScanWords)
	for input.Scan() {
		words[input.Text()]++
	}
	// sort frequencies
	revWords := make(map[int][]string)
	var nmax int
	for c, n := range words {
		revWords[n] = append(revWords[n], c)
		if n > nmax {
			nmax = n
		}

	}
	// print frequencies
	fmt.Printf("word\tcount\n")
	for i := nmax; i > 0; i-- {
		for _, word := range revWords[i] {
			fmt.Printf("%s\t%d\n", word, i)
		}
	}
}
