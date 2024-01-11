package main

import (
	"fmt"
	"io"
	"strings"
)

func main() {
	lr := LimitReader(strings.NewReader("hello world"), 3)

	buf := make([]byte, 1000)
	n, err := lr.Read(buf)

	fmt.Println(string(buf), n, err)
}

func LimitReader(r io.Reader, n int64) io.Reader {
	return &LimitReaderWrapper{r, 0, int(n)}
}

type LimitReaderWrapper struct {
	reader   io.Reader
	n, limit int
}

func (lr *LimitReaderWrapper) Read(p []byte) (int, error) {
	n, err := lr.reader.Read(p[:lr.limit])
	lr.n += n
	if lr.n > lr.limit {
		err = io.EOF
	}
	return n, err

}
