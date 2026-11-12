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
	got := inspect(Left{})
	want := "Left"
	if got != want {
		t.Error("unexpected type", got)
	}
}

func TestRight(t *testing.T) {
	got := inspect(Right{})
	want := "Right"
	if got != want {
		t.Error("unexpected type", got)
	}
}

func TestForward(t *testing.T) {
	got := inspect(Forward{})
	want := "Forward"
	if got != want {
		t.Error("unexpected type", got)
	}
}

func inspect(instruction Instruction) (t string) {
	switch instruction.(type) {
	case Left:
		t = "Left"
	case Right:
		t = "Right"
	case Forward:
		t = "Forward"
	default:
		t = "unexpected"
	}
	return
}
