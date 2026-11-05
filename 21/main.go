package main

import (
	_ "embed"
	"fmt"
	"log"
	"strconv"
	"strings"
	"unicode"
)

//go:embed day21.txt
var puzzle string

func main() {
	fmt.Println("Part 1:", part1(&puzzle))
	fmt.Println("Part 2:", part2(&puzzle))
}

func part1(input *string) int {
	monkeys := parse(input)
	return eval("root", &monkeys)
}

func part2(input *string) int {
	monkeys := parse(input)
	if root, ok := monkeys["root"]; ok {
		root.operator = "-"
		monkeys["root"] = root
	} else {
		log.Fatal("failed to find root")
	}
	errors := make(map[int]int)
	evaluate := func(x int) int {
		if humn, ok := monkeys["humn"]; ok {
			humn.value = x
			monkeys["humn"] = humn
		} else {
			log.Fatal("failed to find humn")
		}
		return eval("root", &monkeys)
	}

	lo := 1
	hi := 10000000000000

	errors[lo] = evaluate(lo)
	errors[hi] = evaluate(hi)

	// for m := lo; m <= hi; m *= 10 {
	for {
		if lo == hi {
			return lo
		}
		m := (lo + hi) / 2

		m_err := evaluate(m)
		lo_err := errors[lo]
		hi_err := errors[hi]

		if (lo_err < 0 && hi_err < 0) || (lo_err > 0 && hi_err > 0) {
			log.Fatal("impossible to find zero in this range", lo, hi)
		}

		errors[m] = m_err

		if m_err == 0 {
			return m
		} else if (m_err < 0 && lo_err > 0) || (m_err > 0 && lo_err < 0) {
			hi = m
		} else if (m_err < 0 && hi_err > 0) || (m_err > 0 && hi_err < 0) {
			lo = m
		} else {
			log.Fatal("unexpected state")
		}
	}
}

func eval(u string, monkeys *map[string]Monkey) int {
	if monkey, ok := (*monkeys)[u]; ok {
		if !monkey.has_value {
			a := eval(monkey.input1, monkeys)
			b := eval(monkey.input2, monkeys)
			var c int
			switch monkey.operator {
			case "+":
				c = a + b
			case "-":
				c = a - b
			case "*":
				c = a * b
			case "/":
				c = a / b
			default:
				log.Fatal("unexpected operator in", monkey)
			}
			// monkey.has_value = true
			monkey.value = c
		}
		return monkey.value
	} else {
		log.Fatal("did not find", u)
		return 0
	}
}

type Monkey struct {
	input1    string
	input2    string
	operator  string
	value     int
	has_value bool
}

func parse(input *string) map[string]Monkey {
	monkeys := make(map[string]Monkey)
	delim := func(r rune) bool {
		return unicode.IsSpace(r) || r == ':'
	}
	for line := range strings.Lines(*input) {
		tokens := strings.FieldsFunc(line, delim)
		var monkey Monkey
		switch len(tokens) {
		case 2:
			value, err := strconv.Atoi(tokens[1])
			if err != nil {
				log.Fatal("failed to parse value in", line)
			}
			monkey.value = value
			monkey.has_value = true
		case 4:
			monkey.input1 = tokens[1]
			monkey.operator = tokens[2]
			monkey.input2 = tokens[3]
		default:
			log.Fatal("unexpected tokens length in", line)
		}
		monkeys[tokens[0]] = monkey
	}
	return monkeys
}
