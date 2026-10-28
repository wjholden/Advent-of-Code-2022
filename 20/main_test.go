package main

import "testing"

var example string = `1
2
-3
3
-2
0
4`

func Test1(t *testing.T) {
	want := 4 - 3 + 2
	got := part1(&example)
	if want != got {
		t.Errorf("incorrect result for part 1 (got %d, want %d)", got, want)
	}
}

func TestLength(t *testing.T) {
	ct := LLparse(&example)
	in := parse(&example)
	if ct.Len() != len(in) {
		t.Error("something is wrong with your length calculation")
	}
}
