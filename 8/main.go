package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Node struct {
	label string

	left  string
	right string

	L *Node
	R *Node
}

func main() {
	f, err := os.Open("input")
	if err != nil {
		panic(err)
	}

	scanner := bufio.NewScanner(f)

	var instruction string

	var nodes []Node
	for scanner.Scan() {
		line := scanner.Text()

		if instruction == "" {
			instruction = line
			continue
		}

		if line == "" {
			continue
		}

		var label, l, r string
		fmt.Sscanf(line, "%s = %s %s", &label, &l, &r)
		l = strings.Trim(l, "(,")
		r = strings.Trim(r, ")")
		nodes = append(nodes, Node{label, l, r, nil, nil})
	}

	fmt.Println(part1(nodes, instruction))
	fmt.Println(part2(nodes, instruction))
}

func part2(nodes []Node, instruction string) int64 {
	var curs []*Node
	for i := range nodes {
		if strings.HasSuffix(nodes[i].label, "A") {
			curs = append(curs, &nodes[i])
		}

		for j := range nodes {
			if nodes[j].label == nodes[i].left {
				nodes[i].L = &nodes[j]
			}

			if nodes[j].label == nodes[i].right {
				nodes[i].R = &nodes[j]
			}
		}
	}

	out := int64(1)
	for _, start := range curs {
		var steps int
		var step int
		curr := start
		cycle := -1
		fmt.Println(curr.label)
		for {
			if strings.HasSuffix(curr.label, "Z") {
				if cycle < 0 {
					cycle = 0
				} else {
					break
				}
			}
			switch instruction[step] {
			case 'L':
				curr = curr.L
			case 'R':
				curr = curr.R
			default:
				panic(instruction[step])
			}
			step = (step + 1) % len(instruction)
			steps++

			if cycle >= 0 {
				cycle++
			}
		}

		out = lcm(out, int64(cycle))
		fmt.Println(start.label, cycle, out)
	}

	return out
}

func gcd(a, b int64) int64 {
	for b != 0 {
		t := b
		b = a % b
		a = t
	}
	return a
}

func lcm(a, b int64) int64 {
	return a * b / gcd(a, b)
}

func part1(nodes []Node, instruction string) int {
	var start *Node
	for i := range nodes {
		if nodes[i].label == "AAA" {
			start = &nodes[i]
		}

		for j := range nodes {
			if nodes[j].label == nodes[i].left {
				nodes[i].L = &nodes[j]
			}

			if nodes[j].label == nodes[i].right {
				nodes[i].R = &nodes[j]
			}
		}
	}

	var step int
	var steps int
	curr := start
	for curr.label != "ZZZ" {
		fmt.Println(curr.label, string(instruction[step]))
		switch instruction[step] {
		case 'L':
			curr = curr.L
		case 'R':
			curr = curr.R
		default:
			panic(instruction[step])
		}

		step = (step + 1) % len(instruction)
		steps++
	}

	return steps
}
