package main

import (
	_ "embed"
	"testing"
)

//go:embed example.txt
var example string

// Table-Driven Testing
// https://go.dev/wiki/TableDrivenTests
func TestParser(t *testing.T) {
	var parsetests = []struct {
		in  string
		out Snafu
	}{
		{"1", 1},
		{"2", 2},
		{"1=", 3},
		{"1-", 4},
		{"10", 5},
		{"11", 6},
		{"12", 7},
		{"2=", 8},
		{"2-", 9},
		{"20", 10},
		{"1=0", 15},
		{"1-0", 20},
		{"1=11-2", 2022},
		{"1-0---0", 12345},
		{"1121-1110-1=0", 314159265},
	}

	for _, tt := range parsetests {
		t.Run(tt.in, func(t *testing.T) {
			got := NewSnafu(tt.in)
			if got != tt.out {
				t.Errorf("got %d, want %d", got, tt.out)
			}
		})
	}
}

func Test1(t *testing.T) {
	want := 4890
	got := 0
	for _, s := range parse(example) {
		got += int(s)
	}
	if got != want {
		t.Errorf("incorrect sum - got %d, want %d", got, want)
	}
}
