package main

import (
	"crypto/sha256"
	"fmt"
)

func main() {
	c1 := sha256.Sum256([]byte("x"))
	c2 := sha256.Sum256([]byte("a"))
	fmt.Println(compareSHA256(&c1, &c2))
}

func compareSHA256(ptr1, ptr2 *[32]byte) uint {
	diff := uint(0)
	for i := 0; i < 32; i++ {
		if ptr1[i] != ptr2[i] {
			diff++
		}
	}
	return diff
}
