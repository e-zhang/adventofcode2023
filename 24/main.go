package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"

	"github.com/davidkleiven/gononlin/nonlin"
	"gonum.org/v1/exp/linsolve"
)

const (
	// LOW  = 7
	// HIGH = 27
	LOW  = 200000000000000
	HIGH = 400000000000000
)

type Coord struct {
	x int
	y int
	z int
}

func parseCoord(line string) Coord {
	s := strings.Split(line, ",")

	x, err := strconv.Atoi(strings.TrimSpace(s[0]))
	if err != nil {
		panic(err)
	}
	y, err := strconv.Atoi(strings.TrimSpace(s[1]))
	if err != nil {
		panic(err)
	}
	z, err := strconv.Atoi(strings.TrimSpace(s[2]))
	if err != nil {
		panic(err)
	}

	return Coord{x, y, z}
}

type Hail struct {
	pos Coord
	v   Coord

	m float64
	b float64
}

func (h *Hail) solve() {
	h.m = float64(h.v.y) / float64(h.v.x)
	h.b = float64(h.pos.y) - h.m*float64(h.pos.x)
}

func (h *Hail) Intersect(o Hail) (float64, float64) {
	// m1x + b1 = m2x + b2
	// (m1 - m2)x = (b2 - b1)
	// x = (b2 - b1) / (m1 - m2)

	x := (o.b - h.b) / (h.m - o.m)
	y := h.m*x + h.b
	if (y-float64(h.pos.y))/float64(h.v.y) < 0 ||
		(x-float64(h.pos.x))/float64(h.v.x) < 0 {
		return math.Inf(0), math.Inf(0)
	}

	if (y-float64(o.pos.y))/float64(o.v.y) < 0 ||
		(x-float64(o.pos.x))/float64(o.v.x) < 0 {
		return math.Inf(0), math.Inf(0)
	}

	return x, y
}

func parseHail(line string) Hail {
	s := strings.Split(line, "@")

	p := parseCoord(s[0])
	v := parseCoord(s[1])

	h := Hail{p, v, 0, 0}
	h.solve()
	return h
}

func main() {
	f, err := os.Open("input")
	if err != nil {
		panic(err)
	}

	var hails []Hail
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		h := parseHail(line)
		hails = append(hails, h)
	}

	// part1(hails)
	part2(hails)
}

func part2(hails []Hail) {

	// x(t) = pos.x + v.x * t
	// y(t) = pos.y + v.y * t
	// z(t) = pos.z + v.z * t

	// L_x(t) = x_l + v_l.x * t
	// L_y(t) = y_l + v_l.y * t
	// L_z(t) = z_l + v_l.z * t

	// x_0(t_1) = L_x(t_2)
	// pos_0.x + v_0.x * t = x_l + v_l.x * t
	// x_l + (v_l.x - v_0.x) * t - pos_0.x = 0

	// for i := 0; i < 3; i++ {
	// 	h := hails[i]
	// 	fmt.Printf("x + (vx - %d) * t%d - %d==0\n", h.v.x, i+1, h.pos.x)
	// 	fmt.Printf("y + (vy - %d) * t%d - %d==0\n", h.v.y, i+1, h.pos.y)
	// 	fmt.Printf("z + (vz - %d) * t%d - %d==0\n", h.v.z, i+1, h.pos.z)

	// 	fmt.Printf("x + t%d*vx == %d + %d*t%d\n", i+1, h.pos.x, h.v.x, i+1)
	// 	fmt.Printf("y + t%d*vy == %d + %d*t%d\n", i+1, h.pos.y, h.v.y, i+1)
	// 	fmt.Printf("z + t%d*vz == %d + %d*t%d\n", i+1, h.pos.z, h.v.z, i+1)
	// }

	eqs := nonlin.Problem{
		F: func(out, x []float64) {
			for i := 0; i < 3; i++ {
				h := hails[i]
				out[i*3+0] = x[0] + (x[3]-float64(h.v.x))*x[6+i] - float64(h.pos.x)
				out[i*3+1] = x[1] + (x[4]-float64(h.v.y))*x[6+i] - float64(h.pos.y)
				out[i*3+2] = x[2] + (x[5]-float64(h.v.z))*x[6+i] - float64(h.pos.z)
			}
		},
	}

	settings := linsolve.Settings{MaxIterations: 100000}
	solver := nonlin.NewtonKrylov{
		Maxiter:       10000,
		StepSize:      1,
		Tol:           1e-1,
		InnerSettings: &settings,
	}

	// [{t1: 817432572167, t2: 543347976376, t3: 521380481091, vx: 154, vy: 75, vz: 290, x: 180391926345105, y: 241509806572899, z: 127971479302113}]
	x0 := []float64{
		float64(180391926345105),
		float64(241509806572899),
		float64(127971479302113),
		float64(154),
		float64(75),
		float64(290),
		817432572167.0,
		543347976376.0,
		521380481091.0,
	}
	res := solver.Solve(eqs, x0)
	fmt.Println(res)
}

func part1(hails []Hail) {
	count := 0
	for i := 0; i < len(hails); i++ {
		for j := i + 1; j < len(hails); j++ {
			x, y := hails[i].Intersect(hails[j])
			fmt.Println("---")
			fmt.Println(hails[i], hails[j])
			if math.IsInf(x, 0) || math.IsInf(y, 0) {
				fmt.Println("doesn't intersect")
				continue
			}
			if x < float64(LOW) || x > float64(HIGH) || y < float64(LOW) || y > float64(HIGH) {
				fmt.Println("outside of bounds")
				continue
			}
			fmt.Printf("(%f, %f)\n", x, y)
			count++
		}
	}

	fmt.Println()
	fmt.Println(count)
}
