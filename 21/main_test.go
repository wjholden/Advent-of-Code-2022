package main

import (
	"testing"
)

var test string = `root: pppw + sjmn
dbpl: 5
cczh: sllz + lgvd
zczc: 2
ptdq: humn - dvpt
dvpt: 3
lfqf: 4
humn: 5
ljgn: 2
sjmn: drzm * dbpl
sllz: 4
pppw: cczh / lfqf
lgvd: ljgn * ptdq
drzm: hmdt - zczc
hmdt: 32`

func Test1(t *testing.T) {
	want := 152
	got := part1(&test)
	if want != got {
		t.Errorf("incorrect result for part 1 (want %d, got %d)", want, got)
	}
}

func Test2(t *testing.T) {
	want := 301
	got := part2(&test)
	if want != got {
		t.Errorf("incorrect result for part 2 (want %d, got %d)", want, got)
	}
}
