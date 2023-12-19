package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Rule struct {
	category string
	op       string
	value    int

	action string
}

func (r *Rule) cmp(v int) bool {
	switch r.op {
	case "<":
		return v < r.value
	case ">":
		return v > r.value
	}
	panic(r.op)
}

func (r *Rule) Check(part Ratings) bool {
	match := false
	switch r.category {
	case "x":
		match = r.cmp(part.x)
	case "m":
		match = r.cmp(part.m)
	case "a":
		match = r.cmp(part.a)
	case "s":
		match = r.cmp(part.s)
	default:
		match = true
	}

	return match
}

type Workflow struct {
	name  string
	rules []Rule
}

func (w *Workflow) Apply(part Ratings) string {
	for _, r := range w.rules {
		if r.Check(part) {
			return r.action
		}
	}

	panic(part)
}

func parseWorkflow(line string) Workflow {
	start := strings.Index(line, "{")
	end := strings.Index(line, "}")
	name := line[:start]

	var w Workflow
	w.name = name

	rules := strings.Split(line[start+1:end], ",")
	for _, r := range rules {
		seg := strings.Split(r, ":")
		if len(seg) == 1 {
			w.rules = append(w.rules, Rule{action: r})
			continue
		}
		idx := strings.IndexAny(seg[0], "<>")
		c := seg[0][:idx]
		op := string(seg[0][idx])
		v, err := strconv.Atoi(seg[0][idx+1:])
		if err != nil {
			panic(err)
		}

		w.rules = append(w.rules, Rule{c, op, v, seg[1]})
	}
	return w
}

type Ratings struct {
	x int
	m int
	a int
	s int
}

func parseRatings(line string) Ratings {
	start := strings.Index(line, "{")
	end := strings.Index(line, "}")

	rest := strings.Split(line[start+1:end], ",")
	var r Ratings
	for _, p := range rest {
		seg := strings.Split(p, "=")
		v, err := strconv.Atoi(seg[1])
		if err != nil {
			panic(err)
		}

		switch seg[0] {
		case "x":
			r.x = v
		case "m":
			r.m = v
		case "a":
			r.a = v
		case "s":
			r.s = v
		}
	}

	return r
}

func main() {
	f, err := os.Open("input")
	if err != nil {
		panic(err)
	}

	workflows := map[string]Workflow{}
	ratings := []Ratings{}
	scanner := bufio.NewScanner(f)
	divider := false
	for scanner.Scan() {
		line := scanner.Text()

		if len(line) == 0 {
			divider = true
			continue
		}

		if !divider {
			wf := parseWorkflow(line)
			workflows[wf.name] = wf
		} else {
			ratings = append(ratings, parseRatings(line))
		}
	}

	part2(workflows, ratings)
}

func part1(workflows map[string]Workflow, ratings []Ratings) {
	sum := 0
	for _, part := range ratings {
		fmt.Println(part)
		state := "in"
		for state != "A" && state != "R" {
			fmt.Printf("%s -> ", state)
			wf := workflows[state]
			state = wf.Apply(part)
		}
		fmt.Printf("%s\n", state)

		if state == "A" {
			sum += part.x + part.m + part.a + part.s
		}
	}
	fmt.Println(sum)
}

func part2(workflows map[string]Workflow, ratings []Ratings) {
	sum := checkPossibilities(TOTAL_COMBOS, "in", workflows, map[string]Range{})
	sum2 := checkPossibilities2("in", workflows, map[string]Range{
		"x": Range{MIN, MAX},
		"m": Range{MIN, MAX},
		"a": Range{MIN, MAX},
		"s": Range{MIN, MAX},
	})
	fmt.Println(sum, sum2)
}

const (
	MIN          = 1
	MAX          = 4000
	TOTAL_COMBOS = MAX * MAX * MAX * MAX
)

type Range struct {
	min int
	max int
}

func checkPossibilities(remaining int, name string, workflows map[string]Workflow, ranges map[string]Range) int {
	if name == "A" {
		return remaining
	}

	if name == "R" {
		return 0
	}

	copyRanges := map[string]Range{}
	for k, v := range ranges {
		copyRanges[k] = v
	}

	wf := workflows[name]
	possibilities := 0
	for _, r := range wf.rules {
		if r.category == "" {
			possibilities += checkPossibilities(remaining, r.action, workflows, copyRanges)
			remaining = 0
		} else {
			rem, ok := copyRanges[r.category]
			if !ok {
				rem = Range{MIN, MAX}
			}

			if r.op == ">" {
				c := (rem.max - r.value) * remaining / (rem.max - rem.min + 1)
				copyRanges[r.category] = Range{r.value + 1, rem.max}
				possibilities += checkPossibilities(c, r.action, workflows, copyRanges)
				copyRanges[r.category] = Range{rem.min, r.value}
				remaining -= c

			} else if r.op == "<" {
				c := (r.value - rem.min) * remaining / (rem.max - rem.min + 1)
				copyRanges[r.category] = Range{rem.min, r.value - 1}
				possibilities += checkPossibilities(c, r.action, workflows, copyRanges)
				copyRanges[r.category] = Range{r.value, rem.max}
				remaining -= c
			}
		}
	}

	return possibilities
}

func checkPossibilities2(name string, workflows map[string]Workflow, ranges map[string]Range) int {
	if name == "A" {
		possibilities := 1
		for _, r := range ranges {
			possibilities *= r.max - r.min + 1
		}
		return possibilities
	}

	if name == "R" {
		return 0
	}

	copyRanges := map[string]Range{}
	for k, v := range ranges {
		copyRanges[k] = v
	}

	wf := workflows[name]
	possibilities := 0
	for _, r := range wf.rules {
		if r.category == "" {
			possibilities += checkPossibilities2(r.action, workflows, copyRanges)
		} else {
			rem := copyRanges[r.category]
			if r.op == ">" {
				copyRanges[r.category] = Range{r.value + 1, rem.max}
				possibilities += checkPossibilities2(r.action, workflows, copyRanges)
				copyRanges[r.category] = Range{rem.min, r.value}
			} else if r.op == "<" {
				copyRanges[r.category] = Range{rem.min, r.value - 1}
				possibilities += checkPossibilities2(r.action, workflows, copyRanges)
				copyRanges[r.category] = Range{r.value, rem.max}
			}
		}
	}

	return possibilities
}
