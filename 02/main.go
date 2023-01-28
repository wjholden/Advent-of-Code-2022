package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	// There's probably some clever way to do this, but it's 6am in Europe...
	score := make(map[string]int)
	score["A X"] = 3
	score["A Y"] = 6
	score["A Z"] = 0
	score["B X"] = 0
	score["B Y"] = 3
	score["B Z"] = 6
	score["C X"] = 6
	score["C Y"] = 0
	score["C Z"] = 3

	// Yeah...again there's probably some cool modulo arithmetic solution, but
	// the problem is so small that I don't mind typing out all 9 possible
	// problem states. For part 2, we translate the new input interpretation
	// (that X, Y, and Z tell us what the outcome of the game should be)
	// into the old input interpretation (that X, Y, and Z told us what move
	// we should make). This allows us to use the same scoring functions from
	// before.
	part2 := make(map[string]string)
	// X: lose
	// Y: draw
	// Z: win
	part2["A X"] = "A Z"
	part2["A Y"] = "A X"
	part2["A Z"] = "A Y"
	part2["B X"] = "B X"
	part2["B Y"] = "B Y"
	part2["B Z"] = "B Z"
	part2["C X"] = "C Y"
	part2["C Y"] = "C Z"
	part2["C Z"] = "C X"

	totalScore1 := 0
	totalScore2 := 0

	f, err := os.Open("input.txt")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		totalScore1 += score[line] + shapeScore(line)
		totalScore2 += score[part2[line]] + shapeScore(part2[line])
	}

	fmt.Println("Part 1:", totalScore1)
	fmt.Println("Part 2:", totalScore2)
}

func shapeScore(line string) int {
	switch line[len(line)-1:] {
	case "X":
		return 1
	case "Y":
		return 2
	case "Z":
		return 3
	default:
		panic("Unexpected value")
	}
}
