package main

import (
	_ "embed"
	"fmt"
	"testing"
)

//go:embed example.txt
var example string

func Test1(t *testing.T) {
	fmt.Println(NewBoard(example))
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
