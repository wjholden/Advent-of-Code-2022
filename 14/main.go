// Day 14 themes: mutable state, if/elseif chains, large problems, coordinate systems, paths
package main

import (
	_ "embed"
	"fmt"
	"strconv"
	"strings"
	"sync"
)

var minX = 500
var minY = 0
var maxX = 500
var maxY = 0

//go:embed day14.txt
var puzzle string

type Point struct {
	x int
	y int
}

var cave map[Point]string = map[Point]string{}

func getPoints(s string) []Point {
	points := make([]Point, 1+strings.Count(s, " -> "))
	for i, point := range strings.Split(s, " -> ") {
		xy := strings.Split(point, ",")
		x, _ := strconv.Atoi(xy[0])
		y, _ := strconv.Atoi(xy[1])
		if x < minX {
			minX = x
		}
		if x > maxX {
			maxX = x
		}
		if y > maxY {
			maxY = y
		}
		points[i] = Point{x, y}
	}
	return points
}

func getPaths(s string) [][]Point {
	s = strings.TrimSpace(s)
	paths := make([][]Point, 1+strings.Count(s, "\n"))
	for i, path := range strings.Split(s, "\n") {
		paths[i] = getPoints(path)
	}
	return paths
}

func makeRocks(path []Point, ch chan<- Point, wg *sync.WaitGroup) {
	defer (*wg).Done()
	x := path[0].x
	y := path[0].y
	for i := 0; i < len(path)-1; i++ {
		nextPoint := path[i+1]
		// Hey, it's our old friend DeMorgan!
		// Stop when both x1==x2 && y1==y2.
		for x != nextPoint.x || y != nextPoint.y {
			ch <- Point{x, y}
			switch {
			case x < nextPoint.x:
				x++
			case x > nextPoint.x:
				x--
			case y < nextPoint.y:
				y++
			case y > nextPoint.y:
				y--
			}
		}
	}
	// also do the very last one
	ch <- Point{x, y}
}

func makeCave(ch <-chan Point) {
	for point := range ch {
		if _, ok := cave[point]; ok {
			_ = ok
			//fmt.Println("hmm, we already knew about", point)
		}
		cave[point] = "#"
	}
}

func showCave() {
	for y := minY; y <= maxY; y++ {
		for x := minX; x <= maxX; x++ {
			mark, ok := cave[Point{x, y}]
			if ok {
				fmt.Print(mark)
			} else {
				fmt.Print(".")
			}
		}
		fmt.Println()
	}
}

func fillCave(part1 bool) bool {
	x := 500
	y := 0

	// If the starting point (500,0) is filled with sand then the cave is completely filled.
	if caveHas(Point{x, y}, 0, 0) {
		return false
	}

	for {
		position := Point{x, y}
		if part1 && y > maxY {
			// now the cave is so full that sand is falling off into the abyss
			return false
		} else if !caveHasDown(position) {
			y++
		} else if !caveHasDownLeft(position) {
			x--
			y++
		} else if !caveHasDownRight(position) {
			x++
			y++
		} else {
			cave[position] = "O"
			if x < minX {
				minX = x
			}
			if maxX < x {
				maxX = x
			}
			return true
		}
	}
}

func caveHas(point Point, dx, dy int) bool {
	if point.y == maxY+1 {
		return true
	}

	_, ok := cave[Point{point.x + dx, point.y + dy}]
	return ok
}

func caveHasDown(point Point) bool {
	return caveHas(point, 0, +1)
}

func caveHasDownLeft(point Point) bool {
	return caveHas(point, -1, +1)
}

func caveHasDownRight(point Point) bool {
	return caveHas(point, +1, +1)
}

func main() {
	paths := getPaths(puzzle)

	ch := make(chan Point)
	go makeCave(ch)

	// There's really no reason to do this in parallel. I just wanted to.
	var wg sync.WaitGroup
	for _, path := range paths {
		wg.Add(1)
		go makeRocks(path, ch, &wg)
	}
	wg.Wait()
	close(ch)

	//showCave()
	//fmt.Println()

	var part1 int
	for fillCave(true) {
		//showCave()
		//fmt.Println()
		part1++
	}
	showCave()
	fmt.Println("Part 1:", part1)

	var part2 int
	for fillCave(false) {
		part2++
	}
	maxY += 2
	//showCave()
	fmt.Println("Part 2:", part1+part2)
}

const sample = `
498,4 -> 498,6 -> 496,6
503,4 -> 502,4 -> 502,9 -> 494,9
`
