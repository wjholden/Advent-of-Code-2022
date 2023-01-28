package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
)

func main() {
	// two compartments for each rucksack
	c1 := make([]string, 0)
	c2 := make([]string, 0)

	// read input from file
	f, err := os.Open("input.txt")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		// I encountered a nasty bounds surprise: the right side of the range is exclusive.
		// I guess this is normal in procedural programming, but I wasn't expecting
		// it after working with scientific languages (Julia, R, Mathematica) for
		// the last few years.
		left := line[:len(line)/2]
		right := line[len(line)/2:]
		c1 = append(c1, left)
		c2 = append(c2, right)
	}

	// iterate over the two halves, find the set intersection, and
	// add up the "priority" values.
	part1 := 0
	for i, v1 := range c1 {
		v2 := c2[i]
		intersection := string_intersection(v1, v2)
		if len(intersection) == 0 {
			panic(c1[i] + c2[i])
		}
		part1 += priority_sum(intersection)
	}
	// 7536 was too low
	fmt.Println("Part 1:", part1)

	// for part 2 we need to iterate over groups of 3
	group_intersections := make([][]string, len(c1)/3)
	part2 := 0
	for i := 0; i < len(c1)/3; i++ {
		g := make([]string, 3)
		for j := 0; j < 3; j++ {
			k := 3*i + j
			g[j] = c1[k] + c2[k]
		}
		g1and2 := strings.Join(string_intersection(g[0], g[1]), "")
		// You might ask why I've kept a copy of the intersection instead of
		// anonymously passing it into the `priority_sum` function. The reason
		// is that keeping the named variable around makes it a little easier
		// for print debugging as I go. This is just a habit I've picked up
		// from working with REPLs. A pipeline would be nice, but I really doubt
		// I'll get one from Go.
		group_intersections[i] = string_intersection(g1and2, g[2])
		part2 += priority_sum(group_intersections[i])
	}
	fmt.Println("Part 2:", part2)
}

func next_different(s []string, i int) int {
	k := i
	for ; k < len(s) && s[i] == s[k]; k++ {
	}
	return k
}

func string_intersection(str1 string, str2 string) []string {
	// sort both strings, then move left to right comparing characters
	s1 := strings.Split(str1, "")
	s2 := strings.Split(str2, "")
	sort.Strings(s1)
	sort.Strings(s2)
	intersection := make([]string, 0)
	// This would have been the easiest one-liner in Python or Julia...
	for i, j := 0, 0; i < len(s1) && j < len(s2); {
		if s1[i] == s2[j] {
			c := s1[i]
			// Convert character to byte for easier arithmetic later.
			// In retrospect, it would have been easier to work with byte arrays
			// throughout.
			intersection = append(intersection, c)
			i = next_different(s1, i)
			j = next_different(s2, j)
		} else if s1[i] < s2[j] {
			i = next_different(s1, i)
		} else {
			j = next_different(s2, j)
		}
	}
	return intersection
}

func priority_sum(B []string) int {
	total := 0
	for _, s := range B {
		b := byte(s[0])
		if 'a' <= b && b <= 'z' {
			total += int(b) - 'a' + 1
		} else if 'A' <= b && b <= 'Z' {
			total += int(b) - 'A' + 27
		} else {
			panic(b)
		}
	}
	return total
}
