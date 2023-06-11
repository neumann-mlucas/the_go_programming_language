package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	files := os.Args[1:]
	if len(files) == 0 {
		counts := countLines(os.Stdin)
		printCounts("Stdin", counts)
	} else {
		for _, arg := range files {
			f, err := os.Open(arg)
			if err != nil {
				fmt.Fprintf(os.Stderr, "dub2: %v\n", err)
				continue
			}
			counts := countLines(f)
			printCounts(arg, counts)
			f.Close()
		}
	}
}

func countLines(f *os.File) map[string]int {
	counts := make(map[string]int)
	input := bufio.NewScanner(f)
	for input.Scan() {
		counts[input.Text()]++
	}
	return counts
}

func printCounts(file string, counts map[string]int) {
	for line, n := range counts {
		if n > 1 {
			fmt.Printf("%s\t%d\t%s\n", file, n, line)
		}
	}
}
