package main

import (
	"bytes"
	"fmt"
	"strings"
)

func main() {
	fmt.Println(comma("1234567890.123"))
}

func comma(s string) (string, error) {
	var out bytes.Buffer

	lenWholeNumber := strings.Index(s, ".")
	if lenWholeNumber == -1 {
		lenWholeNumber = len(s)
	}
	commaUpperBound := 3 * (lenWholeNumber / 3)
	if lenWholeNumber%3 == 0 {
		commaUpperBound -= 3
	}

	out.WriteString(s[:lenWholeNumber-commaUpperBound])
	for i := commaUpperBound; i >= 3; i -= 3 {
		out.Write([]byte{','})
		out.Write([]byte(s[lenWholeNumber-i : lenWholeNumber-i+3]))
	}

	out.WriteString(s[lenWholeNumber:])
	return out.ReadString('\n')
}
