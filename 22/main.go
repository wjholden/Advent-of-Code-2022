package main

import (
	_ "embed"
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

type Vec3 struct {
	x, y, z int
}

func (v *Vec3) ToDense() mat.VecDense {
	return *mat.NewVecDense(3, []float64{
		float64(v.x),
		float64(v.y),
		float64(v.z),
	})
}

func NewVec3(v *mat.VecDense) Vec3 {
	return Vec3{
		x: int(math.Round(v.AtVec(0))),
		y: int(math.Round(v.AtVec(1))),
		z: int(math.Round(v.AtVec(2))),
	}
}

type BetterCube struct {
	// Association of normal vectors to heading transformations.
	A map[Vec3]mat.Matrix
	// Additional change in x, y, and z necessary to transform from row/col.
	b     map[Vec3]mat.VecDense
	board *Board
}

func (c *BetterCube) IsOn(xyz *mat.VecDense) bool {
	x := xyz.AtVec(0)
	y := xyz.AtVec(1)
	z := xyz.AtVec(2)
	l := float64(c.board.cube_len) - 1
	// At least one must be either zero or the length value.
	if !(x == 0 || x == l || y == 0 || y == l || z == 0 || z == l) {
		return false
	}
	return 0 <= x && x <= l && 0 <= y && y <= l && 0 <= z && z <= l
}

func (c *BetterCube) GetRowCol(xyz, N *mat.VecDense) Point {
	n := NewVec3(N)
	A := c.A[n]
	b := c.b[n]
	rc := mat.NewVecDense(2, nil)
	y := mat.NewVecDense(3, nil)
	y.SubVec(xyz, &b)
	err := rc.SolveVec(A, y)
	if err != nil {
		log.Fatal("failed to solve linear system")
	}
	return Point{
		row: int(math.Round(rc.AtVec(0))),
		col: int(math.Round(rc.AtVec(1))),
	}
}

func (c *BetterCube) move(p Vec3, H *mat.VecDense, N Vec3) (Vec3, *mat.VecDense, Vec3) {
	// If we follow the current heading, will we fall off the cube?
	p1 := p.ToDense()
	p2 := mat.NewVecDense(3, nil)
	p2.AddVec(&p1, H)
	// Does this take us off of the cube?
	if c.IsOn(p2) {
		// No. Check if there is a wall at this new position.
		n := N.ToDense()
		point := c.GetRowCol(p2, &n)
		if isWall, ok := c.board.walls[point]; ok && isWall {
			// This is a wall. We can't move.
			return p, H, N
		} else {
			// No wall. Move to the new position.
			return NewVec3(p2), H, N
		}
	}
	// This would have taken us off the cube. We need to rotate around to the other side.

	// H' = -N, N' = H
	h_new := mat.NewVecDense(3, []float64{float64(N.x) * -1, float64(N.y) * -1, float64(N.z) * -1})
	n_new := NewVec3(H)

	// Try to move. If we get the same position back, then we know that
	// there was a wall right around the corner that we couldn't go through.
	n_next := n_new.ToDense()
	p_next := c.GetRowCol(&p1, &n_next)
	if isWall, ok := c.board.walls[p_next]; ok && isWall {
		// There is a wall right around the corner. Stay exactly where you are.
		return p, H, N
	} else {
		// No wall around the corner. Don't move, but return the new heading and
		// normal vectors.
		return p, h_new, n_new
	}
}

func (c *BetterCube) part2() int {
	dr, dc := 0, 1
	p := Vec3{x: 0, y: c.board.cube_len - 1, z: 0}
	N := Vec3{x: 0, y: 1, z: 0}

	for _, instruction := range c.board.instructions {
		switch i := instruction.(type) {
		case GoLeft:
			dr, dc = -dc, dr // dx, dy = -dy, dx but we are in rows and columns.
		case GoRight:
			dr, dc = dc, -dr
		case GoForward:
			for x := i.x; x > 0; x-- {
				// H = A * (dr dc)
				H := mat.NewVecDense(3, nil)
				A, ok := c.A[N]
				if !ok {
					log.Fatal("failed to find A for", N)
				}
				H.MulVec(A, mat.NewVecDense(2, []float64{float64(dr), float64(dc)}))

				// Overwrite p and N unconditionally. Don't bother trying to stop early.
				p, H, N = c.move(p, H, N)

				// Overwrite dr and dc with new heading.
				dr_dc := mat.NewVecDense(2, nil)
				if err := dr_dc.SolveVec(c.A[N], H); err != nil {
					log.Fatal("failed to solve dr/dc from heading")
				}
				dr = int(math.Round(dr_dc.AtVec(0)))
				dc = int(math.Round(dr_dc.AtVec(1)))
			}
		}
	}

	// Now we just need to get r and c back out of our vectors.
	p2 := p.ToDense()
	n2 := N.ToDense()
	position := c.GetRowCol(&p2, &n2)
	return 1000*(1+position.row) + 4*(1+position.col) + Facing(dr, dc)
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
