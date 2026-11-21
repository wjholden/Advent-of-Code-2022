package main

import (
	_ "embed"
	"testing"

	"gonum.org/v1/gonum/mat"
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

func Test2(t *testing.T) {
	board := NewBoard(example)

	l := float64(board.cube_len)

	top := Vec3{0, 1, 0}
	front := Vec3{0, 0, 1}
	left := Vec3{-1, 0, 0}
	back := Vec3{0, 0, -1}
	under := Vec3{0, -1, 0}
	right := Vec3{1, 0, 0}

	A := map[Vec3]mat.Matrix{
		top:   mat.NewDense(3, 2, []float64{0, 1, 0, 0, 1, 0}),   // was {0, 0, 1, 1, 0, 0}
		front: mat.NewDense(3, 2, []float64{0, 1, -1, 0, 0, 0}),  // was {0, -1, 0, 1, 0, 0}
		left:  mat.NewDense(3, 2, []float64{0, 0, -1, 0, 0, 1}),  // was {0, -1, 0, 0, 0, 1}
		back:  mat.NewDense(3, 2, []float64{0, -1, -1, 0, 0, 0}), // was {0, -1, 0, -1, 0, 0}
		under: mat.NewDense(3, 2, []float64{0, 1, 0, 0, -1, 0}),  // was {0, 0, -1, 1, 0, 0}
		right: mat.NewDense(3, 2, []float64{0, 0, 0, 1, -1, 0}),  // was {0, 0, -1, 0, 1, 0}
	}

	b := map[Vec3]mat.VecDense{
		top:   *mat.NewVecDense(3, []float64{-2 * l, l - 1, 0}),
		front: *mat.NewVecDense(3, []float64{-2 * l, 2*l - 1, l - 1}),
		left:  *mat.NewVecDense(3, []float64{0, 2*l - 1, -l}),
		back:  *mat.NewVecDense(3, []float64{l - 1, 2*l - 1, 0}),
		under: *mat.NewVecDense(3, []float64{-2 * l, 0, 3*l - 1}),
		right: *mat.NewVecDense(3, []float64{l - 1, -3 * l, 3*l - 1}),
	}

	cube := BetterCube{
		board: &board,
		A:     A,
		b:     b,
	}

	got := cube.part2()
	want := 1000*5 + 4*7 + 3
	if got != want {
		t.Errorf("incorrect result for part 2 (got: %d, want %d)", got, want)
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
