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
	part1 := 0
	part2 := 0
	for scanner.Scan() {
		line := scanner.Text()

		historyStr := strings.Split(line, " ")

		var history []int
		for _, vStr := range historyStr {
			v, err := strconv.Atoi(vStr)
			if err != nil {
				panic(err)
			}
			history = append(history, v)
		}

		var lasts []int
		var firsts []int
		for !allZeros(history) {
			lasts = append(lasts, history[len(history)-1])
			firsts = append(firsts, history[0])
			history = difference(history)
		}

		inc := 0
		for i := len(lasts) - 1; i >= 0; i-- {
			inc += lasts[i]
		}
		part1 += inc
		// fmt.Println(inc)

		inc = 0
		for i := len(firsts) - 1; i >= 0; i-- {
			inc = firsts[i] - inc
		}
		part2 += inc
	}

	fmt.Println(part1)
	fmt.Println(part2)
}

func allZeros(d []int) bool {
	for _, v := range d {
		if v != 0 {
			return false
		}
	}

	return true
}

func difference(d []int) []int {
	var out []int

	for i := 1; i < len(d); i++ {
		out = append(out, d[i]-d[i-1])
	}

	return out
}
