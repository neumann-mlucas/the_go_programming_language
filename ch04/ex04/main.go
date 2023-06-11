package main

import "fmt"

func main() {
	arr := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	fmt.Println(arr)
	rotate(arr, 10)
	fmt.Println(arr)
}

func rotate(s []int, n int) {
	temp := make([]int, len(s))
	for i := 0; i < len(s); i++ {
		j := (i + n) % len(s)
		temp[j] = s[i]
	}
	copy(s, temp)
}
