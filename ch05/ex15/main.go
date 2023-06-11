package main

import "fmt"

func main() {
	args := []int{1, 2, 3, 4}
	fmt.Println("max: ", max(args...))
	fmt.Println("max: ", min(args...))
}

func max(args ...int) int {
	if len(args) == 0 {
		panic("No Args")
	}
	argmax := args[0]
	for _, val := range args[1:] {
		if val > argmax {
			argmax = val
		}
	}
	return argmax
}

func min(args ...int) int {
	if len(args) == 0 {
		panic("No Args")
	}
	argmin := args[0]
	for _, val := range args[1:] {
		if val < argmin {
			argmin = val
		}
	}
	return argmin
}

func max2(base int, args ...int) int {
	argmax := base
	for _, val := range args {
		if val > argmax {
			argmax = val
		}
	}
	return argmax
}
