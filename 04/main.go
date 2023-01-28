package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
)

func main() {
	// read input from file
	f, err := os.Open("input.txt")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	// For variety, let's read from stdin instead of a file.
	// (On second thought, this makes it harder to run them all in a batch.
	// I changed this to the same file-oriented model in 2023).
	// scanner := bufio.NewScanner(os.Stdin)
	scanner := bufio.NewScanner(f)
	lines := make([]string, 0)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	// Parse each line as 4 integers and build a matrix
	regex := regexp.MustCompile("-|,")
	M := make([][]int, 0)
	for _, line := range lines {
		line_s := regex.Split(line, -1)
		line_i := make([]int, 0)
		for _, s := range line_s {
			x, err := strconv.Atoi(s)
			if err != nil {
				panic(err)
			}
			line_i = append(line_i, x)
		}
		M = append(M, line_i)
	}

	// Find overlapping pairs for part 1
	overlap1 := 0
	overlap2 := 0
	for _, row := range M {
		start1, start2, end1, end2 := row[0], row[2], row[1], row[3]
		if (start1 <= start2 && end2 <= end1) || (start2 <= start1 && end1 <= end2) {
			overlap1++
		}
		// Ohhh it's tricky!
		if (start1 <= start2 && start2 <= end1) || (start2 <= end1 && end1 <= end2) || (start2 <= start1 && start1 <= end2) || (start1 <= end2 && end2 <= end1) {
			overlap2++
		}
	}

	fmt.Println("Part 1:", overlap1)
	fmt.Println("Part 2:", overlap2) // 715 is too low
}
