package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"unicode"
	// "unicode/utf8"
)

func main() {
	counts := make(map[rune]int)       // count of Unicode characters
	categories := make(map[string]int) // count of Unicode characters types
	utflen := make(map[int]int)        // count of lengths of UTF-8 encodings
	invalid := 0                       // count of invalid UTF-8 characters

	in := bufio.NewReader(os.Stdin)
	for {
		r, n, err := in.ReadRune() // returns rune, nbytes, error
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Fprintf(os.Stderr, "charcount: %v\n", err)
			os.Exit(1)
		}

		if r == unicode.ReplacementChar && n == 1 {
			invalid++
			continue
		}
		countType(r, categories)
		counts[r]++
		utflen[n]++
	}

	fmt.Printf("rune\tcount\n")
	for c, n := range counts {
		fmt.Printf("%q\t%d\n", c, n)
	}

	fmt.Printf("\ntype\tcount\n")
	for c, n := range categories {
		fmt.Printf("%s\t%d\n", c, n)
	}

	fmt.Print("\nlen\tcount\n")
	for i, n := range utflen {
		if i > 0 {
			fmt.Printf("%d\t%d\n", i, n)
		}
	}

	if invalid > 0 {
		fmt.Printf("\n%d invalid UTF-8 characters\n", invalid)
	}
}

func countType(r rune, categories map[string]int) {
	if unicode.IsControl(r) {
		categories["Control"]++
	}
	if unicode.IsDigit(r) {
		categories["Digit"]++
	}
	if unicode.IsGraphic(r) {
		categories["Graphic"]++
	}
	if unicode.IsLetter(r) {
		categories["Letter"]++
	}
	if unicode.IsLower(r) {
		categories["Lower"]++
	}
	if unicode.IsMark(r) {
		categories["Mark"]++
	}
	if unicode.IsNumber(r) {
		categories["Number"]++
	}
	if unicode.IsPrint(r) {
		categories["Print"]++
	}
	if unicode.IsPunct(r) {
		categories["Punct"]++
	}
	if unicode.IsSpace(r) {
		categories["Space"]++
	}
	if unicode.IsSymbol(r) {
		categories["Symbol"]++
	}
	if unicode.IsTitle(r) {
		categories["Title"]++
	}
	if unicode.IsUpper(r) {
		categories["Upper"]++
	}
}
