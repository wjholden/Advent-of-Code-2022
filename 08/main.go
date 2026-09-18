// https://go.dev/play/p/rnFmT1BqpZg
package main

import (
	_ "embed"
	"fmt"
	"strconv"
	"strings"
)

const sample = `30373
25512
65332
33549
35390`

//go:embed day08.txt
var puzzle string

func main() {
	M, rows, cols := parseInput(puzzle)
	//fmt.Println(M, rows, cols)
	part1 := 0
	part2 := 0
	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			//fmt.Printf("M[%d][%d]=%d: %t %t %t %t\n", r, c, M[r][c], up(M, r, c), down(M, r, c), left(M, r, c), right(M, r, c))
			var tallest [4]bool
			var distance [4]int
			tallest[0], distance[0] = up(M, r, c)
			tallest[1], distance[1] = down(M, r, c)
			tallest[2], distance[2] = left(M, r, c)
			tallest[3], distance[3] = right(M, r, c)
			if tallest[0] || tallest[1] || tallest[2] || tallest[3] {
				part1++
			}
			scenicScore := distance[0] * distance[1] * distance[2] * distance[3]
			//fmt.Printf("M[%d][%d]=%d: score = %d\n", r, c, M[r][c], scenicScore)
			part2 = max(part2, scenicScore)
		}
	}
	fmt.Printf("Part 1: %d\n", part1)
	fmt.Printf("Part 2: %d\n", part2)
}

func parseInput(x string) ([][]int, int, int) {
	inputSet := x
	cols := strings.Index(inputSet, "\n")
	inputSet = strings.ReplaceAll(inputSet, "\n", "")
	rows := (len(inputSet)) / cols
	X := strings.Split(inputSet, "")

	M := make([][]int, rows)
	for r := 0; r < rows; r++ {
		M[r] = make([]int, cols)
		for c := 0; c < cols; c++ {
			pos := cols*r + c
			y, err := strconv.Atoi(X[pos])
			if err != nil {
				panic(fmt.Sprintf("failed to parse %s at position %d", X[pos], pos))
			}
			M[r][c] = y
		}
	}

	return M, rows, cols
}

func max(a, b int) int {
	if a >= b {
		return a
	}
	return b
}

// You don't actually have to go all the way to the edge in most cases.
// You just need to continue in that direction until you've found a taller
// tree. Once you've found a taller tree, return false.
// Otherwise, continue until you reach the edge and return true.
func tallest(M [][]int, row, col, dr, dc int) (bool, int) {
	var viewingDistance int
	x := M[row][col]
	r := row + dr
	c := col + dc
	for 0 <= r && r < len(M) && 0 <= c && c < len(M[0]) {
		viewingDistance++
		if M[r][c] >= x {
			return false, viewingDistance
		}
		r += dr
		c += dc
	}
	return true, viewingDistance
}

func up(M [][]int, row, col int) (bool, int) {
	return tallest(M, row, col, -1, 0)
}

func down(M [][]int, row, col int) (bool, int) {
	return tallest(M, row, col, +1, 0)
}

func left(M [][]int, row, col int) (bool, int) {
	return tallest(M, row, col, 0, -1)
}

func right(M [][]int, row, col int) (bool, int) {
	return tallest(M, row, col, 0, +1)
}
