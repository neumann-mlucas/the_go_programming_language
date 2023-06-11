package main

import (
	"fmt"
	"unicode"
	"unicode/utf8"
)

func main() {
	s := "Hello　世界"
	fmt.Println(s)
	u := removeUnicodeSpace([]byte(s))
	fmt.Println(string(u))
}

func removeUnicodeSpace(s []byte) []byte {
	for i := 0; i < len(s); {
		r, size := utf8.DecodeRune(s[i:])
		if unicode.IsSpace(r) {
			s[i] = ' '
			copy(s[i+1:], s[i+size:])
			s = s[:len(s)-size+1]
			i++
		} else {
			i += size
		}
	}
	return s
}
