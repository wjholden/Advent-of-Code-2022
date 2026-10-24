package main

import (
	_ "embed"
	"flag"
	"fmt"
	"log"
	"slices"
	"strings"
)

//go:embed day17.txt
var puzzle string

//go:embed example.txt
var example string

type Point struct {
	x int
	y int
}

type Shape []Point

var max_y []int

var positions map[Point]rune

var verbose bool
var show_strudel bool
var show_last bool

type State struct {
	y1       int
	y2       int
	y3       int
	y4       int
	y5       int
	y6       int
	y7       int
	shape_id int
	dir_id   int
}

type Tower struct {
	count  int
	height int
}

func main() {
	verbose_flag := flag.Bool("verbose", false, "Show additional debugging information")
	strudel_flag := flag.Bool("strudel", false, "Create musical score for https://strudel.cc")
	showlast_flag := flag.Bool("show_last", false, "Draw the tower at the end")
	example_flag := flag.Bool("example", false, "Use the example input")
	flag.Parse()
	verbose = *verbose_flag
	show_strudel = *strudel_flag
	show_last = *showlast_flag

	if *example_flag {
		puzzle = example
	}

	fmt.Println("Part 1:", run(2023, false))

	if *strudel_flag {
		strudel()
	}

	fmt.Println("Part 2:", run(1000000000000, true))
}

func run(limit int, part2 bool) int {
	positions = map[Point]rune{}

	var cycle_length int
	var cycle_height int
	var dh int

	state_history := make(map[State]Tower)

	max_y = make([]int, 7)
	sg := ShapeGenerator{}
	jg := JetGenerator{}

	// outer:
	for i := 1; i < limit; i++ {
		shape, shape_index := sg.next()
		if verbose {
			fmt.Println("New shape", shape_index)
			shape.print()
		}

		// Cycle length gets initialized to zero. Once we find it, we don't need to
		// bother with all this anymore.
		if part2 && cycle_length == 0 {
			// check if we found the cycle. If we did, add enough that we just need to finish everything.
			yn := y_normalized()
			_, dir := jg.peek()
			state := State{
				y1: yn[0], y2: yn[1], y3: yn[2], y4: yn[3], y5: yn[4], y6: yn[5], y7: yn[6], shape_id: shape_index, dir_id: dir,
			}
			if previous_i, ok := state_history[state]; ok {
				if verbose {
					fmt.Println("Found cycle in", len(state_history), " with ", state)
				}
				cycle_height = slices.Max(max_y) - previous_i.height
				cycle_length = i - previous_i.count

				cycles := (1000000000000 - (i - 1)) / cycle_length
				dh = cycles * cycle_height
				i = (i - 1) + cycles*cycle_length
			} else {
				state_history[state] = Tower{height: slices.Max(max_y), count: i}
				if verbose {
					fmt.Println(state, state_history[state])
				}
			}
		}
	inner:
		for {
			shape.check()
			if verbose {
				fmt.Println(shape)
				shape.print()
			}
			dx := jg.next()
			if shape.can_blow(dx) {
				shape.blow(dx)
			}
			if verbose {
				shape.print()
			}
			if shape.can_drop() {
				shape.drop()
			} else {
				shape.stop()
				if verbose {
					fmt.Println("New maximum y positions", max_y)
				}
				// if i == 2022 {
				// 	part1 = slices.Max(max_y)
				// }
				break inner
			}
		}
	}

	if !part2 && show_last {
		last, _ := sg.next()
		last.print()
	}

	return slices.Max(max_y) + dh // dh is 0 if part1.
}

var letters = []string{"c", "d", "e", "f", "g", "a", "b"}

func row_chord(y int) string {
	current := make([]string, 0)
	for x := range 7 {
		if _, ok := positions[Point{x, y}]; ok {
			current = append(current, letters[x])
		}
	}
	if len(current) == 0 {
		return "-"
	} else if len(current) == 1 {
		return current[0]
	} else {
		return fmt.Sprintf("[%s]", strings.Join(current, ","))
	}
}

func strudel() {
	notes := make([]string, 0)
	limit := slices.Max(max_y)
	for y := 1; y <= limit; y++ {
		notes = append(notes, row_chord(y))
	}
	fmt.Printf("note(\"<%s>*8\").sound(\"piano\")\n", strings.Join(notes, " "))
}

func y_normalized() [7]int {
	y_min := slices.Min(max_y)
	var result [7]int
	for i, y := range max_y {
		result[i] = y - y_min
	}
	return result
}

func (s Shape) check() {
	for _, point := range s {
		if point.x < 0 || point.x >= 7 || point.y < 0 {
			log.Fatal("invalid state", s)
		}
	}
}

func (s Shape) height() int {
	y := 0
	for _, point := range s {
		y = max(y, point.y)
	}
	return y
}

func (s Shape) print() {
	for _, point := range s {
		positions[point] = '@'
	}
	for y := s.height(); y > 0; y-- {
		fmt.Print("|")
		for x := range 7 {
			p := Point{x, y}
			if c, ok := positions[p]; ok {
				fmt.Printf("%c", c)
			} else {
				fmt.Print(".")
			}
		}
		if show_strudel {
			fmt.Println("|", row_chord(y))
		} else {
			fmt.Println("|")
		}
	}
	fmt.Println("+-------+")
	if show_strudel {
		fmt.Println("|cdefgab|")
	}
	for _, point := range s {
		delete(positions, point)
	}
}

func (s Shape) can_drop() bool {
	if verbose {
		fmt.Print("Can it drop? ")
	}
	for _, point := range s {
		x := point.x
		y := point.y
		if _, ok := positions[Point{x: x, y: y - 1}]; ok {
			if verbose {
				fmt.Println("No, it can't drop into this filled space.")
			}
			return false
		}
		if y == 1 {
			if verbose {
				fmt.Println("No, it can't fall through the floor.")
			}
			return false
		}
	}
	if verbose {
		fmt.Println("Yes, it can drop.")
	}
	return true
}

func (s Shape) drop() {
	for i := range s {
		s[i].y--
	}
}

func (s Shape) can_blow(dx int) bool {
	if verbose {
		if dx > 0 {
			fmt.Print("Can it blow right? ")
		} else {
			fmt.Print("Can it blow left? ")
		}
	}
	for _, point := range s {
		if point.x == 0 && dx < 0 {
			if verbose {
				fmt.Println("No, it can't blow left of the wall.")
			}
			return false
		}
		if point.x == 6 && dx > 0 {
			if verbose {
				fmt.Println("No, it can't blow right of the wall.")
			}
			return false
		}
		if _, ok := positions[Point{x: point.x + dx, y: point.y}]; ok {
			if verbose {
				fmt.Println("No, it can't blow because of another object.")
			}
			return false
		}
	}
	if verbose {
		fmt.Println("Yes, it can blow.")
	}
	return true
}

func (s Shape) blow(dx int) {
	for i := range s {
		s[i].x += dx
	}
}

func (s Shape) stop() {
	for _, point := range s {
		positions[point] = '#'
		x := point.x
		y := point.y
		max_y[x] = max(max_y[x], y)
	}
}

type ShapeGenerator struct {
	i int
}

func (sg *ShapeGenerator) next() (Shape, int) {
	x := 2
	y := 4 + slices.Max(max_y)
	index := sg.i % 5
	sg.i++
	switch index {
	case 0:
		// fmt.Println("-")
		return Shape{
			{x, y},
			{x + 1, y},
			{x + 2, y},
			{x + 3, y},
		}, index
	case 1:
		// fmt.Println("+")
		return Shape{
			{x + 1, y + 2},
			{x, y + 1},
			{x + 1, y + 1},
			{x + 2, y + 1},
			{x + 1, y},
		}, index
	case 2:
		// fmt.Println("⅃")
		return Shape{
			{x, y},
			{x + 1, y},
			{x + 2, y},
			{x + 2, y + 1},
			{x + 2, y + 2},
		}, index
	case 3:
		// fmt.Println("|")
		return Shape{
			{x, y},
			{x, y + 1},
			{x, y + 2},
			{x, y + 3},
		}, index
	case 4:
		// fmt.Println("■")
		return Shape{
			{x, y},
			{x + 1, y},
			{x + 1, y + 1},
			{x, y + 1},
		}, index
	}
	return nil, -1 // unreachable
}

type JetGenerator struct {
	i int
}

func (j *JetGenerator) next() int {
	n := len(puzzle)
	defer func() { j.i++ }()
	for {
		switch puzzle[j.i%n] {
		case '<':
			return -1
		case '>':
			return +1
		default:
			j.i++
		}
	}
}

func (j *JetGenerator) peek() (int, int) {
	n := len(puzzle)
	k := j.i
	for {
		switch puzzle[k%n] {
		case '<':
			return -1, k % n
		case '>':
			return +1, k % n
		default:
			k++
		}
	}
}
