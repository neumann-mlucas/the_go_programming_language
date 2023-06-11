package main

import (
	"bytes"
	"fmt"
)

const (
	INTSIZE = 32 << (^uint(0) >> 63)
)

func main() {
	var x IntSet
	x.Add(2)
	x.Add(4)
	x.Add(128)
	x.Add(1024)
	fmt.Println(x.String())
}

// An IntSet is a set of small non-negative integers.
type IntSet struct {
	words []uint
}

// Has reports whether the set contains the non-negative value x.
func (s *IntSet) Has(x int) bool {
	word, bit := x/INTSIZE, uint(x%INTSIZE)
	return word < len(s.words) && s.words[word]&(1<<bit) != 0
}

// Add adds the non-negative value x to the set.
func (s *IntSet) Add(x int) {
	word, bit := x/INTSIZE, uint(x%INTSIZE)
	for word >= len(s.words) {
		s.words = append(s.words, 0)
	}
	s.words[word] |= 1 << bit
}

// UnionWith sets s to the union of s and t.
func (s *IntSet) UnionWith(t *IntSet) {
	for i, tword := range t.words {
		if i < len(s.words) {
			s.words[i] |= tword
		} else {
			s.words = append(s.words, tword)
		}
	}
}

// String returns the set as a string of the form "{1 2 3}".
func (s *IntSet) String() string {
	var buf bytes.Buffer
	buf.WriteByte('{')
	for i, word := range s.words {
		if word == 0 {
			continue
		}
		for j := 0; j < INTSIZE; j++ {
			if word&(1<<uint(j)) != 0 {
				if buf.Len() > len("{") {
					buf.WriteByte(' ')
				}
				fmt.Fprintf(&buf, "%d", INTSIZE*i+j)
			}
		}
	}
	buf.WriteByte('}')
	return buf.String()
}
