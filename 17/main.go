package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
	"strconv"
)

const (
	UP int = iota
	RIGHT
	DOWN
	LEFT
	NONE
)

type Coord struct {
	row int
	col int
}

func (c *Coord) Add(o Coord) Coord {
	return Coord{c.row + o.row, c.col + o.col}
}

type State struct {
	loc   Coord
	dir   int
	count int
}

func main() {
	f, err := os.Open("input")
	if err != nil {
		panic(err)
	}

	grid := [][]int{}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()

		row := []int{}
		for _, s := range line {
			v, err := strconv.Atoi(string(s))
			if err != nil {
				panic(err)
			}
			row = append(row, v)
		}

		grid = append(grid, row)
	}

	fmt.Println(len(grid), len(grid[0]))

	fmt.Println(shortestPath(grid, 1, 3))
	fmt.Println(shortestPath(grid, 4, 10))
}

func getNextMoves(curr *State) []State {
	switch curr.dir {
	case RIGHT:
		return []State{
			State{curr.loc.Add(Coord{0, 1}), RIGHT, curr.count + 1},
			State{curr.loc.Add(Coord{-1, 0}), UP, 1},
			State{curr.loc.Add(Coord{1, 0}), DOWN, 1},
		}
	case DOWN:
		return []State{
			State{curr.loc.Add(Coord{1, 0}), DOWN, curr.count + 1},
			State{curr.loc.Add(Coord{0, 1}), RIGHT, 1},
			State{curr.loc.Add(Coord{0, -1}), LEFT, 1},
		}
	case LEFT:
		return []State{
			State{curr.loc.Add(Coord{0, -1}), LEFT, curr.count + 1},
			State{curr.loc.Add(Coord{-1, 0}), UP, 1},
			State{curr.loc.Add(Coord{1, 0}), DOWN, 1},
		}
	case UP:
		return []State{
			State{curr.loc.Add(Coord{-1, 0}), UP, curr.count + 1},
			State{curr.loc.Add(Coord{0, 1}), RIGHT, 1},
			State{curr.loc.Add(Coord{0, -1}), LEFT, 1},
		}
	}
	panic(curr.dir)
	return nil
}

func shortestPath(grid [][]int, minStep, maxStep int) int {
	start := Coord{0, 0}
	end := Coord{len(grid) - 1, len(grid[0]) - 1}
	q := PriorityQueue{
		&Item{State{start, DOWN, 1}, 0, 0},
		&Item{State{start, RIGHT, 1}, 0, 1},
	}
	heap.Init(&q)
	cost := map[State]int{
		q[0].state: 0,
	}
	for len(q) > 0 {
		item := heap.Pop(&q).(*Item)
		node := item.state

		if node.loc == end {
			return item.heat
		}

		next := getNextMoves(&node)
		for _, n := range next {
			if n.loc.row < 0 || n.loc.row >= len(grid) {
				continue
			}

			if n.loc.col < 0 || n.loc.col >= len(grid[0]) {
				continue
			}

			if n.count > maxStep {
				continue
			}

			if n.dir != node.dir && node.count < minStep {
				continue
			}

			if n.loc == end && n.count < minStep {
				continue
			}

			heat := item.heat + grid[n.loc.row][n.loc.col]
			h, ok := cost[n]
			if !ok || heat < h {

				cost[n] = heat

				heap.Push(&q, &Item{state: n, heat: heat})
			}
		}
	}

	panic("never reached target")
	return -1
}
