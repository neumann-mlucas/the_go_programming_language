package main

import (
	"fmt"
	"io"
	"os"
)

var stdout io.Writer = os.Stdout

func main() {
	fmt.Printf("return: %d\n", noReturn())
}

func noReturn() (r int) {
	defer func() {
		recover()
	}()
	panic("oops!")
}
