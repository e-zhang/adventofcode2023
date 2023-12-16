package main

import (
	"bufio"
	"fmt"
	"os"
)

const DEBUG = false

type Direction int

const (
	UP Direction = iota
	DOWN
	LEFT
	RIGHT
)

type Coord struct {
	row int
	col int
}

func (c *Coord) Add(o Coord) Coord {
	return Coord{c.row + o.row, c.col + o.col}
}

type Beam struct {
	loc Coord
	dir Direction
}

func (b *Beam) Move(grid [][]rune) bool {
	var d Coord
	switch b.dir {
	case UP:
		d = Coord{-1, 0}
	case DOWN:
		d = Coord{1, 0}
	case LEFT:
		d = Coord{0, -1}
	case RIGHT:
		d = Coord{0, 1}
	}

	next := b.loc.Add(d)
	if next.row < 0 || next.row >= len(grid) {
		return false
	}

	if next.col < 0 || next.col >= len(grid[0]) {
		return false
	}

	b.loc = next
	return true
}

func doBeams(grid [][]rune, start Beam) int {
	beams := []Beam{start}

	energized := map[Coord]struct{}{
		beams[0].loc: struct{}{},
	}

	seen := map[Beam]struct{}{
		beams[0]: struct{}{},
	}

	for len(beams) > 0 {
		beam := beams[0]
		beams = beams[1:]
		for beam.Move(grid) {
			if _, ok := seen[beam]; ok {
				break
			}
			energized[beam.loc] = struct{}{}
			seen[beam] = struct{}{}

			switch grid[beam.loc.row][beam.loc.col] {
			case '.':
			case '/':
				switch beam.dir {
				case RIGHT:
					beam.dir = UP
				case DOWN:
					beam.dir = LEFT
				case UP:
					beam.dir = RIGHT
				case LEFT:
					beam.dir = DOWN
				}
			case '\\':
				switch beam.dir {
				case RIGHT:
					beam.dir = DOWN
				case DOWN:
					beam.dir = RIGHT
				case UP:
					beam.dir = LEFT
				case LEFT:
					beam.dir = UP
				}
			case '|':
				switch beam.dir {
				case LEFT, RIGHT:
					beam.dir = UP
					split := Beam{Coord{beam.loc.row, beam.loc.col}, DOWN}
					beams = append(beams, split)
				}
			case '-':
				switch beam.dir {
				case UP, DOWN:
					beam.dir = RIGHT
					split := Beam{Coord{beam.loc.row, beam.loc.col}, LEFT}
					beams = append(beams, split)
				}
			}
		}
	}

	if DEBUG {
		fmt.Println(len(energized) - 1)
		printGrid(grid, energized)
		fmt.Println()
	}
	return len(energized) - 1
}

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

	fmt.Println(doBeams(grid, Beam{Coord{0, -1}, RIGHT}))
	fmt.Println(part2(grid))
}

func part2(grid [][]rune) int {
	max := 0

	for r := range grid {
		start := Beam{Coord{r, -1}, RIGHT}
		energized := doBeams(grid, start)
		if energized > max {
			max = energized
		}

		start = Beam{Coord{r, len(grid[0])}, LEFT}
		energized = doBeams(grid, start)
		if energized > max {
			max = energized
		}
	}

	for c := range grid[0] {
		start := Beam{Coord{-1, c}, DOWN}
		energized := doBeams(grid, start)
		if energized > max {
			max = energized
		}

		start = Beam{Coord{len(grid), c}, UP}
		energized = doBeams(grid, start)
		if energized > max {
			max = energized
		}
	}

	return max
}

func printGrid(g [][]rune, energized map[Coord]struct{}) {
	for i, r := range g {
		for j, c := range r {
			if _, ok := energized[Coord{i, j}]; ok {
				fmt.Printf("#")
			} else {
				fmt.Printf(string(c))
			}
		}
		fmt.Println()
	}
}
