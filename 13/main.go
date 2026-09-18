package main

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
)

type PacketPair struct {
	a, b *interface{}
}

//go:embed day13.txt
var puzzle string

func parse(s string) []PacketPair {
	n := 1 + strings.Count(s, "\n\n")
	pairs := make([]PacketPair, n)
	for i, lines := range strings.Split(s, "\n\n") {
		var a, b interface{}
		packets := strings.Split(lines, "\n")
		_ = json.Unmarshal([]byte(packets[0]), &a)
		_ = json.Unmarshal([]byte(packets[1]), &b)
		pairs[i] = PacketPair{&a, &b}
	}
	return pairs
}

func collect(pairs []PacketPair) []*interface{} {
	packets := make([]*interface{}, 0)
	for _, pair := range pairs {
		packets = append(packets, pair.a, pair.b)
	}
	return packets
}

func newPacket(key string) *interface{} {
	var a interface{}
	_ = json.Unmarshal([]byte(key), &a)
	return &a
}

func compare(a, b *interface{}) int {
	A := (*a).([]interface{})
	B := (*b).([]interface{})

	n := len(A)
	if n < len(B) {
		n = len(B)
	}

	for i := 0; i < n; i++ {
		if i >= len(B) {
			//fmt.Println("Right side", B, "ran out of items.")
			return -1
		}
		if i >= len(A) {
			//fmt.Println("Left side", A, "ran out of items.")
			return +1
		}
		left := A[i]
		right := B[i]

		x, ok1 := left.(float64)
		y, ok2 := right.(float64)

		if ok1 && ok2 {
			switch {
			case x > y:
				//fmt.Printf("%.0f > %.0f\n", x, y)
				return -1
			case x < y:
				return +1
			default:
				continue
			}
		} else if !ok1 && ok2 {
			tmp := []interface{}{y}
			var newRight interface{} = tmp
			retVal := compare(&left, &newRight)
			if retVal != 0 {
				return retVal
			}
		} else if ok1 && !ok2 {
			tmp := []interface{}{x}
			var newLeft interface{} = tmp
			retVal := compare(&newLeft, &right)
			if retVal != 0 {
				return retVal
			}
		} else {
			retVal := compare(&left, &right)
			if retVal != 0 {
				return retVal
			}
		}
	}

	return 0
}

func main() {
	input := puzzle
	var part1 int
	for i, pair := range parse(input) {
		//fmt.Println(*pair.a, "compare to", *pair.b)
		retVal := compare(pair.a, pair.b)
		if retVal == 1 {
			part1 += 1 + i
		}
		if retVal == 0 {
			fmt.Println("Unexpected value 0!")
		}
		//fmt.Println(retVal)
	}
	fmt.Println("Part 1:", part1)

	packets := collect(parse(input))
	divider1 := newPacket("[[2]]")
	divider2 := newPacket("[[6]]")
	packets = append(packets, divider1, divider2)
	sort.Slice(packets, func(i, j int) bool { return 0 <= compare(packets[i], packets[j]) })
	part2 := 1
	for i, packet := range packets {
		if packet == divider1 || packet == divider2 {
			part2 *= 1 + i
		}
	}
	fmt.Println("Part 2:", part2)
}

const sample = `[1,1,3,1,1]
[1,1,5,1,1]

[[1],[2,3,4]]
[[1],4]

[9]
[[8,7,6]]

[[4,4],4,4]
[[4,4],4,4,4]

[7,7,7,7]
[7,7,7]

[]
[3]

[[[]]]
[[]]

[1,[2,[3,[4,[5,6,7]]]],8,9]
[1,[2,[3,[4,[5,6,0]]]],8,9]`
