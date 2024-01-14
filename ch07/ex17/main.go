package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"strings"
)

func main() {
	args := parseArgs(os.Args[1:])
	dec := xml.NewDecoder(os.Stdin)

	var stack []xml.StartElement
	for {
		tok, err := dec.Token()
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Fprintf(os.Stderr, "xmlselect: %v\n", err)
			os.Exit(1)
		}
		switch tok := tok.(type) {
		case xml.StartElement:
			stack = append(stack, tok) // push
		case xml.EndElement:
			stack = stack[:len(stack)-1] // pop
		case xml.CharData:
			if containsAll(stack, args) {
				fmt.Printf("%s: %s\n", join(stack), tok)
			}
		}
	}
}

func containsAll(x []xml.StartElement, y []Token) bool {
	for len(y) <= len(x) {
		if len(y) == 0 {
			return true
		}
		if y[0].comp(x[0]) {
			y = y[1:]
		}
		x = x[1:]
	}
	return false
}

type Token struct {
	name  string
	attrs map[string]string
}

func parseArgs(args []string) []Token {
	tokens := []Token{}
	for _, arg := range args {
		fields := strings.Fields(arg)
		tokens = append(tokens, Token{name: fields[0], attrs: parseAttrs(fields)})
	}
	return tokens
}

func parseAttrs(fields []string) map[string]string {
	attrs := make(map[string]string)
	if len(fields) <= 1 {
		return attrs
	}
	for _, attr := range fields[1:] {
		if strings.Contains(attr, "=") {
			parts := strings.SplitN(attr, "=", 2)
			attrs[parts[0]] = parts[1]
		}
	}
	return attrs
}

func (t Token) comp(tok xml.StartElement) bool {
	if t.name != "." && t.name != tok.Name.Local {
		return false
	}
	for k, v := range t.attrs {
		if !contains(tok, k, v) {
			return false
		}
	}
	return true
}

func contains(t xml.StartElement, k, v string) bool {
	for _, attr := range t.Attr {
		// fmt.Println(t.Name.Local, t.Attr)
		if attr.Name.Local == k && attr.Value == v {
			return true
		}
	}
	return false
}

func join(stack []xml.StartElement) string {
	var s []string
	for _, t := range stack {
		s = append(s, t.Name.Local)
	}
	return strings.Join(s, " ")
}
