package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func main() {
	b, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		return
	}
	result := expand("foo", string(b), strings.ToUpper)
	fmt.Println(result)
}

func expand(sub, str string, f func(string) string) string {
	var result []string
	for _, substring := range splitString(sub, str) {
		if substring == sub {
			substring = f(substring)
		}
		result = append(result, substring)
	}
	return strings.Join(result, "")

}

func splitString(sub, str string) []string {
	var lastPos int
	var output []string
	for _, i := range findSubstrings(sub, str) {
		output = append(output, str[lastPos:i])
		output = append(output, str[i:i+len(sub)])
		lastPos = i + len(sub)
	}
	output = append(output, str[lastPos:])
	return output
}

func findSubstrings(sub, str string) []int {
	var idxs []int
	for i := 0; i < len(str); i++ {
		if hasPrefix(sub, str[i:]) {
			idxs = append(idxs, i)
			i += len(sub) - 1
		}
	}
	return idxs
}

func hasPrefix(prefix, str string) bool {
	if len(str) < len(prefix) {
		return false
	}
	for i := 0; i < len(prefix); i++ {
		if prefix[i] != str[i] {
			return false
		}
	}
	return true
}
