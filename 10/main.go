package main

import (
	"bufio"
	"fmt"
	"os"
)

var DEBUG = true

func main() {
	f, err := os.Open("test6")
	if err != nil {
		panic(err)
	}

	var grid [][]PipeType
	var start Coord
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		row := make([]PipeType, len(line))

		for i, c := range line {
			row[i] = PipeType(c)
			if row[i] == START {
				start = Coord{len(grid), i}
			}
		}

		grid = append(grid, row)
	}

	fmt.Println(start)
	printGrid(grid)

	loop := map[Coord]struct{}{
		start: struct{}{},
	}

	// check neighbors
	starters := start.Connections(grid)
	if len(starters) != 2 {
		panic(starters)
	}
	front := Node{curr: starters[0], prev: start}
	back := Node{curr: starters[1], prev: start}

	steps := 1
	frontEdges := 0
	backEdges := 0
	for front.curr != back.curr {
		loop[front.curr] = struct{}{}
		loop[back.curr] = struct{}{}

		front = move(front, grid)
		frontEdges += (front.curr.row - front.prev.row) * (front.curr.col + front.prev.col)
		backEdges += (back.curr.row - back.prev.row) * (back.curr.col + back.prev.col)
		back = move(back, grid)
		steps++
	}
	loop[front.curr] = struct{}{}

	fmt.Println(steps)
	var next Coord
	if frontEdges >= backEdges {
		next = starters[0]
	} else {
		next = starters[1]
	}

	fmt.Println(frontEdges, starters[0])
	fmt.Println(backEdges, starters[1])
	fill(start, next, loop, grid)
	raycasting(loop, grid)
}

func raycasting(loop map[Coord]struct{}, grid [][]PipeType) {
	count := 0
	for r := range grid {
		for c := range grid[0] {
			if _, ok := loop[Coord{r, c}]; ok {
				continue
			}

			odd := false
			previousBend := GROUND
			for i := -1; i < c; i++ {
				if _, ok := loop[Coord{r, i}]; ok {
					p := grid[r][i]

					if p == VERTICAL {
						odd = !odd
						continue
					}

					if p == HORIZONTAL {
						continue
					}

					switch previousBend {
					case GROUND:
						if p == SE_BEND || p == NE_BEND {
							previousBend = p
						}
					case SE_BEND:
						// non matching pair
						if p == NW_BEND {
							odd = !odd
							previousBend = GROUND
						}
						// matching pair
						if p == SW_BEND {
							previousBend = GROUND
						}
					case NE_BEND:
						if p == SW_BEND {
							odd = !odd
							previousBend = GROUND
						}
						if p == NW_BEND {
							previousBend = GROUND
						}
					case START:
					default:
						panic(string(previousBend))
					}
				}
			}

			// left over unmatched
			if previousBend != GROUND {
				odd = !odd
			}

			if odd {
				count++
			}
		}
	}

	fmt.Println(count)
}

func fill(start, next Coord, loop map[Coord]struct{}, grid [][]PipeType) {
	n := Node{curr: next, prev: start}
	q := traceLoop(n, loop, grid)

	seen := map[Coord]struct{}{}
	for len(q) != 0 {
		c := q[0]
		q = q[1:]

		seen[c] = struct{}{}

		if DEBUG {
			grid[c.row][c.col] = INSIDE
		}

		for _, d := range []Coord{
			Coord{c.row, c.col - 1}, // W
			Coord{c.row - 1, c.col}, // N
			Coord{c.row, c.col + 1}, // E
			Coord{c.row + 1, c.col}, // S
		} {
			if _, ok := seen[d]; ok {
				continue
			}

			if d.row < 0 || d.row >= len(grid) {
				continue
			}

			if d.col < 0 || d.col >= len(grid[0]) {
				continue
			}

			if _, ok := loop[d]; ok {
				continue
			}

			q = append(q, d)
		}
	}

	if DEBUG {
		printGrid(grid)
	}

	fmt.Println(len(seen))
}

func traceLoop(start Node, loop map[Coord]struct{}, grid [][]PipeType) []Coord {
	q := []Coord{}

	n := start

	for {
		n = move(n, grid)

		dir := Coord{n.curr.row - n.prev.row, n.curr.col - n.prev.col}

		switch grid[n.curr.row][n.curr.col] {
		case HORIZONTAL:
			adj := Coord{n.curr.row - dir.row, n.curr.col}
			_, ok := loop[adj]
			if adj.row >= 0 && adj.row < len(grid) &&
				adj.col >= 0 && adj.col < len(grid[adj.row]) && !ok {
				q = append(q, adj)

			}
		case VERTICAL:
			adj := Coord{n.curr.row, n.curr.col - dir.row}
			_, ok := loop[adj]
			if adj.row >= 0 && adj.row < len(grid) &&
				adj.col >= 0 && adj.col < len(grid[adj.row]) && !ok {
				q = append(q, adj)
			}
		case NE_BEND:
			for _, adj := range []Coord{
				{n.curr.row, n.curr.col - dir.row},
				{n.curr.row + dir.row, n.curr.col},
			} {
				if adj.row < 0 || adj.row >= len(grid) {
					continue
				}
				if adj.col < 0 || adj.col >= len(grid[adj.row]) {
					continue
				}

				if _, ok := loop[adj]; ok {
					continue
				}

				q = append(q, adj)
			}
		case NW_BEND:
			for _, adj := range []Coord{
				{n.curr.row + dir.col, n.curr.col},
				{n.curr.row, n.curr.col + dir.col},
			} {
				if adj.row < 0 || adj.row >= len(grid) {
					continue
				}
				if adj.col < 0 || adj.col >= len(grid[adj.row]) {
					continue
				}
				if _, ok := loop[adj]; ok {
					continue
				}

				q = append(q, adj)
			}
		case SE_BEND:
			for _, adj := range []Coord{
				{n.curr.row - dir.col, n.curr.col},
				{n.curr.row, n.curr.col + dir.col},
			} {
				if adj.row < 0 || adj.row >= len(grid) {
					continue
				}
				if adj.col < 0 || adj.col >= len(grid[adj.row]) {
					continue
				}
				if _, ok := loop[adj]; ok {
					continue
				}

				q = append(q, adj)
			}
		case SW_BEND:
			for _, adj := range []Coord{
				{n.curr.row + dir.row, n.curr.col},
				{n.curr.row, n.curr.col - dir.row},
			} {
				if adj.row < 0 || adj.row >= len(grid) {
					continue
				}
				if adj.col < 0 || adj.col >= len(grid[adj.row]) {
					continue
				}
				if _, ok := loop[adj]; ok {
					continue
				}

				q = append(q, adj)
			}
		}

		if n.curr == start.prev {
			break
		}
	}
	printGrid(grid)

	return q
}

func move(n Node, grid [][]PipeType) Node {
	next := n.curr.Connections(grid)
	if len(next) != 2 {
		fmt.Println(n.curr, next)
		panic(next)
	}

	curr := n.curr
	if next[0] != n.prev {
		n.curr = next[0]
	} else {
		n.curr = next[1]
	}
	n.prev = curr
	return n
}

func printGrid(g [][]PipeType) {
	fmt.Println()
	for _, r := range g {
		for _, c := range r {
			fmt.Printf("%s", string(c))
		}
		fmt.Println()
	}
	fmt.Println()
}

type Node struct {
	curr Coord
	prev Coord
}

type Coord struct {
	row int
	col int
}

func (c *Coord) Connections(grid [][]PipeType) []Coord {
	var conn []Coord
	for _, d := range []Coord{
		Coord{c.row, c.col - 1}, // W
		Coord{c.row - 1, c.col}, // N
		Coord{c.row, c.col + 1}, // E
		Coord{c.row + 1, c.col}, // S
	} {
		if d.row < 0 || d.row >= len(grid) {
			continue
		}

		if d.col < 0 || d.col >= len(grid[0]) {
			continue
		}

		if c.IsConnected(d, grid[d.row][d.col]) && d.IsConnected(*c, grid[c.row][c.col]) {
			conn = append(conn, d)
		}
	}

	return conn
}

func (c *Coord) IsConnected(other Coord, pipe PipeType) bool {
	switch pipe {
	case VERTICAL:
		return c.IsNorth(other) || c.IsSouth(other)
	case HORIZONTAL:
		return c.IsEast(other) || c.IsWest(other)
	case NE_BEND:
		return c.IsEast(other) || c.IsNorth(other)
	case NW_BEND:
		return c.IsWest(other) || c.IsNorth(other)
	case SW_BEND:
		return c.IsWest(other) || c.IsSouth(other)
	case SE_BEND:
		return c.IsEast(other) || c.IsSouth(other)
	case GROUND:
		return false
	case START:
		return true
	}
	panic(string(pipe))
}

func (c *Coord) IsNorth(other Coord) bool {
	return c.row == other.row-1 && c.col == other.col
}

func (c *Coord) IsSouth(other Coord) bool {
	return c.row == other.row+1 && c.col == other.col
}

func (c *Coord) IsEast(other Coord) bool {
	return c.row == other.row && c.col == other.col+1
}

func (c *Coord) IsWest(other Coord) bool {
	return c.row == other.row && c.col == other.col-1
}

type PipeType rune

const (
	VERTICAL   PipeType = '|'
	HORIZONTAL PipeType = '-'
	NE_BEND    PipeType = 'L'
	NW_BEND    PipeType = 'J'
	SW_BEND    PipeType = '7'
	SE_BEND    PipeType = 'F'
	GROUND     PipeType = '.'
	START      PipeType = 'S'
	INSIDE     PipeType = 'I'
)
