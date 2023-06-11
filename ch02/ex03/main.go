package main

import (
	"fmt"
	"time"
)

// pc[i] is the population count of i.
var pc [256]byte

func init() {
	for i := range pc {
		pc[i] = pc[i/2] + byte(i&1)
	}
}

// PopCount returns the population count (number of set bits) of x.
func PopCount(x uint64) int {
	return int(pc[byte(x>>(0*8))] +
		pc[byte(x>>(1*8))] +
		pc[byte(x>>(2*8))] +
		pc[byte(x>>(3*8))] +
		pc[byte(x>>(4*8))] +
		pc[byte(x>>(5*8))] +
		pc[byte(x>>(6*8))] +
		pc[byte(x>>(7*8))])
}

// PopCountLoop is an alternative implementation of PopCount.
func PopCountLoop(x uint64) int {
	var n byte
	for i := uint(0); i < 8; i++ {
		n += pc[byte(x>>(i*8))]
	}
	return int(n)
}

func main() {

	start := time.Now()
	for i := uint64(0); i < 100_000_000; i++ {
		PopCount(i)

	}
	fmt.Printf("PopCount\t-> time elapsed: %.2fs\n", time.Since(start).Seconds())

	start = time.Now()
	for i := uint64(0); i < 100_000_000; i++ {
		PopCountLoop(i)

	}
	fmt.Printf("PopCountLoop\t-> time elapsed: %.2fs\n", time.Since(start).Seconds())
}
