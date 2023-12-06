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

	var times, distances []int
	var t2, d2 int
	for scanner.Scan() {
		line := scanner.Text()
		splits := strings.Split(line, " ")
		if splits[0] == "Time:" {
			times, t2 = parseRace(splits[1:])
		} else {
			distances, d2 = parseRace(splits[1:])
		}
	}

	ways := 1
	for i, t := range times {
		dist := distances[i]
		wins := calculateWins(t, dist)
		ways *= wins
	}
	fmt.Println(ways)

	fmt.Println(t2, d2)
	fmt.Println(calculateWins(t2, d2))
}

func calculateWins(time, distance int) int {
	var wins int
	for hold := 0; hold <= time; hold++ {
		d := hold * (time - hold)
		if d > distance {
			wins++
		}
	}

	return wins
}

func parseRace(s []string) ([]int, int) {
	var nums []int
	var p2val string
	for _, v := range s {
		if v == "" {
			continue
		}

		n, err := strconv.Atoi(v)
		if err != nil {
			panic(err)
		}
		nums = append(nums, n)
		p2val += v
	}

	val, err := strconv.Atoi(p2val)
	if err != nil {
		panic(err)
	}

	return nums, val
}
