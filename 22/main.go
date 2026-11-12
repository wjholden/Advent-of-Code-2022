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
	board := NewBoard(puzzle)
	fmt.Println("Part 1:", board.part1())
}

type Bound struct {
	min int
	max int
}

type Point struct {
	row int
	col int
}

func (p *Point) move(dr, dc int, board *Board) Point {
	row := p.row + dr
	col := p.col + dc

	if dr != 0 {
		if row < board.col_bounds[col].min {
			row = board.col_bounds[col].max
		}
		if row > board.col_bounds[col].max {
			row = board.col_bounds[col].min
		}
	}

	if dc != 0 {
		if col < board.row_bounds[row].min {
			col = board.row_bounds[row].max
		}
		if col > board.row_bounds[row].max {
			col = board.row_bounds[row].min
		}
	}

	return Point{row, col}
}

type Face map[Dir]Dir

type Pattern int

const (
	EXAMPLE Pattern = iota
	PUZZLE
)

type Cube struct {
	pattern   Pattern
	length    int
	TopFace   Face
	LeftFace  Face
	BackFace  Face
	UnderFace Face
	RightFace Face
	FrontFace Face
}

type Dir struct {
	dr int
	dc int
}

var (
	UpDir    = Dir{dr: -1, dc: 0}
	DownDir  = Dir{dr: 1, dc: 0}
	LeftDir  = Dir{dr: 0, dc: -1}
	RightDir = Dir{dr: 0, dc: 1}
)

var ExampleCube = Cube{
	pattern: EXAMPLE,
	length: 4,
	TopFace: ,
}

func (c *Cube) move(p Point, dr, dc int, board *Board) (Point, int, int) {
	// 1) Find the next position.
	// Warping isn't necessary if we don't change tiles.
	// Does the move cross a boundary? If not, just add.
	// If yes, then look up the destination face in a table.

	// 2)
}

func (c *Cube) face(p Point) *Face {
	row := p.row / c.length
	col := p.col / c.length
	switch c.pattern {
	case EXAMPLE:
		switch {
		case row == 0 && col == 2:
			return &c.TopFace
		case row == 1 && col == 0:
			return &c.BackFace
		case row == 1 && col == 1:
			return &c.LeftFace
		case row == 1 && col == 2:
			return &c.FrontFace
		case row == 2 && col == 2:
			return &c.UnderFace
		case row == 2 && col == 3:
			return &c.RightFace
		default:
			log.Fatal("could not find a face for this position on the cube")
		}
	case PUZZLE:
		switch {
		case row == 0 && col == 1:
			return &c.TopFace
		case row == 0 && col == 2:
			return &c.RightFace
		case row == 1 && col == 1:
			return &c.FrontFace
		case row == 2 && col == 0:
			return &c.LeftFace
		case row == 2 && col == 1:
			return &c.UnderFace
		case row == 3 && col == 0:
			return &c.BackFace
		default:
			log.Fatal("could not find a face for this position on the cube")
		}
	default:
		log.Fatal("this program uses hardcoded patterns for layout (sorry, this is a hard one)")
	}
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
	// How far a row can go left and right
	row_bounds []Bound

	// How far a column can go up and down
	col_bounds   []Bound
	walls        map[Point]bool
	path         string
	instructions []Instruction

	cube_len int
}

func (board *Board) part1() int {
	position := Point{row: 0, col: board.row_bounds[0].min}
	dr, dc := 0, 1

	for _, instruction := range board.instructions {
		switch i := instruction.(type) {
		case Left:
			dr, dc = -dc, dr // dx, dy = -dy, dx but we are in rows and columns.
		case Right:
			dr, dc = dc, -dr
		case Forward:
			for x := i.x; x > 0; x-- {
				p := position.move(dr, dc, board)
				if wall, ok := board.walls[p]; ok && wall {
					break
				} else {
					position = p
				}
			}
		}
	}

	return 1000*(1+position.row) + 4*(1+position.col) + Facing(dr, dc)
}

func Facing(dr, dc int) (facing int) {
	switch {
	case dr == 0 && dc == 1:
		facing = 0
	case dr == 1 && dc == 0:
		facing = 1
	case dr == 0 && dc == -1:
		facing = 2
	case dr == -1 && dc == 0:
		facing = 3
	default:
		log.Fatal("invalid direction")
	}
	return
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
		right = len(line) - 1

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
	board.walls = make(map[Point]bool)
	row = 0
	col := 0
	for _, c := range input {
		switch c {
		case '\n':
			row++
			col = 0
			continue
		case '#':
			board.walls[Point{row, col}] = true
		case '.':
			board.walls[Point{row, col}] = false
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

	// the cube faces are either unfolded 4 tall/3 wide or 3 tall/4 wide.
	board.cube_len = width - height
	if board.cube_len < 0 {
		board.cube_len = -board.cube_len
	}

	return
}
