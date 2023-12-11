package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

const (
	EMPTY  = '.'
	GALAXY = '#'
)

const MULTIPLIER = 1000000

func main() {
	f, err := os.Open("input")
	if err != nil {
		panic(err)
	}

	scanner := bufio.NewScanner(f)
	var grid []string
	for scanner.Scan() {
		line := scanner.Text()

		grid = append(grid, line)
	}

	g1 := make([]string, len(grid))
	copy(g1, grid)
	part1Naive(g1)
	fmt.Println(calculateDistances(grid, 2))
	fmt.Println(calculateDistances(grid, MULTIPLIER))
}

func abs(x int) int {
	if x < 0 {
		return -1 * x
	}
	return x
}

// isBetween returns if x is between y and z
func isBetween(x, y, z int) bool {
	return x > y && x < z || x < y && x > z
}

func calculateDistances(grid []string, multiplier int) int {
	emptyRows, emptyCols := expand(grid)
	galaxies := findGalaxies(grid)

	paths := 0
	for i, g := range galaxies {
		for _, g2 := range galaxies[i+1:] {
			rows := 0
			for _, r := range emptyRows {
				if isBetween(r, g.row, g2.row) {
					rows++
				}
			}

			cols := 0
			for _, c := range emptyCols {
				if isBetween(c, g.col, g2.col) {
					cols++
				}
			}

			d := abs(g.row-g2.row) + rows*(multiplier-1) + abs(g.col-g2.col) + cols*(multiplier-1)
			paths += d
		}
	}

	return paths
}

func part1Naive(grid []string) {
	grid = expandPhysically(grid)

	galaxies := findGalaxies(grid)

	paths := 0
	for i, g := range galaxies {
		for _, g2 := range galaxies[i+1:] {
			d := abs(g.row-g2.row) + abs(g.col-g2.col)
			// fmt.Printf("distance from %v to %v: %d\n", g, g2, d)
			paths += d
		}
	}

	fmt.Println(paths)
}

func printGrid(grid []string) {
	fmt.Println()
	for _, r := range grid {
		fmt.Println(r)
	}
	fmt.Println()
}

type Coord struct {
	row int
	col int
}

func findGalaxies(grid []string) []Coord {
	var galaxies []Coord
	for i, r := range grid {
		for j, c := range r {
			if c == GALAXY {
				galaxies = append(galaxies, Coord{i, j})
			}
		}
	}

	return galaxies
}

func expand(grid []string) ([]int, []int) {
	var emptyRows []int
	for i, r := range grid {
		if !strings.ContainsRune(r, GALAXY) {
			emptyRows = append(emptyRows, i)
		}
	}

	var emptyCols []int
	for j := 0; j < len(grid[0]); j++ {
		hasGalaxy := false
		for i := 0; i < len(grid); i++ {
			if grid[i][j] == GALAXY {
				hasGalaxy = true
				break
			}
		}

		if !hasGalaxy {
			emptyCols = append(emptyCols, j)
		}
	}

	return emptyRows, emptyCols
}

func expandPhysically(grid []string) []string {
	printGrid(grid)

	for i := 0; i < len(grid); i++ {
		if !strings.ContainsRune(grid[i], GALAXY) {
			grid = append(grid, "")
			copy(grid[i+1:], grid[i:])
			i++
		}
	}

	for j := 0; j < len(grid[0]); j++ {
		hasGalaxy := false
		for i := 0; i < len(grid); i++ {
			if grid[i][j] == GALAXY {
				hasGalaxy = true
				break
			}
		}

		if !hasGalaxy {
			for i := 0; i < len(grid); i++ {
				grid[i] = grid[i][:j] + string(EMPTY) + grid[i][j:]
			}
			j++
		}
	}

	printGrid(grid)
	return grid
}
