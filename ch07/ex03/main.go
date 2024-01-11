package main

import "fmt"

func main() {
	left, right := &tree{2, nil, nil}, &tree{3, nil, nil}
	root := &tree{1, left, right}

	fmt.Println(root)

}

type tree struct {
	value       int
	left, right *tree
}

func (t *tree) String() string {
	base := fmt.Sprintf("-> %d ", t.value)

	if t.left != nil {
		base += t.left.String()
	}
	if t.right != nil {
		base += t.right.String()
	}
	return base
}
