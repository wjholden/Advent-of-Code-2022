package main

import (
	_ "embed"
	"fmt"
	"log"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

type Point struct {
	x, y int
}

var sensors = map[Point]Point{}
var beacons = map[Point]bool{}

const puzzleY = 2000000

//go:embed day15.txt
var puzzle string

func distance(a, b Point) int {
	dx := a.x - b.x
	dy := a.y - b.y
	if dx < 0 {
		dx *= -1
	}
	if dy < 0 {
		dy *= -1
	}
	return dx + dy
}

func parse(s string) {
	pattern := `Sensor at x=(?P<sx>-?\d+), y=(?P<sy>-?\d+): closest beacon is at x=(?P<bx>-?\d+), y=(?P<by>-?\d+)`
	r := regexp.MustCompile(pattern)
	for _, line := range strings.Split(strings.TrimSpace(s), "\n") {
		match := r.FindStringSubmatch(line)[1:]
		sx, _ := strconv.Atoi(match[0])
		sy, _ := strconv.Atoi(match[1])
		bx, _ := strconv.Atoi(match[2])
		by, _ := strconv.Atoi(match[3])
		sensor := Point{sx, sy}
		beacon := Point{bx, by}
		sensors[sensor] = beacon
		beacons[beacon] = true
	}
}

func main() {
	parse(puzzle)
	y := puzzleY

	distances := make([]int, 0)
	for sensor, beacon := range sensors {
		d := distance(sensor, beacon)
		distances = append(distances, d)
	}

	xmin := 1 << 30
	xmax := -xmin
	for sensor, beacon := range sensors {
		d := distance(sensor, beacon)
		xmin = min(xmin, sensor.x-d)
		xmax = max(xmax, sensor.x+d)
	}
	var part1 int
	for x := xmin; x <= xmax; x++ {
		if cannotBe(x, y) {
			part1++
		}
	}
	fmt.Println("Part 1:", part1)

	// Meh. I'm OK with using a constraint solver for this one.
	cmd := exec.Command("python", "part2.py")
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("Command failed: %v", err)
	}
	fmt.Print(string(output))

}

func cannotBe(x, y int) bool {
	point := Point{x, y}
	for sensor, beacon := range sensors {
		if point != beacon && distance(point, sensor) <= distance(sensor, beacon) {
			return true
		}
	}
	return false
}

func min(a, b int) int {
	if a <= b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a >= b {
		return a
	}
	return b
}

const sampleY = 10
const sample = `Sensor at x=2, y=18: closest beacon is at x=-2, y=15
Sensor at x=9, y=16: closest beacon is at x=10, y=16
Sensor at x=13, y=2: closest beacon is at x=15, y=3
Sensor at x=12, y=14: closest beacon is at x=10, y=16
Sensor at x=10, y=20: closest beacon is at x=10, y=16
Sensor at x=14, y=17: closest beacon is at x=10, y=16
Sensor at x=8, y=7: closest beacon is at x=2, y=10
Sensor at x=2, y=0: closest beacon is at x=2, y=10
Sensor at x=0, y=11: closest beacon is at x=2, y=10
Sensor at x=20, y=14: closest beacon is at x=25, y=17
Sensor at x=17, y=20: closest beacon is at x=21, y=22
Sensor at x=16, y=7: closest beacon is at x=15, y=3
Sensor at x=14, y=3: closest beacon is at x=15, y=3
Sensor at x=20, y=1: closest beacon is at x=15, y=3`
