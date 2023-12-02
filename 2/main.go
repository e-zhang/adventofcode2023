package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

const (
	GREENS = 13
	REDS   = 12
	BLUES  = 14
)

func main() {
	f, err := os.Open("input")
	if err != nil {
		panic(err)
	}

	scanner := bufio.NewScanner(f)

	possibles := []int{}
	powers := []int{}
	for scanner.Scan() {
		line := scanner.Text()

		var id int
		fmt.Sscanf(line, "Game %d: ", &id)
		subsets := strings.Split(line, ":")[1]
		sets := strings.Split(subsets, ";")

		if part1(sets) {
			possibles = append(possibles, id)
		}

		g, r, b := part2(sets)
		powers = append(powers, g*r*b)
	}

	ids := 0
	for _, n := range possibles {
		ids += n
	}
	fmt.Printf("Part 1: %d\n", ids)

	power := 0
	for _, n := range powers {
		power += n
	}
	fmt.Printf("Part 2: %d\n", power)
}

func part2(sets []string) (g, r, b int) {
	var greens, reds, blues int
	for _, set := range sets {
		cubes := strings.Split(set, ",")
		for _, cube := range cubes {
			var n int
			var color string
			fmt.Sscanf(cube, "%d %s", &n, &color)
			switch color {
			case "blue":
				if n > blues {
					blues = n
				}
			case "green":
				if n > greens {
					greens = n
				}
			case "red":
				if n > reds {
					reds = n
				}
			}
		}
	}

	return greens, reds, blues
}

func part1(sets []string) bool {
	for _, set := range sets {
		cubes := strings.Split(set, ",")
		for _, cube := range cubes {
			var n int
			var color string
			fmt.Sscanf(cube, "%d %s", &n, &color)
			switch color {
			case "blue":
				if n > BLUES {
					return false
				}
			case "red":
				if n > REDS {
					return false
				}
			case "green":
				if n > GREENS {
					return false
				}
			}
		}
	}
	return true
}
