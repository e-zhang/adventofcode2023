package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

type MapType int

const (
	NONE MapType = iota
	SEED_TO_SOIL
	SOIL_TO_FERTILIZER
	FERTILIZER_TO_WATER
	WATER_TO_LIGHT
	LIGHT_TO_TEMPERATURE
	TEMPERATURE_TO_HUMIDITY
	HUMIDITY_TO_LOCATION
)

type Lookup struct {
	dst    int
	src    int
	length int
}

func main() {
	f, err := os.Open("input")
	if err != nil {
		panic(err)
	}

	var seeds []int
	typ := NONE

	maps := map[MapType][]Lookup{
		SEED_TO_SOIL:            []Lookup{},
		SOIL_TO_FERTILIZER:      []Lookup{},
		FERTILIZER_TO_WATER:     []Lookup{},
		WATER_TO_LIGHT:          []Lookup{},
		LIGHT_TO_TEMPERATURE:    []Lookup{},
		TEMPERATURE_TO_HUMIDITY: []Lookup{},
		HUMIDITY_TO_LOCATION:    []Lookup{},
	}

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			typ = NONE
			continue
		}

		switch typ {
		case NONE:
			switch {
			case strings.HasPrefix(line, "seeds:"):
				seeds = parseSeeds(line)
			case strings.HasPrefix(line, "seed-to-soil"):
				typ = SEED_TO_SOIL
			case strings.HasPrefix(line, "soil-to-fertilizer"):
				typ = SOIL_TO_FERTILIZER
			case strings.HasPrefix(line, "fertilizer-to-water"):
				typ = FERTILIZER_TO_WATER
			case strings.HasPrefix(line, "water-to-light"):
				typ = WATER_TO_LIGHT
			case strings.HasPrefix(line, "light-to-temperature"):
				typ = LIGHT_TO_TEMPERATURE
			case strings.HasPrefix(line, "temperature-to-humidity"):
				typ = TEMPERATURE_TO_HUMIDITY
			case strings.HasPrefix(line, "humidity-to-location"):
				typ = HUMIDITY_TO_LOCATION
			}
		default:
			maps[typ] = parseMap(maps[typ], line)
		}
	}

	fmt.Println(part1(seeds, maps))
	fmt.Println(part2(seeds, maps))
}

type Interval struct {
	start int
	end   int
}

func compareIntervals(a, b Interval) []Interval {
	var out []Interval

	before := Interval{a.start, min(b.start, a.end)}
	if before.end > before.start {
		out = append(out, before)
	}

	after := Interval{max(a.start, b.end), a.end}
	if after.end > after.start {
		out = append(out, after)
	}

	middle := Interval{max(a.start, b.start), min(a.end, b.end)}
	if middle.end > middle.start {
		out = append(out, middle)
	}

	return out
}

func part2(seeds []int, maps map[MapType][]Lookup) int {
	min := math.MaxInt
	for i := 0; i < len(seeds); i += 2 {
		start := seeds[i]
		l := seeds[i+1]

		intervals := []Interval{{start, start + l}}
		for _, typ := range []MapType{
			SEED_TO_SOIL,
			SOIL_TO_FERTILIZER,
			FERTILIZER_TO_WATER,
			WATER_TO_LIGHT,
			LIGHT_TO_TEMPERATURE,
			TEMPERATURE_TO_HUMIDITY,
			HUMIDITY_TO_LOCATION,
		} {
			var next []Interval
			for _, lk := range maps[typ] {
				remaining := intervals[:0]
				for _, iv := range intervals {
					segments := compareIntervals(iv, Interval{lk.src, lk.src + lk.length})
					for _, seg := range segments {
						if seg.start >= lk.src && seg.end <= lk.src+lk.length {
							shift := lk.dst - lk.src
							next = append(next, Interval{seg.start + shift, seg.end + shift})
						} else {
							remaining = append(remaining, seg)
						}
					}
				}
				intervals = remaining
			}
			next = append(next, intervals...)
			intervals = next
		}

		for _, iv := range intervals {
			if iv.start < min {
				min = iv.start
			}
		}
	}

	return min
}

func part1(seeds []int, maps map[MapType][]Lookup) int {
	min := math.MaxInt
	for _, seed := range seeds {
		key := lookupSeed(seed, maps)
		if key < min {
			min = key
		}
	}

	return min
}

func lookupSeed(seed int, maps map[MapType][]Lookup) int {
	key := seed
	for _, typ := range []MapType{
		SEED_TO_SOIL,
		SOIL_TO_FERTILIZER,
		FERTILIZER_TO_WATER,
		WATER_TO_LIGHT,
		LIGHT_TO_TEMPERATURE,
		TEMPERATURE_TO_HUMIDITY,
		HUMIDITY_TO_LOCATION,
	} {
		for _, entry := range maps[typ] {
			if key >= entry.src && key < entry.src+entry.length {
				key = entry.dst + (key - entry.src)
				break
			}
		}
	}

	return key
}

func parseMap(m []Lookup, s string) []Lookup {
	var src, dst, l int
	fmt.Sscanf(s, "%d %d %d", &dst, &src, &l)
	return append(m, Lookup{dst, src, l})
	return m
}

func parseSeeds(s string) []int {
	v := strings.Split(s, ":")
	seedStrings := strings.Split(v[1], " ")

	var seeds []int
	for _, seed := range seedStrings {
		if seed == "" {
			continue
		}
		a, _ := strconv.Atoi(seed)
		seeds = append(seeds, a)
	}

	return seeds
}
