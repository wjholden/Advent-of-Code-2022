package main

import (
	_ "embed"
	"fmt"
	"strings"
)

//go:embed day25.txt
var puzzle string

type Snafu int

func NewSnafu(s string) Snafu {
	v := 0

	for i, c := range s {
		d := 0
		switch c {
		case '2':
			d = 2
		case '1':
			d = 1
		case '0':
			d = 0
		case '-':
			d = -1
		case '=':
			d = -2
		}
		for range len(s) - i - 1 {
			d *= 5
		}
		v += d
	}

	return Snafu(v)
}

func Reformat(s string) Snafu {
	s = strings.ReplaceAll(s, "-2", "=")
	s = strings.ReplaceAll(s, "-1", "-")
	fmt.Println(s)
	return NewSnafu(s)
}

func parse(input string) []Snafu {
	results := make([]Snafu, 0)
	for line := range strings.Lines(input) {
		snafu := NewSnafu(strings.TrimSpace(line))
		results = append(results, snafu)
	}
	return results
}

func main() {
	got := 0
	for _, s := range parse(puzzle) {
		got += int(s)
	}
	fmt.Println(got)
	fmt.Println(NewSnafu("20-1-0=-2=-2220=0011"))
	fmt.Println("So there was an efficient algorithm to do this...I just didn't know it.")
}
