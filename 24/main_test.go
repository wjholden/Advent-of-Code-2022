package main

import (
	_ "embed"
	"testing"
)

//go:embed complex_example.txt
var complex_example string

func Test1(t *testing.T) {
	want := 18
	blizzards := parse(complex_example)
	got := bfs(blizzards, start(), end())
	if want != got {
		t.Errorf("incorrect result (got: %d, want %d)", got, want)
	}
}

func Test2(t *testing.T) {
	blizzards := parse(complex_example)
	x := bfs(blizzards, start(), end())
	if x != 18 {
		t.Errorf("incorrect result (got: %d, want %d)", x, 18)
	}
	y := bfs(blizzards, end(), start())
	if y != 23 {
		t.Errorf("incorrect result (got: %d, want %d)", y, 23)
	}
	z := bfs(blizzards, start(), end())
	if z != 13 {
		t.Errorf("incorrect result (got: %d, want %d)", z, 13)
	}
}
