package main

import (
	_ "embed"
	"fmt"
)

//go:embed day24.txt
var puzzle string

func main() {
	blizzards := parse(puzzle)
	part1 := bfs(blizzards, start(), end())
	part2 := part1 + bfs(blizzards, end(), start())
	part2 += bfs(blizzards, start(), end())
	fmt.Println("Part 1:", part1)
	fmt.Println("Part 2:", part2)
}

type Point struct {
	x, y int
}

func (p Point) moves() []Point {
	adj := make([]Point, 0)
	// We always have the option to stay still.
	adj = append(adj, p)

	// special case: start
	if p == start() {
		adj = append(adj, Point{x: p.x, y: p.y - 1})
		// We can't move left, right, or up from the start.
		return adj
	}

	// special case: end
	if p == end() {
		adj = append(adj, Point{x: p.x, y: p.y + 1})
		// We can't move left, right, or down from the end.
		return adj
	}

	// go left unless we are at the left wall
	if p.x > 1 {
		adj = append(adj, Point{x: p.x - 1, y: p.y})
	}
	// go right unless we are at the right wall
	if p.x < width-2 {
		adj = append(adj, Point{x: p.x + 1, y: p.y})
	}
	// go down unless we are at the bottom wall
	if p.y > 1 {
		adj = append(adj, Point{x: p.x, y: p.y - 1})
	}
	// go up unless we are at the top wall
	if p.y < height-2 {
		adj = append(adj, Point{x: p.x, y: p.y + 1})
	}
	// special case: goal
	if p.x == width-2 && p.y == 1 {
		adj = append(adj, end())
	}
	// special case: entrance
	if p.x == 1 && p.y == height-2 {
		adj = append(adj, start())
	}
	return adj
}

type Blizzard struct {
	Point
	dx, dy int
}

var height, width int

func start() Point {
	return Point{x: 1, y: height - 1}
}

func end() Point {
	return Point{x: width - 2, y: 0}
}

func (b *Blizzard) tick() {
	b.x = b.x + b.dx
	if b.x == 0 {
		b.x = width - 2
	}
	if b.x == width-1 {
		b.x = 1
	}

	b.y = b.y + b.dy
	if b.y == 0 {
		b.y = height - 2
	}
	if b.y == height-1 {
		b.y = 1
	}
}

func parse(input string) []Blizzard {
	var row, col int
	blizzards := make([]Blizzard, 0)
	for _, c := range input {
		switch c {
		case '.', '#':
		case '\r':
			continue
		case '\n':
			width = max(width, col)
			col = 0
			row++
			continue
		case '>':
			blizzards = append(blizzards, Blizzard{
				x: col, y: -row, dx: 1, dy: 0,
			})
		case '<':
			blizzards = append(blizzards, Blizzard{
				x: col, y: -row, dx: -1, dy: 0,
			})
		case '^':
			blizzards = append(blizzards, Blizzard{
				x: col, y: -row, dx: 0, dy: 1,
			})
		case 'v':
			blizzards = append(blizzards, Blizzard{
				x: col, y: -row, dx: 0, dy: -1,
			})
		}
		col++
	}
	height = row
	// Let's just make all those y values positive.
	for i := range blizzards {
		blizzards[i].y += height - 1
	}
	return blizzards
}

func bfs(blizzards []Blizzard, start, goal Point) int {
	positions := make(map[Point]struct{})
	positions[start] = struct{}{}
	positions_with_blizzards := make(map[Point]int)

	for round := 0; ; round++ {
		// Did we make it to the goal?
		if _, ok := positions[goal]; ok {
			return round
		}

		// Update all the blizzards and store their positions in a map.
		clear(positions_with_blizzards)
		for i := range blizzards {
			blizzards[i].tick()
			positions_with_blizzards[blizzards[i].Point]++
		}

		// Ok, now go through all the reachable positions and determine where we can go.
		next_positions := make(map[Point]struct{})
		for p := range positions {
			for _, move := range p.moves() {
				if _, ok := positions_with_blizzards[move]; !ok {
					// The adjacent tile has no blizzard.
					next_positions[move] = struct{}{}
				}
			}
		}
		// Commit the next round of positions and continue.
		positions = next_positions
	}
}
