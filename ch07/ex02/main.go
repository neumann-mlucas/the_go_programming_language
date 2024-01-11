package main

import (
	"bytes"
	"fmt"
	"io"
)

func main() {
	var buf bytes.Buffer
	wrapper, n := CountingWriter(&buf)
	for i := 0; i < 10; i++ {
		fmt.Fprint(wrapper, "hello")
		fmt.Println(*n)
	}
}

func CountingWriter(w io.Writer) (io.Writer, *int64) {
	n := int64(0)
	wrapper := CounterWrapper{w, &n}
	return &wrapper, &n

}

type CounterWrapper struct {
	writer io.Writer
	n      *int64
}

func (c *CounterWrapper) Write(p []byte) (int, error) {
	n, err := c.writer.Write(p)
	*c.n += int64(n)
	return n, err
}
