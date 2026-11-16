package main

import (
	_ "embed"
	"testing"
)

//go:embed example.txt
var example string

func Test1(t *testing.T) {
	board := NewBoard(example)
	got := board.part1()
	want := 1000*6 + 4*8 + 0
	if got != want {
		t.Errorf("incorrect result for part 1 (got: %d, want %d)", got, want)
	}
}

func TestLeft(t *testing.T) {
	got := inspect(GoLeft{})
	want := "Left"
	if got != want {
		t.Error("unexpected type", got)
	}
}

func TestRight(t *testing.T) {
	got := inspect(GoRight{})
	want := "Right"
	if got != want {
		t.Error("unexpected type", got)
	}
}

func TestForward(t *testing.T) {
	got := inspect(GoForward{})
	want := "Forward"
	if got != want {
		t.Error("unexpected type", got)
	}
}

func inspect(instruction Instruction) (t string) {
	switch instruction.(type) {
	case GoLeft:
		t = "Left"
	case GoRight:
		t = "Right"
	case GoForward:
		t = "Forward"
	default:
		t = "unexpected"
	}
	return
}
