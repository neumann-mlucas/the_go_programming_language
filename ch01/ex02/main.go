package main

import (
	"fmt"
	"os"
)

func main() {
	s, sep := "", ""
	for i, arg := range os.Args {
        s += fmt.Sprint(sep, i," ", arg)
		sep = "\n"
	}
	fmt.Println(s)
}
