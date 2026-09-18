// https://go.dev/play/p/FTrnAcK6BY1
package main

import (
	_ "embed"
	"fmt"
	"math"
	"strconv"
	"strings"
)

type move struct {
	direction string
	distance  int
}

type point struct {
	x, y int
}

//go:embed day09.txt
var puzzle string

func main() {
	moves := parse(strings.TrimSpace(puzzle))
	var knots [10]point
	var H *point
	H = &knots[0]
	positionHistory1 := make(map[point]bool)
	positionHistory1[point{0, 0}] = true
	positionHistory2 := make(map[point]bool)
	positionHistory2[point{0, 0}] = true
	for _, move := range moves {
		d := move.distance
		for d > 0 {
			switch move.direction {
			case "U":
				H.y += 1
			case "D":
				H.y -= 1
			case "L":
				H.x -= 1
			case "R":
				H.x += 1
			}
			//fmt.Printf("H is at (%d, %d)\n", H.x, H.y)

			for i := 0; i < 9; i++ {
				h := &knots[i]
				t := &knots[i+1]
				if !adjacent(h, t) {
					moveTail(h, t)
					if i == 0 {
						positionHistory1[point{t.x, t.y}] = true
					}
					if i == 8 {
						positionHistory2[point{t.x, t.y}] = true
					}
				}
			}
			d--
		}
	}
	fmt.Printf("Part 1: %d\n", len(positionHistory1))
	fmt.Printf("Part 1: %d\n", len(positionHistory2))
	//printTailPath(positionHistory)
}

func printTailPath(history map[point]bool) {
	// find bounds
	var xmax, ymax int
	for p, _ := range history {
		if p.x > xmax {
			xmax = p.x
		}
		if p.y > ymax {
			ymax = p.y
		}
	}
	for y := ymax; y >= 0; y-- {
		for x := 0; x <= xmax; x++ {
			if history[point{x, y}] {
				fmt.Print("#")
			} else {
				fmt.Print(".")
			}
		}
		fmt.Println()
	}
}

func adjacent(p1, p2 *point) bool {
	dx := math.Abs(float64(p1.x) - float64(p2.x))
	dy := math.Abs(float64(p1.y) - float64(p2.y))
	return dx <= 1 && dy <= 1
}

func moveTail(H, T *point) {
	var dx, dy int
	if H.x < T.x {
		dx = -1
	} else if H.x > T.x {
		dx = +1
	}

	if H.y < T.y {
		dy = -1
	} else if H.y > T.y {
		dy = +1
	}

	T.x += dx
	T.y += dy
}

func parse(s string) []move {
	var moves []move
	for _, line := range strings.Split(s, "\n") {
		tokens := strings.Split(line, " ")
		if len(tokens) != 2 {
			panic("input line is too short: " + line)
		}
		d, err := strconv.Atoi(tokens[1])
		if err != nil {
			panic("input line does not contain numeric distance: " + line)
		}
		moves = append(moves, move{tokens[0], d})
	}
	return moves
}

const sample = `R 4
U 4
L 3
D 1
R 4
D 1
L 5
R 2`

const sample2 = `R 5
U 8
L 8
D 3
R 17
D 10
L 25
U 20`
