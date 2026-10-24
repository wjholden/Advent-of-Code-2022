package main

import (
	_ "embed"
	"testing"
)

//go:embed example.txt
var example string

func Test1(t *testing.T) {
	cubes := parse(example)
	result, _ := solve(&cubes)
	expect := 64
	if result != expect {
		t.Errorf("incorrect result for part 1 (got %d, expected %d)", result, expect)
	}
}

func Test2(t *testing.T) {
	cubes := parse(example)
	_, result := solve(&cubes)
	expect := 58
	if result != expect {
		t.Errorf("incorrect result for part 2 (got %d, expected %d)", result, expect)
	}
}
