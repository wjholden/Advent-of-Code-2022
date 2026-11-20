package main

import (
	_ "embed"
	"errors"
	"fmt"
	"log"
	"math"
	"strings"
	"unicode"

	"gonum.org/v1/gonum/mat"
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

type Face struct {
	min_r, max_r, min_c, max_c int
}

type Cube struct {
	T                      [6]mat.Matrix
	H                      [6]mat.Matrix
	N                      [6]mat.Matrix
	board                  *Board
	length                 int
	ft, fu, fl, fr, fb, ff Face
}

func (c *Cube) GetFace(p *Point) (int, error) {
	switch {
	case p.IsInside(&c.ft):
		return 0, nil
	case p.IsInside(&c.fu):
		return 1, nil
	case p.IsInside(&c.fl):
		return 2, nil
	case p.IsInside(&c.fr):
		return 3, nil
	case p.IsInside(&c.fb):
		return 4, nil
	case p.IsInside(&c.ff):
		return 5, nil
	default:
		return -1, errors.New("this point is not on a face")
	}
}

func (c *Cube) IsOn(xyz mat.VecDense) bool {
	x := xyz.At(0, 0)
	y := xyz.At(1, 0)
	z := xyz.At(2, 0)
	return 0 <= x && x <= float64(c.length) &&
		0 <= y && y <= float64(c.length) &&
		0 <= z && z <= float64(c.length)
}

func (p *Point) IsInside(f *Face) bool {
	return f.min_r <= p.row && p.row <= f.max_r &&
		f.min_c <= p.col && p.col <= f.max_c
}

func (p *Point) move3D(dr, dc int, cube *Cube) Point {
	// First, try to just push it in the direction it was going.
	p_new := Point{row: p.row + dr, col: p.col + dc}
	_, err := cube.GetFace(&p_new)
	if err != nil {
		// We're still on the cube. Is there a wall on this position?
		if _, ok := cube.board.walls[p_new]; ok {
			// Yes, there is a wall here.
			return *p
		} else {
			return p_new
		}
	}
	// We're off the cube. We need to warp to another side.
	face, err := cube.GetFace(p)
	if err != nil {
		log.Fatal("invariant violation - we got off the cube")
	}
	t := cube.T[face]
	v := mat.NewVecDense(4, []float64{float64(p.row), float64(p.col), float64(cube.length), 1})
	xyz := mat.NewVecDense(3, nil)
	xyz.MulVec(t, v)
	h := cube.H[face]
	d_xyz := mat.NewVecDense(3, nil)
	d_rc := mat.NewVecDense(2, []float64{float64(dr), float64(dc)})
	d_xyz.MulVec(h, d_rc)
	p_xyz := mat.NewVecDense(3, nil)
	p_xyz.AddVec(xyz, d_xyz)
	// Are we still on a face?
	if cube.IsOn(*p_xyz) {
		// Is there a wall here?
		rcl1 := mat.NewVecDense(4, nil)
		if rcl1.SolveVec(h, p_xyz) != nil {
			log.Fatal("failed to reverse row and column")
		}
		p_new.row = int(rcl1.AtVec(0))
		p_new.col = int(rcl1.AtVec(1))
		if _, ok := cube.board.walls[p_new]; ok {
			// There was a wall in our way.
			return *p
		}
		// There was no wall in our way.
		return p_new
	}
	// We need to warp.
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

// I really miss Rust's enumerated type right now.
type Instruction interface {
	isInstruction()
}

type GoLeft struct{}
type GoRight struct{}
type GoForward struct{ x int }

func (GoLeft) isInstruction()    {}
func (GoRight) isInstruction()   {}
func (GoForward) isInstruction() {}

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
		case GoLeft:
			dr, dc = -dc, dr // dx, dy = -dy, dx but we are in rows and columns.
		case GoRight:
			dr, dc = dc, -dr
		case GoForward:
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
			instruction = GoForward{x}
		case a == 'R':
			instruction = GoRight{}
		case a == 'L':
			instruction = GoLeft{}
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
