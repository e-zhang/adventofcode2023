package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Card struct {
	winners []int
	numbers []int
}

func main() {
	f, err := os.Open("input")
	if err != nil {
		panic(err)
	}

	var cards []Card

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()

		card := strings.Split(line, ":")
		values := strings.Split(card[1], "|")

		winners := parseNum(strings.Trim(values[0], " "))
		numbers := parseNum(strings.Trim(values[1], " "))

		cards = append(cards, Card{winners, numbers})
	}

	fmt.Println(part1(cards))
	fmt.Println(part2(cards))
}

func part2(cards []Card) int {
	matching := make([]int, len(cards))

	for i, c := range cards {
		matching[i] = matches(c.numbers, c.winners)
	}

	total := 0
	for i := range cards {
		total += countMatches(matching, i)
	}

	return total
}

func countMatches(matching []int, i int) int {
	sum := 1
	m := matching[i]
	for j := i + 1; j < i+1+m; j++ {
		sum += countMatches(matching, j)
	}
	return sum
}

func parseNum(s string) []int {
	var out []int
	numbers := strings.Split(s, " ")
	for _, n := range numbers {
		if n == "" {
			continue
		}
		v, err := strconv.Atoi(n)
		if err != nil {
			panic(err)
		}
		out = append(out, v)
	}

	return out
}

// finds number of matches from a in b
func matches(a, b []int) int {
	m := 0
	for _, n := range a {
		for _, w := range b {
			if n == w {
				m += 1
				break
			}
		}
	}

	return m
}

func part1(cards []Card) int {
	sum := 0
	for _, card := range cards {
		m := matches(card.numbers, card.winners)
		if m == 0 {
			continue
		}
		points := 1
		for i := 0; i < m-1; i++ {
			points *= 2
		}
		sum += points
	}

	return sum
}
