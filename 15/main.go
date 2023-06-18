// You can edit this code!
// Click here and start typing.
package main

import (
	"fmt"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
)

type Point struct {
	x, y int
}

var sensors = map[Point]Point{}
var beacons = map[Point]bool{}

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

	//show(0, 20, 0, 20)

	// We think the solution is going to be on some corner away from sensors.
	corners := make([]Point, 0)
	for sensor, beacon := range sensors {
		d := distance(sensor, beacon)
		sx, sy := sensor.x, sensor.y
		corners = append(corners, Point{sx - d - 1, sy}, Point{sx + d + 1, sy}, Point{sx, sy - d - 1}, Point{sx, sy + d + 1})
	}
	fmt.Println(2*y, corners)
corner:
	for _, corner := range corners {
		if 0 <= corner.x && corner.x <= 2*y && 0 <= corner.y && corner.y <= 2*y {
			for sensor, beacon := range sensors {
				// If this "corner" position is inside of a diamond, then it is not a candidate solution.
				if distance(corner, sensor) <= distance(sensor, beacon) {
					continue corner
				}
			}
			fmt.Println("Part 2:", 4000000*corner.x+corner.y, corner)
		}
	}
	//svg(20)
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

func show(xmin, xmax, ymin, ymax int) {
	for y := ymin; y <= ymax; y++ {
		for x := xmin; x <= xmax; x++ {
			p := Point{x, y}
			if _, ok1 := sensors[p]; ok1 {
				fmt.Print("S")
				continue
			}
			if _, ok2 := beacons[p]; ok2 {
				fmt.Print("B")
				continue
			}
			if cannotBe(x, y) {
				fmt.Print("#")
			} else {
				fmt.Print(".")
			}
		}
		fmt.Println()
	}
}

func svg(size int) {
	fmt.Printf("<!DOCTYPE html>\n<html>\n<body>\n<svg height='%d' width='%d'>\n", size, size)
	for sensor, beacon := range sensors {
		sx, sy := sensor.x, sensor.y
		d := distance(sensor, beacon)
		fmt.Printf("<polygon points='%d,%d %d,%d %d,%d %d,%d' style='fill: rgb(%d, %d, %d)' />\n",
			sx-d, sy, sx, sy+d, sx+d, sy, sx, sy-d, rand.Intn(255), rand.Intn(255), rand.Intn(255))
	}
	fmt.Println("</svg>\n</body>\n</html>")
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

const puzzleY = 2000000
const puzzle = `Sensor at x=2885528, y=2847539: closest beacon is at x=2966570, y=2470834
Sensor at x=2224704, y=1992385: closest beacon is at x=2018927, y=2000000
Sensor at x=3829144, y=1633329: closest beacon is at x=2966570, y=2470834
Sensor at x=43913, y=426799: closest beacon is at x=152363, y=369618
Sensor at x=2257417, y=2118161: closest beacon is at x=2386559, y=2090397
Sensor at x=8318, y=3994839: closest beacon is at x=-266803, y=2440278
Sensor at x=69961, y=586273: closest beacon is at x=152363, y=369618
Sensor at x=3931562, y=3361721: closest beacon is at x=3580400, y=3200980
Sensor at x=476279, y=3079924: closest beacon is at x=-266803, y=2440278
Sensor at x=2719185, y=2361091: closest beacon is at x=2966570, y=2470834
Sensor at x=2533382, y=3320911: closest beacon is at x=2260632, y=3415930
Sensor at x=3112735, y=3334946: closest beacon is at x=3580400, y=3200980
Sensor at x=1842258, y=3998928: closest beacon is at x=2260632, y=3415930
Sensor at x=3712771, y=3760832: closest beacon is at x=3580400, y=3200980
Sensor at x=1500246, y=2684955: closest beacon is at x=2018927, y=2000000
Sensor at x=3589321, y=142859: closest beacon is at x=4547643, y=-589891
Sensor at x=1754684, y=2330721: closest beacon is at x=2018927, y=2000000
Sensor at x=2476631, y=3679883: closest beacon is at x=2260632, y=3415930
Sensor at x=27333, y=274008: closest beacon is at x=152363, y=369618
Sensor at x=158732, y=2405833: closest beacon is at x=-266803, y=2440278
Sensor at x=2955669, y=3976939: closest beacon is at x=3035522, y=4959118
Sensor at x=1744196, y=13645: closest beacon is at x=152363, y=369618
Sensor at x=981165, y=1363480: closest beacon is at x=2018927, y=2000000
Sensor at x=2612279, y=2151377: closest beacon is at x=2386559, y=2090397
Sensor at x=3897, y=2076376: closest beacon is at x=-266803, y=2440278
Sensor at x=2108479, y=1928318: closest beacon is at x=2018927, y=2000000
Sensor at x=1913043, y=3017841: closest beacon is at x=2260632, y=3415930
Sensor at x=2446778, y=785075: closest beacon is at x=2386559, y=2090397
Sensor at x=2385258, y=2774943: closest beacon is at x=2386559, y=2090397
Sensor at x=3337656, y=2916144: closest beacon is at x=3580400, y=3200980
Sensor at x=380595, y=66906: closest beacon is at x=152363, y=369618
Sensor at x=1593628, y=3408455: closest beacon is at x=2260632, y=3415930`
