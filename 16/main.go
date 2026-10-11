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
var good_valves = []string{}

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

	n := len(valves)
	w := floyd_warshall()
	for label, valve := range valves {
		if valve.flow_rate > 0 {
			good_valves = append(good_valves, label)
		}
	}

	p1mask := make([]bool, len(w))
	for i := range p1mask {
		p1mask[i] = true
	}

	fmt.Println("Part 1:", flow(0, 0, START, 0, &w, &p1mask))

	var candidates []int
	mask_count := len(masks(n))
	bar := progressbar.Default(int64(mask_count / 2))
	for j, mask := range masks(len(valves)) {
		if j > mask_count/2 {
			// De-duplicate work where the elephant gets the exact opposite of what we're taking.
			break
		}
		memo = make(map[State]int)
		my_time := flow(0, 4, START, 0, &w, &mask)
		for i := range mask {
			mask[i] = !mask[i]
		}
		memo = make(map[State]int)
		elephant_time := flow(0, 4, START, 0, &w, &mask)
		candidates = append(candidates, my_time+elephant_time)
		bar.Add(1)
	}
	fmt.Println("Part 2:", slices.Max(candidates))
}

type State struct {
	time        int
	position_id int
	open_valves uint64
}

var memo = map[State]int{}

const TIME_CAP = 30

// Got some really excellent automated feedback on this one. My original version used
// strings to represent the open valves. Slices aren't hashable in Go, so to use them
// the memo table I was joining them with a comma. It worked, but it required a lot
// more computation and allocations than we want. The assistant also pointed out that
// we don't need to recalculate the flow rate on every invocation - it's OK to pass
// this along as an argument. Finally, the mask thing. I probably would not have
// thought of this on my own. The mask is a simple means of selecting the valves that
// I take or the elephant takes.
func flow(current_flow int, time int, position string, open_valves uint64, w *[][]float64, mask *[]bool) int {
	position_id := valve_indices[position]
	key := State{time, position_id, open_valves}
	if value, ok := memo[key]; ok {
		return value
	}
	if time >= TIME_CAP {
		return 0
	}
	var candidates []int
	i := valve_indices[position]

	// Warp to any unopened valve with a nonzero flow.
	//
	// Wow! An earlier version of this loop went through all the valves and continued
	// over the ones with zero flow rate. We get almost 3x speed boost by instead
	// keeping track of the good ones and iterating only over those.
	for _, label := range good_valves {
		if label == position {
			continue
		}
		j := valve_indices[label]
		if (1<<j)&open_valves != 0 {
			continue
		}
		if (*mask)[j] == false {
			continue // skip anything that is not included in our mask
		}
		dt := int((*w)[i][j])
		if time+dt+1 > TIME_CAP {
			continue
		}
		new_open_valves := open_valves | (1 << j)
		// We also get a nice speed boost by delaying this lookup until we absolutely know that we need it.
		flow_rate := valves[label].flow_rate
		candidates = append(candidates, current_flow*(dt+1)+flow(current_flow+flow_rate, time+dt+1, label, new_open_valves, w, mask))
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
