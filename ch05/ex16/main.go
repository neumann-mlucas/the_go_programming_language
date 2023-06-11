package main

import "fmt"

func main() {
	args := []string{"Hello", "World!"}
	fmt.Println("args: ", args)
	fmt.Printf("out: \"%s\"\n", strJoin(" ", args...))
}

func strJoin(sep string, args ...string) string {
	var result string
	for n, arg := range args {
		result += arg
		if n < len(args)-1 {
			result += sep
		}
	}
	return result
}
