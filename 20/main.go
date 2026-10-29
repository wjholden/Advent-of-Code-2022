package main

import (
	_ "embed"
	"fmt"
	"log"
	"strconv"
	"strings"
)

//go:embed day20.txt
var puzzle string

type Link struct {
	index int
	value int
	prev  *Link
	next  *Link
}

func (l *Link) String() string {
	current := l
	a := make([]int, 0)
	for {
		a = append(a, (*current).value)
		current = (*current).next
		if current == l {
			break
		}
	}
	return fmt.Sprint(a)
}

func (l *Link) search_index(i int) *Link {
	current := l
	for current.index != i {
		current = current.next
		if current == l {
			log.Fatalf("index %d not in linked list", i)
		}
	}
	return current
}

func (l *Link) search_value(x int) *Link {
	current := l
	for current.value != x {
		current = current.next
		if current == l {
			log.Fatalf("%d not in linked list", x)
		}
	}
	return current
}

func (l *Link) Len() int {
	length := 1
	for p := l.next; p != l; p = p.next {
		length++
	}
	return length
}

func (current *Link) rotate(x int) (left *Link, right *Link) {
	left = current.prev
	right = current.next

	rotations := x
	if x < 0 {
		rotations = -x
	}

	for ; rotations > 0; rotations-- {
		if x < 0 {
			left = left.prev
		} else {
			right = right.next
		}
	}

	if x < 0 {
		right = left.next
	} else {
		left = right.prev
	}

	return
}

func (l *Link) mix(x int, i int, n int) {
	if x == 0 {
		return
	}
	current := l.search_index(i)
	left := current.prev
	right := current.next
	left.next = right
	right.prev = left

	// left, right = current.rotate(x)
	left, right = current.rotate(x % (n - 1)) // subtract one because we clipped one node out

	current.prev = left
	current.next = right
	left.next = current
	right.prev = current
}

func main() {
	fmt.Println("Part 1: ", part1(&puzzle)) // 2842 too high
}

func part1(input *string) int {
	ct := LLparse(input)
	n := ct.Len()
	for i, x := range parse(input) {
		ct.mix(x, i, n)
	}

	// make sure we didn't drop anything
	if ct.Len() != len(parse(input)) {
		log.Fatal("incorrect sequence length")
	}

	fmt.Println(ct, n)
	_, zp1000 := ct.search_value(0).prev.rotate(1000 % n)
	_, zp2000 := zp1000.prev.rotate(1000 % n)
	_, zp3000 := zp2000.prev.rotate(1000 % n)

	fmt.Println(*zp1000, *zp2000, *zp3000)

	return zp1000.value + zp2000.value + zp3000.value
}

func LLparse(input *string) *Link {
	count := 0
	start := Link{}
	var prev *Link
	for line := range strings.Lines(*input) {
		if x, err := strconv.Atoi(strings.TrimSpace(line)); err == nil {
			if count == 0 {
				start.value = x
				prev = &start
			} else {
				current := Link{index: count, value: x, prev: prev}
				(*prev).next = &current
				prev = &current
			}
			count++
		} else {
			log.Fatal(err)
		}
	}
	start.prev = prev
	prev.next = &start
	return &start
}

func parse(input *string) []int {
	result := make([]int, 0)
	for line := range strings.Lines(*input) {
		if x, err := strconv.Atoi(strings.TrimSpace(line)); err == nil {
			result = append(result, x)
		} else {
			log.Fatal(err)
		}
	}
	return result
}
