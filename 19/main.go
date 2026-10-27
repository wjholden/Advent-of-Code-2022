package main

import (
	_ "embed"
	"fmt"
	"log"
	"regexp"
	"strconv"
)

//go:embed day19.txt
var puzzle string

func main() {
	quality := 0
	blueprints := parse(puzzle)
	for _, blueprint := range blueprints {
		state := State{ore_robots: 1, time: 24}
		memo := make(map[State]int)
		var best int
		geodes := blueprint.collect(state, &memo, &best)
		fmt.Printf("Blueprint %d makes %d geodes\n", blueprint.blueprint_id, geodes)
		quality += blueprint.blueprint_id * geodes
	}
	fmt.Println("Part 1:", quality)

	product := 1
	for _, blueprint := range blueprints[:3] {
		state := State{ore_robots: 1, time: 32}
		memo := make(map[State]int)
		var best int
		geodes := blueprint.collect(state, &memo, &best)
		product *= geodes
	}
	fmt.Println("Part 2:", product)
}

type Blueprint struct {
	blueprint_id   int
	ore_ore        int
	clay_ore       int
	obsidian_ore   int
	obsidian_clay  int
	geode_ore      int
	geode_obsidian int
}

type State struct {
	time            int
	ore_robots      int
	ore             int
	clay_robots     int
	clay            int
	obsidian_robots int
	obsidian        int
	geode_robots    int
	geode           int
}

func (state *State) tick() {
	state.ore += state.ore_robots
	state.clay += state.clay_robots
	state.obsidian += state.obsidian_robots
	state.geode += state.geode_robots
	state.time--
}

func (blueprint *Blueprint) max_ore_usage() int {
	return max(blueprint.ore_ore, blueprint.clay_ore, blueprint.obsidian_ore, blueprint.geode_ore)
}

func (blueprint *Blueprint) max_clay_usage() int {
	return blueprint.obsidian_clay
}

func (blueprint *Blueprint) max_obsidian_usage() int {
	return blueprint.geode_obsidian
}

func (blueprint *Blueprint) collect(state State, memo *map[State]int, best *int) int {
	// Clamp values to the max usable amount. This makes our memo table smaller.
	max_ore := state.time * blueprint.max_ore_usage()
	if state.ore > max_ore {
		state.ore = max_ore
	}

	max_clay := state.time * blueprint.max_clay_usage()
	if state.clay > max_clay {
		state.clay = max_clay
	}

	max_obsidian := state.time * blueprint.max_obsidian_usage()
	if state.obsidian > max_obsidian {
		state.obsidian = max_obsidian
	}

	// Can we beat the best case we've found so far?
	t := state.time
	g := state.geode + t*state.geode_robots + (t*(t-1))/2
	if g <= *best {
		return *best
	}

	if g, ok := (*memo)[state]; ok {
		return g
	}

	var geode int

	if state.geode > *best {
		*best = state.geode
	}

	if state.time == 0 {
		(*memo)[state] = state.geode
		return state.geode
	}

	if state.ore >= blueprint.geode_ore && state.obsidian >= blueprint.geode_obsidian {
		copy := state
		copy.ore -= blueprint.geode_ore
		copy.obsidian -= blueprint.geode_obsidian
		copy.tick()
		copy.geode_robots++
		geode = max(blueprint.collect(copy, memo, best), geode)
	}

	if state.ore >= blueprint.obsidian_ore && state.clay >= blueprint.obsidian_clay && state.obsidian_robots < blueprint.max_obsidian_usage() && state.obsidian < state.time*blueprint.max_obsidian_usage() {
		copy := state
		copy.ore -= blueprint.obsidian_ore
		copy.clay -= blueprint.obsidian_clay
		copy.tick()
		copy.obsidian_robots++
		geode = max(blueprint.collect(copy, memo, best), geode)
	}

	if state.ore >= blueprint.clay_ore && state.clay_robots < blueprint.max_clay_usage() && state.clay < state.time*blueprint.max_clay_usage() {
		copy := state
		copy.ore -= blueprint.clay_ore
		copy.tick()
		copy.clay_robots++
		geode = max(blueprint.collect(copy, memo, best), geode)
	}

	if state.ore >= blueprint.ore_ore && state.ore_robots < blueprint.max_ore_usage() && state.ore < state.time*blueprint.max_ore_usage() {
		copy := state
		copy.ore -= blueprint.ore_ore
		copy.tick()
		copy.ore_robots++
		geode = max(blueprint.collect(copy, memo, best), geode)
	}

	copy := state
	copy.tick()
	geode = max(blueprint.collect(copy, memo, best), geode)

	(*memo)[state] = geode

	return geode
}

func parse(input string) []Blueprint {
	re := regexp.MustCompile(`Blueprint (?P<blueprint_id>\d+): Each ore robot costs (?P<ore_ore>\d+) ore. Each clay robot costs (?P<clay_ore>\d+) ore. Each obsidian robot costs (?P<obsidian_ore>\d+) ore and (?P<obsidian_clay>\d+) clay. Each geode robot costs (?P<geode_ore>\d+) ore and (?P<geode_obsidian>\d+) obsidian.`)
	matches := re.FindAllStringSubmatch(input, -1)
	blueprints := make([]Blueprint, 0)
	for _, match := range matches {
		blueprint := Blueprint{}

		blueprint_id, err := strconv.Atoi(match[1])
		if err != nil {
			log.Fatal(err)
		}
		blueprint.blueprint_id = blueprint_id

		ore_ore, err := strconv.Atoi(match[2])
		if err != nil {
			log.Fatal(err)
		}
		blueprint.ore_ore = ore_ore

		clay_ore, err := strconv.Atoi(match[3])
		if err != nil {
			log.Fatal(err)
		}
		blueprint.clay_ore = clay_ore

		obsidian_ore, err := strconv.Atoi(match[4])
		if err != nil {
			log.Fatal(err)
		}
		blueprint.obsidian_ore = obsidian_ore

		obsidian_clay, err := strconv.Atoi(match[5])
		if err != nil {
			log.Fatal(err)
		}
		blueprint.obsidian_clay = obsidian_clay

		geode_ore, err := strconv.Atoi(match[6])
		if err != nil {
			log.Fatal(err)
		}
		blueprint.geode_ore = geode_ore

		geode_obsidian, err := strconv.Atoi(match[7])
		if err != nil {
			log.Fatal(err)
		}
		blueprint.geode_obsidian = geode_obsidian

		fmt.Println(blueprint)
		blueprints = append(blueprints, blueprint)
	}
	return blueprints
}
