package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func main() {
	// read lines of input
	//scanner := bufio.NewScanner(os.Stdin)
	// read input from file
	f, err := os.Open("input.txt")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	lines := make([]string, 0)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	// parse lines of input
	stacks_input := make([][]string, 0)
	moves := make([][]int, 0)
	stacksComplete := false
	for _, line := range lines {
		if len(line) == 0 {
			stacksComplete = true
		} else if !stacksComplete {
			stacks_input = append(stacks_input, parseStacks(line))
		} else {
			moves = append(moves, parseInstruction(line))
		}
	}
	stacks_input = stacks_input[:len(stacks_input)-1]

	// actually push lines of input to real stacks
	stacks1 := make([][]string, len(stacks_input[0]))
	// iterate over input in reverse because prepending is difficult
	for i := len(stacks_input) - 1; i >= 0; i-- {
		for j, value := range stacks_input[i] {
			if value != " " {
				stacks1[j] = append(stacks1[j], value)
			}
		}
	}

	stacks2 := make([][]string, len(stacks1))
	// Is this a shallow copy?
	//copy(stacks2, stacks1)
	for i := range stacks1 {
		// Apparently it is! Mutating the inner arrays in stacks2 would
		// also cause changes in stacks1.
		stacks2[i] = append(stacks2[i], stacks1[i]...)
	}

	// execute the instructions
	for _, instruction := range moves {
		cnt := instruction[0]
		src := instruction[1]
		dst := instruction[2]

		// part 1 crane logic
		for i := 0; i < cnt; i++ {
			// pop
			x := stacks1[src][len(stacks1[src])-1]
			stacks1[src] = stacks1[src][:len(stacks1[src])-1]
			// push
			stacks1[dst] = append(stacks1[dst], x)
		}

		// part 2 crane logic is a little easier
		{
			y := stacks2[src][len(stacks2[src])-cnt:]
			stacks2[src] = stacks2[src][:len(stacks2[src])-cnt]
			stacks2[dst] = append(stacks2[dst], y...)
		}
	}

	fmt.Print("Part 1: ")
	for _, s1 := range stacks1 {
		fmt.Print(s1[len(s1)-1])
	}
	fmt.Println()

	fmt.Print("Part 2: ")
	for _, s2 := range stacks2 {
		fmt.Print(s2[len(s2)-1])
	}
	fmt.Println()
}

func parseStacks(line string) []string {
	A := make([]string, 0)
	for i, s := range strings.Split(line, "") {
		if i%4 == 1 {
			A = append(A, s)
		}
	}
	return A
}

func parseInstruction(line string) []int {
	// Lesson learned: the + has to go *inside* the parenthesis.
	re := regexp.MustCompile(`move (\d+) from (\d+) to (\d+)`)
	moves := make([]int, 3)
	matches := re.FindAllStringSubmatch(line, -1)[0]
	moves[0], _ = strconv.Atoi(matches[1])
	moves[1], _ = strconv.Atoi(matches[2])
	moves[2], _ = strconv.Atoi(matches[3])
	// change to 0-indexing
	moves[1]--
	moves[2]--
	return moves
}
