package main

import (
	_ "embed"
	"fmt"
	"log"
	"math"
)

//go:embed day23.txt
var puzzle string

func main() {
	fmt.Println("Part 1:", part1(parse(puzzle)))
}

type Point struct {
	x, y int
}

type RuleAndResult struct {
	IsSatisfied bool
	Result      Point
}

// Apparently you don't need to pass a pointer to a map in Go.
// The internal structure of a map has pointers to the data,
// so copying those pointers is already pretty lightweight.
// This saves you the ugly syntax that looks like (*M)[x].
// Passing a map directly into a function does not copy its data.
func (p *Point) next(elves Elves, round int) Point {
	x := p.x
	y := p.y
	_, N := elves[Point{x: x, y: y + 1}]
	_, NE := elves[Point{x: x + 1, y: y + 1}]
	_, NW := elves[Point{x: x - 1, y: y + 1}]
	_, S := elves[Point{x: x, y: y - 1}]
	_, SE := elves[Point{x: x + 1, y: y - 1}]
	_, SW := elves[Point{x: x - 1, y: y - 1}]
	_, W := elves[Point{x: x - 1, y: y}]
	_, E := elves[Point{x: x + 1, y: y}]

	if !N && !NE && !NW && !S && !SE && !SW && !W && !E {
		return *p
	}

	rr := [4]RuleAndResult{
		{
			IsSatisfied: !N && !NE && !NW,
			Result:      Point{x: x, y: y + 1},
		},
		{
			IsSatisfied: !S && !SE && !SW,
			Result:      Point{x: x, y: y - 1},
		},
		{
			IsSatisfied: !W && !NW && !SW,
			Result:      Point{x: x - 1, y: y},
		},
		{
			IsSatisfied: !E && !NE && !SE,
			Result:      Point{x: x + 1, y: y},
		},
	}

	for i := range 4 {
		r := rr[(i+round)%4]
		if r.IsSatisfied {
			return r.Result
		}
	}

	// No rule matched.
	return *p
}

type Elves map[Point]struct{}

func parse(input string) Elves {
	elves := make(Elves)
	row := 0
	col := 0
	for _, char := range input {
		switch char {
		case '#':
			elves[Point{x: col, y: -row}] = struct{}{}
		case '.':
		case '\n':
			row++
			col = 0
			continue
		}
		col++
	}
	return elves
}

func part1(elves Elves) int {
	// 10 rounds for part 1
	for round := range 10 {
		// fmt.Printf("Round %d:\n%s\n\n", round, elves)
		// first half of the round: proposals
		proposal := make(map[Point]int)
		for e := range elves {
			p := e.next(elves, round)
			if n, ok := proposal[p]; ok {
				proposal[p] = n + 1
			} else {
				proposal[p] = 1
			}
		}
		// second half of the round: move iff no conflict
		new_elves := make(Elves)
		for e := range elves {
			p := e.next(elves, round)
			n, ok := proposal[p]
			if !ok {
				log.Fatal("did not find expected proposal")
			}
			var move Point
			if n > 1 {
				// don't move
				move = e
			} else {
				// move to next to proposed position
				move = p
			}
			new_elves[move] = struct{}{}
		}
		if len(elves) != len(new_elves) {
			log.Fatal("we lost or gained an elf")
		}
		// commit proposals
		elves = new_elves
	}

	fmt.Println(elves)
	min_x, max_x, min_y, max_y := elves.bounds()
	return (1+max_x-min_x)*(1+max_y-min_y) - len(elves)
}

func (elves Elves) bounds() (min_x, max_x, min_y, max_y int) {
	min_x = math.MaxInt
	min_y = math.MaxInt
	max_x = math.MinInt
	max_y = math.MinInt
	for p := range elves {
		min_x = min(min_x, p.x)
		min_y = min(min_y, p.y)
		max_x = max(max_x, p.x)
		max_y = max(max_y, p.y)
	}
	return
}

func (elves Elves) String() string {
	min_x, max_x, min_y, max_y := elves.bounds()
	s := ""
	for y := max_y; y >= min_y; y-- {
		for x := min_x; x <= max_x; x++ {
			if _, ok := elves[Point{x, y}]; ok {
				s += "#"
			} else {
				s += "."
			}
		}
		if y > min_y {
			s += "\n"
		}
	}
	return s
}
