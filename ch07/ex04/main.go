package main

import (
	"fmt"
	"io"
)

func main() {
	s := NewReader("hello world")
	buf := make([]byte, 5)

	fmt.Println(s.Read(buf))
	fmt.Println(string(buf))
	fmt.Println(s.Read(buf))
	fmt.Println(string(buf))

}

func NewReader(s string) *StrReader {
	return &StrReader{s, 0}
}

type StrReader struct {
	s string
	p int
}

func (sr *StrReader) Read(p []byte) (int, error) {
	n := copy(p, sr.s)
	sr.s = sr.s[n:]
	if len(sr.s) == 0 {
		return 0, io.EOF
	}
	return n, nil
}
