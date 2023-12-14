package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

const CYCLES = 1000000000

func main() {
	f, err := os.Open("input")
	if err != nil {
		panic(err)
	}

	grid := [][]rune{}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()

		row := []rune{}
		for _, s := range line {
			row = append(row, s)
		}
		grid = append(grid, row)
	}
	tiltN(grid)

	load := 0
	for i, r := range grid {
		for _, c := range r {
			if c == 'O' {
				load += len(grid) - i
			}
		}
	}
	fmt.Println(load)

	seen := map[string]int{}
	start := 0
	l := 0
	for i := 0; i < CYCLES; i++ {
		cycle(grid)
		// printGrid(grid)
		k := toString(grid)
		if v, ok := seen[k]; ok {
			fmt.Println(i, v)
			start = v
			l = i - v
			break
		}
		seen[k] = i
	}

	idx := ((CYCLES - 1 - start) % l)
	fmt.Println(idx)
	var g []string
	for k, v := range seen {
		if v == idx+start {
			g = strings.Split(k, ",")
			break
		}
	}

	load = 0
	for i, r := range g {
		for _, c := range r {
			if c == 'O' {
				load += len(g) - i
			}
		}
	}
	fmt.Println(load)
}

func toString(grid [][]rune) string {
	s := ""
	for i, r := range grid {
		if i > 0 {

		}
		s += ","
		for _, c := range r {
			s += string(c)
		}
	}
	return s
}

func cycle(grid [][]rune) {
	tiltN(grid)
	tiltW(grid)
	tiltS(grid)
	tiltE(grid)
	// printGrid(grid)
}

func printGrid(grid [][]rune) {
	for _, r := range grid {
		for _, c := range r {
			fmt.Printf("%s", string(c))
		}
		fmt.Println()
	}
	fmt.Println()
}

func tiltN(grid [][]rune) {
	for i := 1; i < len(grid); i++ {
		for j := 0; j < len(grid[0]); j++ {
			if grid[i][j] == 'O' {
				grid[i][j] = '.'
				for r := i - 1; r >= -1; r-- {
					if r == -1 || grid[r][j] != '.' {
						grid[r+1][j] = 'O'
						break
					}
				}
			}
		}
	}
}

func tiltS(grid [][]rune) {
	for i := len(grid) - 2; i >= 0; i-- {
		for j := 0; j < len(grid[0]); j++ {
			if grid[i][j] == 'O' {
				grid[i][j] = '.'
				for r := i + 1; r <= len(grid); r++ {
					if r == len(grid) || grid[r][j] != '.' {
						grid[r-1][j] = 'O'
						break
					}
				}
			}
		}
	}
}

func tiltW(grid [][]rune) {
	for i := 0; i < len(grid); i++ {
		for j := 1; j < len(grid[0]); j++ {
			if grid[i][j] == 'O' {
				grid[i][j] = '.'
				for c := j - 1; c >= -1; c-- {
					if c == -1 || grid[i][c] != '.' {
						grid[i][c+1] = 'O'
						break
					}
				}
			}
		}
	}
}

func tiltE(grid [][]rune) {
	for i := 0; i < len(grid); i++ {
		for j := len(grid[0]) - 2; j >= 0; j-- {
			if grid[i][j] == 'O' {
				grid[i][j] = '.'
				for c := j + 1; c <= len(grid[0]); c++ {
					if c == len(grid[0]) || grid[i][c] != '.' {
						grid[i][c-1] = 'O'
						break
					}
				}
			}
		}
	}
}
