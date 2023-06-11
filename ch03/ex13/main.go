package main

import "fmt"

const (
	Kbyte = 1000
	Mbyte = Kbyte * Kbyte
	Gbyte = Mbyte * Kbyte
	Tbyte = Gbyte * Kbyte
)

func main() {
	fmt.Println(Kbyte, Mbyte, Gbyte, Tbyte)
}
