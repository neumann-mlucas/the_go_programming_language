package main

import (
	"fmt"
	"os"
	"sort"
	"strings"
)

var prereqs = map[string][]string{
	"algorithms": {"data structures"},
	"calculus":   {"linear algebra"},
	"compilers": {
		"data structures",
		"formal languages",
		"computer organization",
	},
	"data structures":       {"discrete math"},
	"databases":             {"data structures"},
	"discrete math":         {"intro to programming"},
	"formal languages":      {"discrete math"},
	"networks":              {"operating systems"},
	"operating systems":     {"data structures", "computer organization"},
	"programming languages": {"data structures", "computer organization"},
	"linear algebra":        {"calculus"},        // a circle
	"intro to programming":  {"data structures"}, // another circle
}

func main() {
	order, err := topoSort(prereqs)
	if err != nil {
		fmt.Fprint(os.Stderr, err)
	}
	for i, course := range order {
		fmt.Printf("%d:\t%s\n", i+1, course)
	}
}

func topoSort(m map[string][]string) ([]string, error) {
	var order []string
	seen := make(map[string]bool)

	var visitAll func(items []string)
	visitAll = func(items []string) {
		for _, item := range items {
			if !seen[item] {
				seen[item] = true
				visitAll(m[item])
				order = append(order, item)
			}
		}
	}

	var keys []string
	for key := range m {
		path, circle := detectCircle(key, m)
		if circle {
			return nil, fmt.Errorf("Circle detect: %s", strings.Join(path, " => "))
		}
		keys = append(keys, key)
	}
	sort.Strings(keys)
	visitAll(keys)

	return order, nil
}

func detectCircle(key string, m map[string][]string) ([]string, bool) {
	seen := make(map[string]bool)
	var hasCircle bool
	var path []string

	var visitAll func(item string)
	visitAll = func(item string) {
		if !seen[item] {
			seen[item] = true
			path = append(path, item)
			for _, prereq := range m[item] {
				hasCircle = key == prereq
				visitAll(prereq)
			}
		}
	}

	visitAll(key)
	return path, hasCircle
}
