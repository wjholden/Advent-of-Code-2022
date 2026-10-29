package main

import "testing"

var example string = `1
2
-3
3
-2
0
4`

// https://www.reddit.com/r/adventofcode/comments/zr29qd/comment/j11bcdy/
var extra1 string = `0
-1
-1
1`

var extra2 string = `1
2
-3
3
-2
0
8`

func Test1(t *testing.T) {
	want := 4 - 3 + 2
	got := part1(&example)
	if want != got {
		t.Errorf("incorrect result for part 1 (got %d, want %d)", got, want)
	}
}

func TestExtra1(t *testing.T) {
	want := 0
	got := part1(&extra1)
	if want != got {
		t.Errorf("incorrect result for extra testcase (got %d, want %d)", got, want)
	}
}

func TestExtra2(t *testing.T) {
	want := 7
	got := part1(&extra2)
	if want != got {
		t.Errorf("incorrect result for extra testcase (got %d, want %d)", got, want)
	}
}

func TestLength(t *testing.T) {
	ct := LLparse(&example)
	in := parse(&example)
	if ct.Len() != len(in) {
		t.Error("something is wrong with your length calculation")
	}
}
