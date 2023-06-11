package main

import (
	"bytes"
	"fmt"
)

func main() {
	var x IntSet
	x.Add(2)
	x.Add(4)
	x.Add(128)
	x.Add(1024)
	fmt.Println("x = ", x.String(), "length of", x.Len())

	x.Remove(1024)
	x.Remove(4)
	fmt.Println("x = ", x.String(), "length of", x.Len())

	t := x.Copy()
	fmt.Println("t = ", t.String(), "length of", t.Len())

	x.Clear()
	fmt.Println("x = ", x.String(), "length of", x.Len())
	fmt.Println("t = ", t.String(), "length of", t.Len())
}

// An IntSet is a set of small non-negative integers.
type IntSet struct {
	words []uint64
}

// Has reports whether the set contains the non-negative value x.
func (s *IntSet) Has(x int) bool {
	word, bit := x/64, uint(x%64)
	return word < len(s.words) && s.words[word]&(1<<bit) != 0
}

// Add adds the non-negative value x to the set.
func (s *IntSet) Add(x int) {
	word, bit := x/64, uint(x%64)
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
		for j := 0; j < 64; j++ {
			if word&(1<<uint(j)) != 0 {
				if buf.Len() > len("{") {
					buf.WriteByte(' ')
				}
				fmt.Fprintf(&buf, "%d", 64*i+j)
			}
		}
	}
	buf.WriteByte('}')
	return buf.String()
}

// Len returns the length of the set
func (s *IntSet) Len() int {
	var length int
	for _, word := range s.words {
		if word == 0 {
			continue
		}
		for j := 0; j < 64; j++ {
			if word&(1<<uint(j)) != 0 {
				length++
			}
		}
	}
	return length
}

// Remove removes the non-negative value x to the set.
func (s *IntSet) Remove(x int) {
	if s.Has(x) {
		word, bit := x/64, uint(x%64)
		s.words[word] ^= 1 << bit
	}
}

// Clear removes all elements from the set
func (s *IntSet) Clear() {
	for i := range s.words {
		s.words[i] = 0
	}
}

// Copy returns a copy of the set
func (s *IntSet) Copy() *IntSet {
	t := new(IntSet)
	t.words = make([]uint64, len(s.words))
	copy(t.words, s.words)

	return t
}
