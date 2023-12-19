package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const (
	UP    = "U"
	DOWN  = "D"
	LEFT  = "L"
	RIGHT = "R"
)

type Command struct {
	dir   string
	steps int
	color string
}

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

	curr1 := Coord{0, 0}
	border1 := map[Coord]struct{}{
		curr1: struct{}{},
	}
	test := []Coord{curr1}
	curr2 := Coord{0, 0}
	border2 := []Coord{curr2}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()

		cmd := parseCommand(line)

		var d Coord
		switch cmd.dir {
		case UP:
			d = Coord{-1, 0}
		case DOWN:
			d = Coord{1, 0}
		case RIGHT:
			d = Coord{0, 1}
		case LEFT:
			d = Coord{0, -1}
		default:
			panic(cmd.dir)
		}

		for i := 0; i < cmd.steps; i++ {
			curr1 = curr1.Add(d)
			border1[curr1] = struct{}{}
			test = append(test, curr1)
		}

		steps, err := strconv.ParseInt(cmd.color[:5], 16, 64)
		if err != nil {
			panic(cmd.color)
		}
		switch cmd.color[5] {
		case '0':
			d = Coord{0, 1}
		case '1':
			d = Coord{1, 0}
		case '2':
			d = Coord{0, -1}
		case '3':
			d = Coord{-1, 0}
		}

		for i := 0; i < int(steps); i++ {
			curr2 = curr2.Add(d)
			border2 = append(border2, curr2)
		}
	}

	if curr1.row != 0 && curr1.col != 0 {
		panic(curr1)
	}
	if curr2.row != 0 && curr2.col != 0 {
		panic(curr2)
	}
	fmt.Println(len(border1), len(border2))
	fill(border1)
	part2Fill(border1)
	// part2Range(border1)
	part2Shoelace(test)
	part2Shoelace(border2)
}

func parseCommand(line string) Command {
	var c Command
	fmt.Sscanf(line, "%s %d %s", &c.dir, &c.steps, &c.color)
	c.color = strings.Trim(c.color, "()#")
	return c
}

func getMinMax(border map[Coord]struct{}) (int, int, int, int) {
	minR, maxR, minC, maxC := 0, 0, 0, 0

	for k := range border {
		if k.row < minR {
			minR = k.row
		}
		if k.row > maxR {
			maxR = k.row
		}
		if k.col < minC {
			minC = k.col
		}
		if k.col > maxC {
			maxC = k.col
		}
	}

	fmt.Println(minR, maxR, minC, maxC)
	return minR, maxR, minC, maxC
}

func fill(border map[Coord]struct{}) int {
	minR, maxR, minC, maxC := getMinMax(border)

	start := Coord{minR + (maxR-minR)/2, minC + (maxC-minC)/2}
	seen := map[Coord]struct{}{}
	q := []Coord{start}
	for len(q) > 0 {
		v := q[0]
		q = q[1:]

		if _, ok := seen[v]; ok {
			continue
		}
		seen[v] = struct{}{}

		for _, n := range []Coord{
			Coord{0, 1},
			Coord{0, -1},
			Coord{1, 0},
			Coord{-1, 0},
		} {
			next := v.Add(n)
			if next.row < minR || next.row > maxR {
				continue
			}

			if next.col < minC || next.col > maxC {
				continue
			}

			if _, ok := border[next]; ok {
				continue
			}
			q = append(q, next)
		}
	}

	fmt.Println(len(seen) + len(border))
	return -1
}

func part2Shoelace(border []Coord) {
	sum := int64(0)
	for i := 1; i < len(border); i++ {
		pt1 := border[i-1]
		pt2 := border[i]

		determinant := (pt1.col * pt2.row) - (pt1.row * pt2.col)
		sum += int64(determinant)
	}

	fmt.Println(sum, sum/2+int64(len(border))/2+1)
}

func part2Fill(border map[Coord]struct{}) {
	minR, maxR, minC, maxC := getMinMax(border)

	start := Coord{minR - 1, minC - 1}
	fmt.Println(start)
	seen := map[Coord]struct{}{}
	q := []Coord{start}
	for len(q) > 0 {
		v := q[0]
		q = q[1:]

		if _, ok := seen[v]; ok {
			continue
		}
		seen[v] = struct{}{}

		for _, n := range []Coord{
			Coord{0, 1},
			Coord{0, -1},
			Coord{1, 0},
			Coord{-1, 0},
		} {
			next := v.Add(n)
			if next.row < minR-1 || next.row > maxR+1 {
				continue
			}

			if next.col < minC-1 || next.col > maxC+1 {
				continue
			}

			if _, ok := border[next]; ok {
				continue
			}
			q = append(q, next)
		}
	}

	fmt.Println(len(seen))
	rows := maxR - minR + 3
	cols := maxC - minC + 3
	fmt.Println(rows*cols - len(seen))
}
