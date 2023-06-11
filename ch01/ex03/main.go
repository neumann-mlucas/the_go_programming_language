package main

import (
	"fmt"
	"os"
	"strings"
	"time"
)

func main() {
	fastEcho(os.Args)
	if false {
		measureEchoTime(fastEcho)
		measureEchoTime(slowEcho)
	}
}

func measureEchoTime(echo func([]string)) {
	start := time.Now()
	echo(os.Args)
	t := time.Now()
	fmt.Println(t.Sub(start))
}
func fastEcho(args []string) {
	s := strings.Join(args, " ")
	fmt.Println(s)
}

func slowEcho(args []string) {
	s, sep := "", ""
	for _, arg := range args {
		s += sep + arg
		sep = " "
	}
	fmt.Println(s)
}
