package main

import (
	"bufio"
	"fmt"
	"os"
)

type Coord struct {
	row int
	col int
}

func (c *Coord) Add(o Coord) Coord {
	return Coord{c.row + o.row, c.col + o.col}
}

func main() {
	f, err := os.Open("input")
	if err != nil {
		panic(err)
	}

	var grid []string
	var start Coord
	total := 0
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		for c, s := range line {
			if s == 'S' {
				start = Coord{len(grid), c}
			}

			if s == '#' {
				total++
			}
		}
		grid = append(grid, line)
	}

	fmt.Println(move(grid, start, 64))

	rem := 26501365 % len(grid)
	for i := 0; i < 3; i++ {
		plots := move(grid, start, i*len(grid)+rem)
		fmt.Println(i, plots)
	}
}

type State struct {
	pos  Coord
	step int
}

func move(grid []string, start Coord, steps int) int {
	q := []State{State{start, 0}}

	plots := map[Coord]int{}

	for len(q) > 0 {
		curr := q[0]
		q = q[1:]

		if _, ok := plots[curr.pos]; ok {
			continue
		}

		plots[curr.pos] = curr.step

		if curr.step == steps {
			continue
		}

		for _, d := range []Coord{
			Coord{-1, 0},
			Coord{1, 0},
			Coord{0, -1},
			Coord{0, 1},
		} {
			next := curr.pos.Add(d)

			r := (next.row % len(grid))
			if r < 0 {
				r += len(grid)
			}
			c := next.col % len(grid[0])
			if c < 0 {
				c += len(grid[0])
			}

			if grid[r][c] == '#' {
				continue
			}

			q = append(q, State{next, curr.step + 1})
		}
	}

	// printGrid(grid, plots)
	// fmt.Println(plots)
	count := 0
	filtered := map[Coord]struct{}{}
	for k, v := range plots {
		if v%2 == steps%2 {
			// fmt.Println(k, x)
			filtered[k] = struct{}{}
			count++
		}
	}

	// fmt.Println(count, len(plots), steps)
	// fmt.Println(plots)
	// fmt.Println(count)
	// printGrid(grid, filtered)
	return count
}

func printGrid(g []string, o map[Coord]struct{}) {
	minR, maxR, minC, maxC := 0, 0, 0, 0
	for c := range o {
		if minR > c.row || minR == 0 {
			minR = c.row
		}
		if maxR < c.row || maxR == 0 {
			maxR = c.row
		}
		if minC > c.col || minC == 0 {
			minC = c.col
		}
		if maxC < c.col || maxC == 0 {
			maxC = c.col
		}
	}

	rows := maxR - minR
	cols := maxC - minC
	x := (rows)/len(g) + 1
	if y := (cols)/len(g[0]) + 1; y > x {
		x = y
	}

	roff := -x / 2 * len(g)
	coff := -x / 2 * len(g[0])

	fmt.Println(rows, cols, maxR, minR, maxC, minC, x)

	for i := 0; i < x*len(g); i++ {
		for j := 0; j < x*len(g[0]); j++ {
			if _, ok := o[Coord{i + roff, j + coff}]; ok {
				fmt.Printf("O")
			} else {
				r := (i % len(g))
				if r < 0 {
					r += len(g)
				}
				c := j % len(g[0])
				if c < 0 {
					c += len(g[0])
				}

				if g[r][c] == 'S' && (r != x/2*len(g) || c != x/2*len(g[0])) {
					fmt.Printf(".")
				} else {
					fmt.Printf(string(g[r][c]))
				}
			}
		}
		fmt.Println()
	}
}
