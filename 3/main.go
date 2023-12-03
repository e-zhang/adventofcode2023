package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

type Coord struct {
	row int
	col int
}

type Part struct {
	Coord
	end   int
	value int
}

type Symbol struct {
	Coord
	value rune
}

func main() {
	f, err := os.Open("input")
	if err != nil {
		panic(err)
	}

	scanner := bufio.NewScanner(f)

	var parts []Part
	var symbols []Symbol

	lines := 0
	for scanner.Scan() {
		line := scanner.Text()

		start := -1
		for i, c := range line {
			isNumber := c >= '0' && c <= '9'
			if isNumber {
				if start < 0 {
					start = i
				} else if i == len(line)-1 {
					n, err := strconv.Atoi(line[start:])
					if err != nil {
						panic(err)
					}

					parts = append(parts, Part{Coord{lines, start}, i, n})
					fmt.Println("found part", parts[len(parts)-1])
				}
				continue
			}

			if start >= 0 && start < i {
				n, err := strconv.Atoi(line[start:i])
				if err != nil {
					panic(err)
				}

				parts = append(parts, Part{Coord{lines, start}, i - 1, n})
				fmt.Println("found part", parts[len(parts)-1])
				start = -1
			}

			if c != '.' {
				symbols = append(symbols, Symbol{Coord{lines, i}, c})
				fmt.Println("found symbol", Coord{lines, i}, string(c))
			}
		}

		lines += 1
	}

	sum := 0
	ratios := 0
	found := map[Coord]struct{}{}
	for _, symbol := range symbols {
		var touching []int
		for _, part := range parts {
			// check if the part is adjacent to the symbol
			if part.row < symbol.row-1 || part.row > symbol.row+1 ||
				part.end < symbol.col-1 || part.col > symbol.col+1 {
				continue
			}

			fmt.Printf("found part %#v touching symbol %#v\n", part, symbol)
			if _, ok := found[part.Coord]; !ok {
				sum += part.value
				found[part.Coord] = struct{}{}
			}

			touching = append(touching, part.value)
		}

		if symbol.value == '*' && len(touching) == 2 {
			ratios += touching[0] * touching[1]
		}
	}

	fmt.Println(sum)
	fmt.Println(ratios)
}
