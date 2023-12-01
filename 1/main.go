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

	sum := 0
	for scanner.Scan() {
		line := scanner.Text()
		val := part2(line)
		sum += val
		fmt.Println(val, sum)
	}

	fmt.Println(sum)
}

func part2(line string) int {
	var digits []int
	for i := 0; i < len(line); {
		s := line[i:]
		switch {
		case strings.HasPrefix(s, "one"):
			digits = append(digits, 1)
		case strings.HasPrefix(s, "two"):
			digits = append(digits, 2)
		case strings.HasPrefix(s, "three"):
			digits = append(digits, 3)
		case strings.HasPrefix(s, "four"):
			digits = append(digits, 4)
		case strings.HasPrefix(s, "five"):
			digits = append(digits, 5)
		case strings.HasPrefix(s, "six"):
			digits = append(digits, 6)
		case strings.HasPrefix(s, "seven"):
			digits = append(digits, 7)
		case strings.HasPrefix(s, "eight"):
			digits = append(digits, 8)
		case strings.HasPrefix(s, "nine"):
			digits = append(digits, 9)
		default:
			c := line[i]
			if c > '0' && c <= '9' {
				digits = append(digits, int(c-'0'))
			}
		}
		i += 1
	}

	fmt.Println(digits)
	return digits[0]*10 + digits[len(digits)-1]
}

func part1(line string) int {
	var digits []rune
	for _, c := range line {
		if c >= '0' && c <= '9' {
			digits = append(digits, c)
		}
	}

	val, err := strconv.Atoi(string(digits[0]) + string(digits[len(digits)-1]))
	if err != nil {
		panic(err)
	}

	return val
}
