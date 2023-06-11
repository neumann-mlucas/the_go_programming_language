package main

import (
	"fmt"
	"unicode/utf8"
)

func main() {
	s := []byte("Hello　世界")
	fmt.Println(string(s), s)
	u := reverseBytes(s)
	fmt.Println(string(u), u)
}

func reverseBytes(a []byte) []byte {
	if utf8.RuneCount(a) == 1 {
		return a
	}
	_, s := utf8.DecodeRune(a)
	return append(reverseBytes(a[s:]), a[:s]...)
}
