package main

import (
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

type monkey struct {
	div   int
	op    string
	t     int
	f     int
	items []int
}

func main() {
	monkeys := parse(puzzle)
	inspections := make([]int, len(monkeys))
	for round := 0; round < 20; round++ {
		//fmt.Println("Round", round, monkeys)
		for i := range monkeys {
			// Why pointers? Because Go will pass by value, not reference.
			// So, you're getting a copy of the monkey object. If you make any
			// changes to it, those won't propagate back to the original.
			m := &monkeys[i]
			for _, oldWorry := range m.items {
				inspections[i]++
				newWorry := eval(m.op, oldWorry)
				newWorry /= 3
				var dst *monkey
				if newWorry%m.div == 0 {
					dst = &monkeys[m.t]
				} else {
					dst = &monkeys[m.f]
				}
				dst.items = append(dst.items, newWorry)
			}
			m.items = m.items[:0] // clear array
		}
	}
	sort.Sort(sort.Reverse(sort.IntSlice(inspections)))
	fmt.Println("Part 1:", inspections[0]*inspections[1])
}

const pattern = `Monkey (?P<id>\d):
  Starting items: (?P<items>.+)
  Operation: new = (?P<op>.+)
  Test: divisible by (?P<div>\d+)
    If true: throw to monkey (?P<t>\d)
    If false: throw to monkey (?P<f>\d)`

func parse(input string) []monkey {
	M := strings.Split(input, "\n\n")
	n := len(M)
	if n != strings.Count(input, "Monkey") {
		panic("something is wrong parsing Monkey")
	}
	monkeys := make([]monkey, n)
	var r = regexp.MustCompile(pattern)
	for i := 0; i < n; i++ {
		var m monkey
		match := r.FindStringSubmatch(M[i])

		// Starting items
		holding := strings.Split(match[2], ", ")
		for _, itemHeld := range holding {
			worry, _ := strconv.Atoi(itemHeld)
			m.items = append(m.items, worry)
		}

		// Operation
		m.op = match[3]

		// Test
		div, _ := strconv.Atoi(match[4])
		m.div = div
		t, _ := strconv.Atoi(match[5])
		m.t = t
		f, _ := strconv.Atoi(match[6])
		m.f = f

		monkeys[i] = m
	}
	return monkeys
}

var er = regexp.MustCompile(`\+|\*`)

func eval(f string, old int) int {
	var x [2]int
	fn := er.Split(f, -1)
	if len(fn) != len(x) {
		panic("error splitting " + f)
	}
	for i, v := range fn {
		v = strings.TrimSpace(v)
		if v == "old" {
			x[i] = old
		} else {
			x[i], _ = strconv.Atoi(v)
		}
	}
	//fmt.Println(x)

	switch {
	case strings.Contains(f, "+"):
		return x[0] + x[1]
	case strings.Contains(f, "*"):
		return x[0] * x[1]
	default:
		panic("no applicable operator in " + f)
	}
}

const sample = `Monkey 0:
  Starting items: 79, 98
  Operation: new = old * 19
  Test: divisible by 23
    If true: throw to monkey 2
    If false: throw to monkey 3

Monkey 1:
  Starting items: 54, 65, 75, 74
  Operation: new = old + 6
  Test: divisible by 19
    If true: throw to monkey 2
    If false: throw to monkey 0

Monkey 2:
  Starting items: 79, 60, 97
  Operation: new = old * old
  Test: divisible by 13
    If true: throw to monkey 1
    If false: throw to monkey 3

Monkey 3:
  Starting items: 74
  Operation: new = old + 3
  Test: divisible by 17
    If true: throw to monkey 0
    If false: throw to monkey 1`

const puzzle = `Monkey 0:
  Starting items: 59, 74, 65, 86
  Operation: new = old * 19
  Test: divisible by 7
    If true: throw to monkey 6
    If false: throw to monkey 2

Monkey 1:
  Starting items: 62, 84, 72, 91, 68, 78, 51
  Operation: new = old + 1
  Test: divisible by 2
    If true: throw to monkey 2
    If false: throw to monkey 0

Monkey 2:
  Starting items: 78, 84, 96
  Operation: new = old + 8
  Test: divisible by 19
    If true: throw to monkey 6
    If false: throw to monkey 5

Monkey 3:
  Starting items: 97, 86
  Operation: new = old * old
  Test: divisible by 3
    If true: throw to monkey 1
    If false: throw to monkey 0

Monkey 4:
  Starting items: 50
  Operation: new = old + 6
  Test: divisible by 13
    If true: throw to monkey 3
    If false: throw to monkey 1

Monkey 5:
  Starting items: 73, 65, 69, 65, 51
  Operation: new = old * 17
  Test: divisible by 11
    If true: throw to monkey 4
    If false: throw to monkey 7

Monkey 6:
  Starting items: 69, 82, 97, 93, 82, 84, 58, 63
  Operation: new = old + 5
  Test: divisible by 5
    If true: throw to monkey 5
    If false: throw to monkey 7

Monkey 7:
  Starting items: 81, 78, 82, 76, 79, 80
  Operation: new = old + 3
  Test: divisible by 17
    If true: throw to monkey 3
    If false: throw to monkey 4`
