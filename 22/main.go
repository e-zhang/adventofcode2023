package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

type Coord struct {
	x int
	y int
	z int
}

func parseCoord(s string) Coord {
	c := strings.Split(s, ",")

	x, err := strconv.Atoi(c[0])
	if err != nil {
		panic(err)
	}
	y, err := strconv.Atoi(c[1])
	if err != nil {
		panic(err)
	}
	z, err := strconv.Atoi(c[2])
	if err != nil {
		panic(err)
	}

	return Coord{x, y, z}
}

type Brick struct {
	name string

	start Coord
	end   Coord
}

func (b *Brick) Fall() bool {
	if b.start.z > 1 {
		b.start.z--
		b.end.z--
		return true
	}

	return false
}

func (b *Brick) SupportedBy(o Brick) bool {
	if o.end.z != b.start.z-1 {
		return false
	}

	x := b.start.x <= o.end.x && b.end.x >= o.start.x
	y := b.start.y <= o.end.y && b.end.y >= o.start.y
	return x && y
}

func parseBrick(idx int, line string) Brick {
	s := strings.Split(line, "~")

	start := parseCoord(s[0])
	end := parseCoord(s[1])

	if end.z < start.z || end.x < start.x || end.y < start.y {
		panic(line)
	}

	name := string('A' + idx)
	// name := strconv.Itoa(idx)
	return Brick{name, start, end}
}

func main() {
	f, err := os.Open("input")
	if err != nil {
		panic(err)
	}

	var bricks []Brick
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		b := parseBrick(len(bricks), line)
		bricks = append(bricks, b)
	}

	sort.Slice(bricks, func(i, j int) bool {
		return bricks[i].start.z < bricks[j].start.z
	})

	invariant(bricks)
	bricks, supports, supporters := settle(bricks)
	invariant(bricks)

	part1(bricks, supports, supporters)
	part2(bricks, supports, supporters)
}

func invariant(bricks []Brick) {
	for _, b := range bricks {
		for _, o := range bricks {
			if b.name == o.name {
				continue
			}
			if b.start.z == o.start.z {
				x := b.start.x <= o.end.x && b.end.x >= o.start.x
				y := b.start.y <= o.end.y && b.end.y >= o.start.y
				if x && y {
					fmt.Println(b.name, o.name)
					panic(fmt.Sprintf("%t %t", x, y))
				}
			}
		}
	}
}

func settle(bricks []Brick) ([]Brick, map[string][]string, map[string][]string) {
	supporters := map[string][]string{}
	supports := map[string][]string{}

	for i, b := range bricks {
		canFall := true
		for canFall {
			for j := 0; j < i; j++ {
				other := bricks[j]
				if ok := b.SupportedBy(other); ok {
					canFall = false
					supporters[b.name] = append(supporters[b.name], other.name)
					supports[other.name] = append(supports[other.name], b.name)
				}
			}

			if canFall {
				canFall = b.Fall()
			}
		}
		bricks[i] = b
	}

	return bricks, supports, supporters
}

func part2(bricks []Brick, supports, supporters map[string][]string) {
	count := 0
	for _, b := range bricks {
		state := map[string]struct{}{}
		chain(b.name, state, supports, supporters)
		count += len(state) - 1
	}
	fmt.Println(count)
}

func chain(b string, state map[string]struct{}, supports, supporters map[string][]string) {
	state[b] = struct{}{}
	for _, n := range supports[b] {
		s := supporters[n]
		rem := len(s)
		for _, n := range s {
			if _, ok := state[n]; ok {
				rem--
			}
		}

		if rem > 0 {
			continue
		}

		chain(n, state, supports, supporters)
	}
}

func part1(bricks []Brick, supports, supporters map[string][]string) {
	disintegrates := map[string]struct{}{}
	for _, b := range bricks {
		v := supports[b.name]
		if len(v) == 0 {
			disintegrates[b.name] = struct{}{}
			continue
		}

		ok := true
		for _, n := range v {
			s := supporters[n]
			if len(s) == 1 {
				ok = false
				break
			}
		}
		if ok {
			disintegrates[b.name] = struct{}{}
		}
	}

	fmt.Println(len(disintegrates))
}
