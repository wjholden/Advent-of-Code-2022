// https://go.dev/play/p/RsdHBqnRcW2
package main

import (
	_ "embed"
	"fmt"
	"strings"
)

//go:embed day06.txt
var puzzle string

func main() {
	puzzle := strings.ReplaceAll(puzzle, "\n", "")
	fmt.Printf("Part 1: %d\n", findMarker(puzzle, 4))
	fmt.Printf("Part 2: %d\n", findMarker(puzzle, 14))
}

func findMarker(s string, n int) int {
	for i := n; i < len(s); i++ {
		if distinct(s[(i-n):i], n) {
			return i
		}
	}
	return -1
}

func distinct(s string, n int) bool {
	set := make(map[rune]bool)
	for _, r := range s {
		set[r] = true
	}
	return len(set) == n
}

func segment(s string) {
	for i := 0; i < len(s); i += 80 {
		fmt.Println(s[i:min(i+80, len(s))])
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
