package main

import (
	_ "embed"
	"fmt"
	"log"
	"maps"
	"math"
	"regexp"
	"slices"
	"strconv"
	"strings"

	"github.com/schollz/progressbar/v3"
)

//go:embed day16.txt
var puzzle string

type Valve struct {
	flow_rate int
	neighbors []string
}

var valves = map[string]Valve{}

const START = "AA"

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

	w := floyd_warshall()
	// fmt.Println(w)

	p1mask := make([]bool, len(w))
	for i := range p1mask {
		p1mask[i] = true
	}

	fmt.Println("Part 1:", flow(0, START, []string{}, &w, &p1mask))

	// fmt.Println(masks(len(w)))
	// fmt.Println(len(masks(len(w))))

	// Slow but it works!
	var candidates []int
	bar := progressbar.Default(int64(len(masks(len(valves)))))
	for _, mask := range masks(len(valves)) {
		memo = make(map[State]int)
		my_time := flow(4, START, []string{}, &w, &mask)
		for i := range mask {
			mask[i] = !mask[i]
		}
		memo = make(map[State]int)
		elephant_time := flow(4, START, []string{}, &w, &mask)
		candidates = append(candidates, my_time+elephant_time)
		bar.Add(1)
	}
	fmt.Println("Part 2:", slices.Max(candidates))
}

type State struct {
	time        int
	position    string
	open_valves string
}

var memo = map[State]int{}

const TIME_CAP = 30

func flow(time int, position string, open_valves []string, w *[][]float64, mask *[]bool) int {
	slices.Sort(open_valves)
	key := State{time, position, strings.Join(open_valves, ",")}
	if value, ok := memo[key]; ok {
		// fmt.Println("Cache hit!")
		return value
	}
	// fmt.Printf("Searching %s time %d\n", position, time)
	var current_flow int
	for _, s := range open_valves {
		current_flow += valves[s].flow_rate
	}
	if time >= TIME_CAP {
		return 0
	}
	var candidates []int
	i := valve_indices[position]

	// Warp to any unopened valve with a nonzero flow.
	for label, valve := range valves {
		if label == position || valve.flow_rate == 0 || slices.Contains(open_valves, label) {
			continue
		}
		j := valve_indices[label]
		if (*mask)[j] == false {
			continue // skip anything that is not included in our mask
		}
		dt := int((*w)[i][j])
		if time+dt+1 > TIME_CAP {
			continue
		}
		// Lots of very subtle bugs here.
		lets_go_open_it := slices.Clone(open_valves)
		lets_go_open_it = append(lets_go_open_it, label)
		candidates = append(candidates, current_flow*(dt+1)+flow(time+dt+1, label, lets_go_open_it, w, mask))
	}
	var value int
	if len(candidates) == 0 {
		// Looks like nothing was available to open. Keep the current flow going up to the time cap.
		value = current_flow * (TIME_CAP - time)
	} else {
		value = slices.Max(candidates)
	}
	memo[key] = value
	return value
}

var valve_indices = map[string]int{}
var flow_rates = []int{}

// Why floats? Because Floyd-Warshall needs infinities that do what you expect.
// If you use integers, then you'll need to handle overflow very carefully.
func floyd_warshall() [][]float64 {
	// Our input is an adjacency list. We'll need a stable way to translate
	// node names to rows and columns in our distances matrix later.
	keys := slices.Collect(maps.Keys(valves))
	slices.Sort(keys)
	for i, s := range keys {
		valve_indices[s] = i
		flow_rates = append(flow_rates, valves[s].flow_rate)
	}

	// Buncha boilerplate to initialize our distances matrix.
	n := len(keys)
	d := make([][]float64, 0)
	for r := range n {
		row := make([]float64, n)
		current := keys[r]
		valve := valves[current]
		for c := range n {
			if r == c {
				row[c] = 0
			} else if slices.Contains(valve.neighbors, keys[c]) {
				row[c] = 1
			} else {
				row[c] = math.Inf(1)
			}
		}
		d = append(d, row)
	}

	// This is the core of the Floyd-Warshall algorithm.
	for k := range n {
		for i := range n {
			for j := range n {
				x := d[i][k] + d[k][j]
				if x < d[i][j] {
					d[i][j] = x
				}
			}
		}
	}

	return d
}

func masks(n int) [][]bool {
	if n == 0 {
		return [][]bool{{}}
	} else {
		var result [][]bool
		for _, mask := range masks(n - 1) {
			// Only create a true branch for nonzero flow rates.
			if flow_rates[n-1] > 0 {
				tmask := slices.Clone(mask)
				tmask = append(tmask, true)
				result = append(result, tmask)
			}
			fmask := append(mask, false)
			result = append(result, fmask)
		}
		return result
	}
}
