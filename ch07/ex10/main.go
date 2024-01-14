package main

import (
	"fmt"
	"sort"
)

func IsPalindrome(s sort.Interface) bool {
	for i := 0; i < s.Len(); i++ {
		j := s.Len() - 1 - i
		if !(!s.Less(i, j) && !s.Less(j, i)) {
			return false
		}
	}
	return true
}

func main() {
	var seq sort.Interface

	seq = sort.IntSlice([]int{2, 1, 0, 1, 3})
	fmt.Printf("seq (%v) is a palindrome: %v\n", seq, IsPalindrome(seq))

	seq = sort.IntSlice([]int{2, 1, 0, 1, 2})
	fmt.Printf("seq (%v) is a palindrome: %v\n", seq, IsPalindrome(seq))
}
