package main

import (
	_ "embed"
	"fmt"
	"log"
	"math"
	"strings"
	"unicode"
)

//go:embed day22.txt
var puzzle string

func main() {
	fmt.Println(NewBoard(puzzle))
}

type Bound struct {
	min int
	max int
}

type Point struct {
	row int
	col int
}

// I really miss Rust's enumerated type right now.
type Instruction interface {
	isInstruction()
}

type Left struct{}
type Right struct{}
type Forward struct{ x int }

func (Left) isInstruction()    {}
func (Right) isInstruction()   {}
func (Forward) isInstruction() {}

type Board struct {
	row_bounds   []Bound
	col_bounds   []Bound
	walls        map[Point]struct{}
	path         string
	instructions []Instruction
}

func (board *Board) part1() int {
	position := Point{row: 0, col: board.row_bounds[0].min}
	direction := Point{row: 1, col: 0}

}

func BoardDimensions(input string) (height, width int) {
	for line := range strings.Lines(input) {
		line = strings.TrimSuffix(line, "\n")
		if len(line) == 0 {
			break
		}
		height++
		width = max(width, len(line))
	}
	fmt.Println("board dimensions", height, width)
	return
}

func NewBoard(input string) (board Board) {
	height, width := BoardDimensions(input)

	board.row_bounds = make([]Bound, height)
	board.col_bounds = make([]Bound, width)

	row := 0
	for line := range strings.Lines(input) {
		var left, right int
		line = strings.TrimSuffix(line, "\n")

		if len(line) == 0 {
			break
		}

		// find first non-space
		for r, c := range line {
			if !unicode.IsSpace(c) {
				left = r
				break
			}
		}

		// there are no trailing spaces in the input
		right = len(line)

		board.row_bounds[row].min = left
		board.row_bounds[row].max = right

		row++
	}

	for c := range width {
		lo := math.MaxInt
		hi := 0
		for r, rb := range board.row_bounds {
			if rb.min <= c && c <= rb.max {
				lo = min(lo, r)
				hi = max(hi, r)
			}
		}
		board.col_bounds[c].min = lo
		board.col_bounds[c].max = hi
	}

	// now to find the walls in the input
	board.walls = make(map[Point]struct{})
	row = 0
	col := 0
	for _, c := range input {
		switch c {
		case '\n':
			row++
			col = 0
			continue
		case '#':
			board.walls[Point{row, col}] = struct{}{}
		}
		col++
	}

	// finally, the path
	tokens := strings.Split(input, "\n\n")
	if len(tokens) != 2 {
		log.Fatal()
	}
	board.path = strings.TrimSpace(tokens[1])

	board.instructions = make([]Instruction, 0)
	for i := 0; i < len(board.path); i++ {
		var instruction Instruction
		a := rune(board.path[i])
		switch {
		case unicode.IsNumber(a):
			x := int(a - '0')
			// maybe the next one is also a letter
			if i+1 < len(board.path) {
				if b := rune(board.path[i+1]); unicode.IsNumber(b) {
					x = x*10 + int(b-'0')
					i++
				}
			}
			instruction = Forward{x}
		case a == 'R':
			instruction = Right{}
		case a == 'L':
			instruction = Left{}
		default:
			log.Fatal("unexpected instruction", a)
		}
		board.instructions = append(board.instructions, instruction)
	}

	return
}
