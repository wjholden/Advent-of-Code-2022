package main

import (
	_ "embed"
	"testing"
)

//go:embed example_large.txt
var example_large string

//go:embed example_small.txt
var example_small string

func Test1(t *testing.T) {
	want := 110
	got, _ := solve(parse(example_large))
	if want != got {
		t.Errorf("incorrect result for part 1 (got %d, want %d)", got, want)
	}
}

func TestSmall(t *testing.T) {
	want := 5*6 - 5
	got, _ := solve(parse(example_small))
	if want != got {
		t.Errorf("incorrect result for part 1 (got %d, want %d)", got, want)
	}
}

func Test2(t *testing.T) {
	want := 20
	_, got := solve(parse(example_large))
	if want != got {
		t.Errorf("incorrect result for part 2 (got %d, want %d)", got, want)
	}
}
