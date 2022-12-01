package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
)

func main() {
	// https://stackoverflow.com/questions/6141604/go-readline-string
	f, err := os.Open("input.txt")
	if err != nil {
		panic(err)
	}
	// https://go.dev/tour/flowcontrol/12
	defer f.Close()

	scanner := bufio.NewScanner(f)
	// https://go.dev/tour/moretypes/15
	elves := make([]int, 1)
	max := 0
	for scanner.Scan() {
		if line := scanner.Text(); len(line) > 0 {
			calories, err := strconv.Atoi(line)
			if err != nil {
				panic(err)
			}
			elves[len(elves)-1] += calories
			if elves[len(elves)-1] > max {
				max = elves[len(elves)-1]
			}
		} else {
			elves = append(elves, 0)
		}
	}

	fmt.Println("Part 1", max)

	// Sort the array and take the last 3 for part 2.
	sort.Sort(sort.Reverse(sort.IntSlice(elves)))

	// Does Go not have a `sum()` function?
	total := 0
	for _, element := range elves[:3] {
		total += element
	}
	fmt.Println(elves[:3], total)
}
