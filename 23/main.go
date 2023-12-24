package main

import (
	"bufio"
	"fmt"
	"os"
)

const (
	PATH        = '.'
	FOREST      = '#'
	SLOPE_UP    = '^'
	SLOPE_DOWN  = 'v'
	SLOPE_RIGHT = '>'
	SLOPE_LEFT  = '<'
)

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

type Coord struct {
	row int
	col int
}

func (c *Coord) Add(d Coord) Coord {
	return Coord{c.row + d.row, c.col + d.col}
}
func (c *Coord) Count(d Coord) int {
	row := abs(d.row - c.row)
	col := abs(d.col - c.col)

	if row != 0 && col != 0 {
		panic(d)
	}

	return row + col
}

func main() {
	f, err := os.Open("input")
	if err != nil {
		panic(err)
	}

	var grid []string

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		grid = append(grid, line)
	}

	var start, end Coord
	for c := 0; c < len(grid[0]); c++ {
		if grid[0][c] == PATH {
			start = Coord{0, c}
		}
		if grid[len(grid)-1][c] == PATH {
			end = Coord{len(grid) - 1, c}
		}
	}

	fmt.Println(start, end)
	// part1(grid, start, end)
	part2(grid, start, end)
}

func part1(grid []string, start, end Coord) {
	seen := map[Coord]struct{}{start: struct{}{}}
	steps := dfs(grid, seen, start, end, 0, true)
	fmt.Println(steps)
}

type State struct {
	seen  map[Coord]struct{}
	steps int
}

func copyMap(m map[Coord]struct{}) map[Coord]struct{} {
	n := make(map[Coord]struct{})
	for k, v := range m {
		n[k] = v
	}
	return n
}

func part2(grid []string, start, end Coord) {
	seen := map[Coord]struct{}{start: struct{}{}}
	steps := dfs(grid, seen, start, end, 0, false)
	fmt.Println(steps)
}

func printGrid(grid []string, seen map[Coord]struct{}) {
	for i, r := range grid {
		for j, c := range r {
			if _, ok := seen[Coord{i, j}]; ok {
				fmt.Printf("O")
			} else {
				fmt.Printf(string(c))
			}
		}
		fmt.Println()
	}
}

func adjacency(grid []string, curr Coord) []Coord {
	next := moves(grid, curr, false)
	for i, n := range next {
		m := moves(grid, n, false)
		for {
			if len(m) != 2 {
				break
			}

			p := m[0]
			if p == n {
				p = m[1]
			}

			if !(p.col == n.col && n.col == curr.col ||
				p.row == n.row && n.row == curr.row) {
				break
			}

			m = moves(grid, p, false)
			n = p
		}

		next[i] = n
	}

	return next
}

func dfs(grid []string, seen map[Coord]struct{}, curr Coord, end Coord, steps int, slopes bool) int {
	if curr == end {
		// printGrid(grid, seen)
		// fmt.Println(steps, len(seen))
		return steps
	}

	max := 0
	for _, n := range adjacency(grid, curr) {
		if _, ok := seen[n]; ok {
			continue
		}

		seen[n] = struct{}{}
		s := dfs(grid, seen, n, end, steps+curr.Count(n), slopes)
		if s > max {
			max = s
		}
		delete(seen, n)
	}

	return max
}

func moves(grid []string, loc Coord, slopes bool) []Coord {
	if slopes {
		switch grid[loc.row][loc.col] {
		case SLOPE_DOWN:
			return []Coord{loc.Add(Coord{1, 0})}
		case SLOPE_UP:
			return []Coord{loc.Add(Coord{-1, 0})}
		case SLOPE_LEFT:
			return []Coord{loc.Add(Coord{0, -1})}
		case SLOPE_RIGHT:
			return []Coord{loc.Add(Coord{0, 1})}
		}
	}

	var next []Coord
	for _, n := range []Coord{
		Coord{1, 0},
		Coord{-1, 0},
		Coord{0, 1},
		Coord{0, -1},
	} {
		nxt := loc.Add(n)
		if nxt.row < 0 || nxt.row >= len(grid) {
			continue
		}

		if nxt.col < 0 || nxt.col >= len(grid[0]) {
			continue
		}

		if grid[nxt.row][nxt.col] == FOREST {
			continue
		}

		next = append(next, nxt)
	}

	return next
}
