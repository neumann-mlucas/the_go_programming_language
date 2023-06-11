package main

import "fmt"

func main() {
	fmt.Println(isAnagram("arara", "arara"))
}

func isAnagram(sa, sb string) bool {
	if len(sa) != len(sb) {
		return false
	}

	for i := 0; i < len(sa); i++ {
		if sa[i] != sb[len(sb)-i-1] {
			return false
		}
	}
	return true
}
