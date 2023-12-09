package main

import (
	"bufio"
	"fmt"
	"math"
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
	waysQuad := 1
	for i, t := range times {
		dist := distances[i]
		ways *= calculateWins(t, dist)
		waysQuad *= calculateWinsQuadratic(t, dist)
	}
	fmt.Println(ways, waysQuad)

	fmt.Println(t2, d2)
	fmt.Println(calculateWins(t2, d2), calculateWinsQuadratic(t2, d2))
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

func calculateWinsQuadratic(time, distance int) int {
	// distance  = (hold * time - hold ^ 2)
	// 0 = -hold^2 + hold * time - distance
	x1, x2 := quadratic(-1, time, -1*distance)

	return int(x2) - int(x1) + 1
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

func quadratic(a, b, c int) (float64, float64) {
	x1 := (float64(-1*b) + math.Sqrt(float64(b*b-4*a*c))) / float64(2*a)
	x2 := (float64(-1*b) - math.Sqrt(float64(b*b-4*a*c))) / float64(2*a)

	return math.Floor(x1 + 1.0), math.Ceil(x2 - 1.0)
	return x1, x2
}
