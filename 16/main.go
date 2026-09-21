package main

import (
	_ "embed"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
)

//go:embed day16.txt
var puzzle string

type Valve struct {
	flow_rate int
	neighbors []string
}

var valves = map[string]Valve{}

func main() {
	pattern := `Valve (?P<valve>[A-Z]{2}) has flow rate=(?P<rate>\d+); tunnels? leads? to valves? (?P<neighbors>.+)`
	re := regexp.MustCompile(pattern)
	matches := re.FindAllStringSubmatch(puzzle, -1)
	for _, match := range matches {
		label := match[1]
		flow_rate, err := strconv.Atoi(match[2])
		if err != nil {
			log.Fatal(err)
		}
		var neighbors []string
		for _, neighbor := range strings.Split(match[3], ", ") {
			neighbors = append(neighbors, neighbor)
		}
		valves[label] = Valve{flow_rate: flow_rate, neighbors: neighbors}
	}
	fmt.Println(valves)
}
