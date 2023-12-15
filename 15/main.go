package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	f, err := os.Open("input")
	if err != nil {
		panic(err)
	}

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		sequence := strings.Split(line, ",")
		part1(sequence)
		part2(sequence)

	}
}

type Lens struct {
	label string
	focal int
}

func part2(sequence []string) {
	boxes := make([][]Lens, 256)

	for _, s := range sequence {
		idx := strings.IndexAny(s, "=-")
		if idx < 0 {
			panic(s)
		}

		label := s[:idx]
		box := hash(label)
		switch s[idx] {
		case '=':
			f, err := strconv.Atoi(s[idx+1:])
			if err != nil {
				panic(err)
			}

			boxes[box] = addToBox(boxes[box], Lens{label, f})
		case '-':
			boxes[box] = removeFromBox(boxes[box], label)
		}

		// fmt.Println(s)
		// printBoxes(boxes)
	}

	power := 0
	for i, b := range boxes {
		for j, l := range b {
			power += (i + 1) * (j + 1) * l.focal
		}
	}
	fmt.Println(power)
}

func printBoxes(boxes [][]Lens) {
	for i, b := range boxes {
		if len(b) == 0 {
			continue
		}

		fmt.Printf("Box %d: ", i)
		for _, l := range b {
			fmt.Printf("[%s %d]", l.label, l.focal)
		}
		fmt.Println()
	}
}

func addToBox(box []Lens, l Lens) []Lens {
	add := true
	for i, lens := range box {
		if lens.label == l.label {
			box[i] = l
			add = false
			break
		}
	}

	if add {
		box = append(box, l)
	}

	return box
}

func removeFromBox(box []Lens, label string) []Lens {
	idx := -1
	for i, l := range box {
		if l.label == label {
			idx = i
			break
		}
	}

	if idx >= 0 {
		box = append(box[:idx], box[idx+1:]...)
	}
	return box
}

func part1(sequence []string) {
	sum := 0
	for _, s := range sequence {
		h := hash(s)
		fmt.Println(s, h)
		sum += h
	}
	fmt.Println(sum)
}

func hash(s string) int {
	v := 0
	for _, c := range s {
		v += int(c)
		v *= 17
		v %= 256
	}
	return v
}
