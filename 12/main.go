package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Key struct {
	s string
	g []int
}

func (k Key) String() string {
	s := k.s + "|"
	for i, n := range k.g {
		if i > 0 {
			s += ","
		}
		s += strconv.Itoa(n)
	}
	return s
}

func main() {
	f, err := os.Open("input")
	if err != nil {
		panic(err)
	}

	part1 := 0
	part2 := 0
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()

		var springs string
		var groupS string
		fmt.Sscanf(line, "%s %s", &springs, &groupS)

		var groups []int
		for _, g := range strings.Split(groupS, ",") {
			v, err := strconv.Atoi(g)
			if err != nil {
				panic(err)
			}
			groups = append(groups, v)
		}

		seen := map[string]int{}

		possibles := doCombos(springs, groups, seen)
		part1 += possibles
		fmt.Println(springs, possibles)

		seen = map[string]int{}
		unfoldedSprings, unfoldedGroups := unfold(springs, groups)
		possibles = doCombos(unfoldedSprings, unfoldedGroups, seen)
		part2 += possibles
		fmt.Println(unfoldedSprings, possibles)
	}

	fmt.Println(part1)
	fmt.Println(part2)
}

func unfold(springs string, groups []int) (string, []int) {
	unfoldedSprings := ""
	var unfoldedGroups []int
	for i := 0; i < 5; i++ {
		if i > 0 {
			unfoldedSprings += "?"
		}
		unfoldedSprings += springs

		unfoldedGroups = append(unfoldedGroups, groups...)
	}

	return unfoldedSprings, unfoldedGroups
}

func sum(n []int) int {
	s := 0
	for _, i := range n {
		s += i
	}
	return s
}

func max(n []int) int {
	max := 0
	for _, i := range n {
		if i > max {
			max = i
		}
	}
	return max
}

func leftover(springs string, groups []int) (string, []int) {
	sz := 0
	idx := 0

	last := 0

	for i, s := range springs {
		switch s {
		case '?':
			return springs[last:], groups[idx:]
		case '.':
			if sz > 0 {
				idx++
			}
			sz = 0
		case '#':
			if sz == 0 {
				last = i
			}
			sz++
		}
	}

	return "", []int{}
}

func doCombos(springs string, groups []int, seen map[string]int) int {
	s, g := leftover(springs, groups)
	if c, ok := seen[Key{s, g}.String()]; ok {
		return c
	}

	i := strings.Index(springs, "?")
	if i < 0 {
		if isValid(springs, groups) {
			return 1
		}
		return 0
	}

	total := sum(groups)

	valid := 0

	numDamaged := strings.Count(springs, "#")
	numUnknown := strings.Count(springs, "?")

	if numDamaged == total {
		remaining := strings.ReplaceAll(springs, "?", ".")
		if isValid(remaining, groups) {
			return 1
		} else {
			return 0
		}
	}

	if numDamaged < total {
		damaged := springs[:i] + "#" + springs[i+1:]
		if isValid(damaged, groups) {
			count := doCombos(damaged, groups, seen)
			s, g = leftover(damaged, groups)
			seen[Key{s, g}.String()] = count
			valid += count
		}
	}

	if (numUnknown + numDamaged) > total {
		operational := springs[:i] + "." + springs[i+1:]
		if isValid(operational, groups) {
			count := doCombos(operational, groups, seen)
			s, g = leftover(operational, groups)
			seen[Key{s, g}.String()] = count
			valid += count
		}
	}

	return valid
}

func isValid(records string, groups []int) bool {
	sz := 0
	idx := 0
	max := max(groups)

	for _, s := range records {
		switch s {
		case '?':
			return true
		case '.':
			if sz > 0 {
				if sz != groups[idx] {
					return false
				}
				idx++
			}
			sz = 0
		case '#':
			if idx >= len(groups) {
				return false
			}
			sz++
			if sz > max {
				return false
			}
		}
	}

	if idx == len(groups) {
		return sz == 0
	}

	if idx == len(groups)-1 {
		return sz == groups[idx]
	}

	return false
}
