package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
)

type Pulse struct {
	source      string
	destination string
	pulse       bool
}

func (p *Pulse) String() string {
	pulse := "low"
	if p.pulse {
		pulse = "high"
	}

	return fmt.Sprintf("%s -%s-> %s", p.source, pulse, p.destination)
}

type Module interface {
	Handle(input string, pulse bool) []Pulse
	Name() string
	Destinations() []string
	Inputs([]string)
	String() string
	Reset()
}

func makePulses(name string, destinations []string, pulse bool) []Pulse {
	pulses := make([]Pulse, len(destinations))
	for i := range pulses {
		pulses[i] = Pulse{
			name,
			destinations[i],
			pulse,
		}
	}
	return pulses
}

type FlipFlop struct {
	name         string
	destinations []string

	state bool
}

func (f *FlipFlop) Name() string {
	return f.name
}

func (f *FlipFlop) Destinations() []string {
	return f.destinations
}

func (f *FlipFlop) Inputs([]string) {}

func (f *FlipFlop) Reset() { f.state = false }

func (f *FlipFlop) Handle(_ string, pulse bool) []Pulse {
	if pulse {
		return nil
	}

	f.state = !f.state

	return makePulses(f.Name(), f.Destinations(), f.state)
}

func (f *FlipFlop) String() string {
	return fmt.Sprintf("%s: %t", f.Name(), f.state)
}

type Conjunction struct {
	name         string
	destinations []string

	memory map[string]bool
}

func (c *Conjunction) Name() string {
	return c.name
}

func (c *Conjunction) Reset() {
	for k := range c.memory {
		c.memory[k] = false
	}
}

func (c *Conjunction) Destinations() []string {
	return c.destinations
}

func (c *Conjunction) Inputs(in []string) {
	for _, n := range in {
		c.memory[n] = false
	}
}

func (c *Conjunction) Handle(input string, pulse bool) []Pulse {
	c.memory[input] = pulse

	send := true
	for _, v := range c.memory {
		send = send && v
	}

	return makePulses(c.Name(), c.Destinations(), !send)
}

func (c *Conjunction) String() string {
	in := []string{}
	for k := range c.memory {
		in = append(in, k)
	}
	sort.Strings(in)

	mem := ""
	for i, n := range in {
		if i > 0 {
			mem += ","
		}
		mem += fmt.Sprintf("%s:%t", n, c.memory[n])
	}

	return fmt.Sprintf("%s: [%s]", c.Name(), mem)
}

type Broadcaster struct {
	destinations []string
}

func (b *Broadcaster) Name() string {
	return "broadcaster"
}

func (b *Broadcaster) Handle(_ string, pulse bool) []Pulse {
	return makePulses(b.Name(), b.Destinations(), pulse)
}

func (b *Broadcaster) Destinations() []string {
	return b.destinations
}

func (b *Broadcaster) Inputs([]string) {}
func (b *Broadcaster) Reset()          {}

func (b *Broadcaster) String() string { return b.Name() }

func parseModule(line string) Module {
	module := strings.Split(line, " -> ")
	dest := strings.Split(module[1], ", ")

	if module[0] == "broadcaster" {
		return &Broadcaster{dest}
	}

	if module[0][0] == '%' {
		return &FlipFlop{module[0][1:], dest, false}
	}

	if module[0][0] == '&' {
		return &Conjunction{module[0][1:], dest, map[string]bool{}}
	}

	panic(line)
}

func main() {
	f, err := os.Open("input")
	if err != nil {
		panic(err)
	}

	modules := map[string]Module{}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		m := parseModule(line)
		modules[m.Name()] = m
	}

	inputs := map[string][]string{}
	for _, m := range modules {
		for _, d := range m.Destinations() {
			in, ok := inputs[d]
			if !ok {
				in = []string{}
			}
			inputs[d] = append(in, m.Name())
		}
	}

	for k, v := range inputs {
		m, ok := modules[k]
		if !ok {
			continue
		}
		m.Inputs(v)
	}

	part1(modules)
	for _, m := range modules {
		m.Reset()
	}
	part2(modules)
}

func findTargetSource(modules map[string]Module, target string) string {
	for _, m := range modules {
		for _, d := range m.Destinations() {
			if d == target {
				return m.Name()
			}
		}
	}

	panic(target)
}

func part2(modules map[string]Module) {
	target := findTargetSource(modules, "rx")
	count := 0
	for _, m := range modules {
		for _, d := range m.Destinations() {
			if d == target {
				count++
				break
			}
		}
	}

	fmt.Println(target, count)

	cycles := map[string]int{}
	i := 0
	for len(cycles) < count {
		seen := button2(modules, target)
		i++
		if len(seen) > 0 {
			// fmt.Println(seen, i)
			cycles[seen] = i
		}
	}

	v := 0
	for _, s := range cycles {
		if v == 0 {
			v = s
			continue
		}

		v = lcm(v, s)
	}

	fmt.Println(v)
}

func gcd(a, b int) int {
	for b != 0 {
		t := b
		b = a % b
		a = t
	}
	return a
}

func lcm(a, b int) int {
	return a * b / gcd(a, b)
}

func button2(modules map[string]Module, target string) string {
	pulses := []Pulse{
		Pulse{"button", "broadcaster", false},
	}

	sender := ""
	for len(pulses) > 0 {
		p := pulses[0]
		pulses = pulses[1:]

		if p.destination == target && p.pulse {
			sender = p.source
		}

		m, ok := modules[p.destination]
		if !ok {
			continue
		}

		out := m.Handle(p.source, p.pulse)
		pulses = append(pulses, out...)
	}

	return sender
}

func part1(modules map[string]Module) {
	totalL, totalH := 0, 0
	for i := 0; i < 1000; i++ {
		// fmt.Println("====")
		low, high := button1(modules)
		totalL += low
		totalH += high
	}

	fmt.Println(totalL * totalH)
}

func button1(modules map[string]Module) (int, int) {
	low, high := 0, 0

	pulses := []Pulse{
		Pulse{"button", "broadcaster", false},
	}

	for len(pulses) > 0 {
		p := pulses[0]
		pulses = pulses[1:]

		if p.pulse {
			high++
		} else {
			low++
		}

		m, ok := modules[p.destination]
		if !ok {
			continue
		}

		out := m.Handle(p.source, p.pulse)
		pulses = append(pulses, out...)

		// for _, o := range out {
		// 	fmt.Println(o.String())
		// }
	}

	return low, high
}
