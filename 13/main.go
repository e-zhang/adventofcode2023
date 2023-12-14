package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	f, err := os.Open("input")
	if err != nil {
		panic(err)
	}

	var grid [][]rune

	part1 := []int{0, 0}
	part2 := []int{0, 0}

	count := 0
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()

		if line == "" {
			c, r := checkReflection(grid, 0, 0)
			part1[0] += c
			part1[1] += r

			c, r = checkSmudge(grid)
			part2[0] += c
			part2[1] += r
			grid = [][]rune{}
			continue
		}

		var row []rune
		for _, s := range line {
			row = append(row, s)
		}
		grid = append(grid, row)
	}
	c, r := checkReflection(grid, 0, 0)
	part1[0] += c
	part1[1] += r

	c, r = checkSmudge(grid)
	part2[0] += c
	part2[1] += r

	fmt.Println(part1[1]*100 + part1[0])
	fmt.Println(part2[1]*100 + part2[0])
}

func printGrid(grid [][]rune) {
	for _, r := range grid {
		for _, c := range r {
			fmt.Printf(string(c))
		}
		fmt.Println()
	}
	fmt.Println()
}

func checkSmudge(grid [][]rune) (int, int) {
	oc, or := checkReflection(grid, 0, 0)

	for r := 0; r < len(grid); r++ {
		for c := 0; c < len(grid[r]); c++ {
			orig := grid[r][c]
			if grid[r][c] == '.' {
				grid[r][c] = '#'
			} else {
				grid[r][c] = '.'
			}
			col, row := checkReflection(grid, or, oc)

			if row+col > 0 {
				return col, row
			}
			grid[r][c] = orig
		}
	}

	return 0, 0
}

func checkReflection(grid [][]rune, ignoreRow, ignoreCol int) (int, int) {
	c := checkVerticalReflection(grid, ignoreCol)
	r := checkHorizontalReflection(grid, ignoreRow)
	return c, r
}

func checkHorizontalReflection(grid [][]rune, ignoreRow int) int {
	line := 0

	for r := 1; r < len(grid); r++ {
		reflected := true
		for i := 0; i < r; i++ {
			if !checkRow(grid, r, i) {
				reflected = false
				break
			}
		}

		if reflected && r != ignoreRow {
			line = r
			break
		}
	}

	return line
}

func checkRow(grid [][]rune, center, offset int) bool {
	if center+offset >= len(grid) {
		return true
	}

	if center-offset-1 < 0 {
		return true
	}

	for c := 0; c < len(grid[center]); c++ {
		if grid[center+offset][c] != grid[center-offset-1][c] {
			return false
		}
	}

	return true
}

func checkVerticalReflection(grid [][]rune, ignoreCol int) int {
	line := 0

	for c := 1; c < len(grid[0]); c++ {
		reflected := true
		for i := 0; i < c; i++ {
			if !checkCol(grid, c, i) {
				reflected = false
				break
			}
		}

		if reflected && c != ignoreCol {
			line = c
			break
		}
	}

	return line
}

func checkCol(grid [][]rune, center, offset int) bool {
	if center+offset >= len(grid[0]) {
		return true
	}

	if center-offset-1 < 0 {
		return true
	}

	for r := 0; r < len(grid); r++ {
		if grid[r][center+offset] != grid[r][center-offset-1] {
			return false
		}
	}

	return true
}
