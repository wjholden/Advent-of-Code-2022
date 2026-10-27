package main

import (
	_ "embed"
	"fmt"
	"testing"
)

//go:embed example.txt
var example string

func Test1(t *testing.T) {
	want := 33
	got := 0

	for _, blueprint := range parse(example) {
		memo := make(map[State]int)
		var best int
		geodes := blueprint.collect(State{ore_robots: 1, time: 24}, &memo, &best)
		fmt.Printf("Blueprint %d makes %d geodes\n", blueprint.blueprint_id, geodes)
		got += blueprint.blueprint_id * geodes
	}

	if want != got {
		t.Errorf("incorrect quality level sum (got %d, want %d)", got, want)
	}
}

func Test2(t *testing.T) {
	want := []int{56, 62}
	for i, blueprint := range parse(example) {
		initial := State{ore_robots: 1, time: 32}
		memo := make(map[State]int)
		var best int
		geodes := blueprint.collect(initial, &memo, &best)
		if geodes != want[i] {
			t.Errorf("incorrect amount for blueprint %d (got %d, want %d)", blueprint.blueprint_id, geodes, want[i])
		}
	}
}
