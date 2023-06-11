package main

import (
	"bytes"
	"fmt"
)

func main() {
	fmt.Println(comma("12345678901"))
}

func comma(s string) (string, error) {
	var out bytes.Buffer

	// round to the nearest multiple of 3
	upperBound := 3 * (len(s) / 3)
	if len(s)%3 == 0 {
		upperBound -= 3
	}

	out.WriteString(s[:len(s)-upperBound])
	for i := upperBound; i >= 3; i -= 3 {
		out.Write([]byte{','})
		out.Write([]byte(s[len(s)-i : len(s)-i+3]))
	}

	return out.ReadString('\n')
}
