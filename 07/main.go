// https://go.dev/play/p/dMN966wRmNZ
package main

import (
	_ "embed"
	"fmt"
	"strconv"
	"strings"
)

const sample = `$ cd /
$ ls
dir a
14848514 b.txt
8504156 c.dat
dir d
$ cd a
$ ls
dir e
29116 f
2557 g
62596 h.lst
$ cd e
$ ls
584 i
$ cd ..
$ cd ..
$ cd d
$ ls
4060174 j
8033020 d.log
5626152 d.ext
7214296 k`

//go:embed day07.txt
var puzzle string

// I need to spend some more time with Go's interfaces.
// This would have been an easy problem to solve in Java.
//
// In Java, I would have defined some "fileItem" and "dirItem"
// types that both implement some "treeItem" supertype.
// Using the "instanceof" operator, it would have been easy
// to distinguish one from another and provide guarantees.
// As written, this program does not use inheritance at all.
// Instead, I just use "size == 0" to distinguish a directory
// from a file.
//
// This problem also might have been neatly solved using SQL.
type item struct {
	name     string
	size     int
	children map[string]*item
}

// Not all directory names are unique!
var dirSizes []int

func (i item) String() string {
	if i.size == 0 {
		return fmt.Sprintf("%s (dir)", i.name)
	} else {
		return fmt.Sprintf("%s (file, size=%d)", i.name, i.size)
	}
}

func main() {
	rootSize := walk(parse(strings.TrimSpace(puzzle)))
	part1, part2 := 0, 70000000
	for _, size := range dirSizes {
		if size <= 100000 {
			part1 += size
		}
		if size < part2 && rootSize-size+30000000 <= 70000000 {
			part2 = size
		}
	}
	fmt.Printf("Part 1: %d\n", part1)
	fmt.Printf("Part 2: %d\n", part2)
}

func walk(tree *item) int {
	size := tree.size
	if size == 0 { // a directory has size 0
		for _, child := range tree.children {
			size += walk(child)
		}
		dirSizes = append(dirSizes, size)
	}
	return size
}

func parse(s string) *item {
	ch := make(chan string)
	go iterate(s, ch)

	root := new(item)
	root.children = make(map[string]*item)
	var stack []*item
	stack = append(stack, root)

	for line := range ch {
		cwd := stack[len(stack)-1]

		tokens := strings.Split(line, " ")
		switch {
		// Initial cd to the root. Really, this could happen anytime, but in the
		// problem we actually just see this as the very first command, which
		// simplifies the problem.
		case line == "$ cd /":
			cwd.name = "/"
		// Go back to the parent directory. We just need to pop one item off the stack.
		case line == "$ cd ..":
			stack = stack[:len(stack)-1]
		// Enter a child directory (if it exists).
		case strings.HasPrefix(line, "$ cd"):
			child, ok := cwd.children[tokens[2]]
			if !ok {
				panic("trying to cd to non-existant directory " + tokens[2])
			}
			stack = append(stack, child)
		// This problem has only two commands (cd and ls), so really we can interpret
		// anything that isn't a command as an output of ls. If this were a harder
		// problem, then we would have to keep track of which command had been issued
		// to distinguish one output from another.
		case line == "$ ls":
			// do nothing
			continue
		// Lines of output from ls start with "dir" if it is a directory.
		case strings.HasPrefix(line, "dir"): // directories
			childItem := new(item)
			childItem.name = tokens[1]
			childItem.children = make(map[string]*item)
			cwd.children[childItem.name] = childItem
		// Lines of output from ls contain a number and filename if it is an ordinary file.
		default:
			childItem := new(item)
			childItem.name = tokens[1]
			i, err := strconv.Atoi(tokens[0])
			if err != nil {
				panic("failed to read filesize in " + line)
			}
			childItem.size = i
			cwd.children[childItem.name] = childItem
		}
	}
	return root
}

func iterate(s string, ch chan<- string) {
	for _, line := range strings.Split(s, "\n") {
		ch <- line
	}
	close(ch)
}
