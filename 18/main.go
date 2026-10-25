package main

import (
	_ "embed"
	"fmt"
	"log"
	"strconv"
	"strings"
)

//go:embed day18.txt
var puzzle string

type Cube struct {
	x int
	y int
	z int
}

type D struct {
	dx int
	dy int
	dz int
}

var deltas []D = []D{
	{dx: 0, dy: 0, dz: 1}, {dx: 0, dy: 0, dz: -1},
	{dx: 0, dy: 1, dz: 0}, {dx: 0, dy: -1, dz: 0},
	{dx: 1, dy: 0, dz: 0}, {dx: -1, dy: 0, dz: 0},
}

// Using a buffer here is intended to save on allocations.
// I can't tell that it actually helps.
func (cube *Cube) neighbors(buf *[]Cube) *[]Cube {
	for i, delta := range deltas {
		(*buf)[i].x = cube.x + delta.dx
		(*buf)[i].y = cube.y + delta.dy
		(*buf)[i].z = cube.z + delta.dz
	}
	return buf
}

func (cube *Cube) in_bounds(max_x int, max_y int, max_z int) bool {
	return -1 <= cube.x && cube.x <= max_x+1 &&
		-1 <= cube.y && cube.y <= max_y+1 &&
		-1 <= cube.z && cube.z <= max_z+1
}

func main() {
	cubes := parse(puzzle)
	part1, part2 := solve(&cubes)
	fmt.Println("Part 1:", part1)
	fmt.Println("Part 2:", part2) // 2556 too low
}

func parse(input string) (cubes []Cube) {
	for line := range strings.Lines(input) {
		xyz := strings.Split(line, ",")
		x, err := strconv.Atoi(xyz[0])
		if err != nil {
			log.Fatal(err)
		}
		y, err := strconv.Atoi(xyz[1])
		if err != nil {
			log.Fatal(err)
		}
		z, err := strconv.Atoi(strings.TrimSpace(xyz[2]))
		if err != nil {
			log.Fatal(err)
		}
		cubes = append(cubes, Cube{x, y, z})
	}
	return
}

func solve(cubes *[]Cube) (int, int) {
	var surface_area int
	volume := map[Cube]struct{}{}
	buf := make([]Cube, 6)

	for _, cube := range *cubes {
		ds := 0
		for _, neighbor := range *cube.neighbors(&buf) {
			if _, ok := volume[neighbor]; ok {
				ds--
			} else {
				ds++
			}
		}
		surface_area += ds
		volume[cube] = struct{}{}
	}

	return surface_area, external_area(&volume)
}

func external_area(volume *map[Cube]struct{}) (external_surface_area int) {
	q := make([]Cube, 0)
	q = append(q, Cube{}) // 0, 0, 0

	discovered := make(map[Cube]struct{})

	var max_x int
	var max_y int
	var max_z int

	for cube := range *volume {
		max_x = max(max_x, cube.x)
		max_y = max(max_y, cube.y)
		max_z = max(max_z, cube.z)
	}

	buf := make([]Cube, 6)

	for {
		if len(q) == 0 {
			return
		}
		current := q[0]
		q = q[1:]

		if _, ok := (*volume)[current]; ok {
			log.Fatal("Should not have entered volume at", current)
		}

		for _, neighbor := range *current.neighbors(&buf) {
			if _, ok := (*volume)[neighbor]; ok {
				external_surface_area++
				continue
			}
			if _, ok := discovered[neighbor]; ok {
				continue
			}
			discovered[neighbor] = struct{}{}
			if neighbor.in_bounds(max_x, max_y, max_z) {
				q = append(q, neighbor)
			}
		}
	}
}
