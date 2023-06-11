package main

import "fmt"

func main() {
	s := []string{"one", "one", "two", "two", "three", "four", "four", "four"}
	fmt.Println(s)
	s = removeDuplicates(s)
	fmt.Println(s)
}

func removeDuplicates(strings []string) []string {
	var prev string
	var i int
	for _, s := range strings {
		strings[i] = s
		if s != prev {
			i++
		}
		prev = s
	}
	return strings[:i]
}
