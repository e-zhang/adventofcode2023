package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"
)

func main() {
	f, err := os.Open("input")
	if err != nil {
		panic(err)
	}

	connections := map[string][]string{}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()

		s := strings.Split(line, ":")
		src := s[0]
		targets := s[1]

		for _, t := range strings.Split(targets, " ") {
			if t == "" {
				continue
			}
			connections[src] = append(connections[src], t)
			connections[t] = append(connections[t], src)
		}
	}

	// var nodes []string
	// for k, v := range connections {
	// 	nodes = append(nodes, k)
	// 	fmt.Println(k, len(v), v)
	// }

	// group1, group2 := map[string]struct{}{}, map[string]struct{}{}
	// ok := assign(nodes, connections, group1, group2, 3)
	// if !ok {
	// 	fmt.Println(group1, group2)
	// 	panic(ok)
	// }

	// fmt.Println(group1)
	// fmt.Println(group2)
	// fmt.Println(len(group1) * len(group2))
	cuts := 0
	var g1, g2 []string
	for cuts != 3 {
		cuts, g1, g2 = minCut(connections)
		fmt.Println(cuts, g1, g2)
	}
	fmt.Println(len(g1), len(g2), len(g1)*len(g2))
}

func randomNodes(g map[string][]string) (string, string) {
	i := rand.Intn(len(g))
	for k, v := range g {
		if i == 0 {
			j := rand.Intn(len(v))
			return k, v[j]
		}
		i--
	}

	panic(i)
}

func minCut(connections map[string][]string) (int, []string, []string) {
	cg := map[string][]string{}
	for k, v := range connections {
		newv := make([]string, len(v))
		copy(newv, v)
		cg[k] = newv
	}

	for len(cg) > 2 {
		u, v := randomNodes(cg)
		// fmt.Println(k, v)
		merged := u + "-" + v
		var combined []string
		for _, i := range append(cg[u], cg[v]...) {
			if i != u && i != v {
				combined = append(combined, i)

				var s []string
				for _, x := range cg[i] {
					if x == u || x == v {
						x = merged
					}
					s = append(s, x)
				}
				cg[i] = s
				// fmt.Println(merged, i, s, combined)
			}
		}

		delete(cg, u)
		delete(cg, v)
		cg[merged] = combined
	}

	cuts := 0
	var g1, g2 []string
	for k, v := range cg {
		cuts = len(v)
		if g1 == nil {
			g1 = strings.Split(k, "-")
		} else {
			g2 = strings.Split(k, "-")
		}
	}
	return cuts, g1, g2
}

func assign(nodes []string, connections map[string][]string, dst, other map[string]struct{}, cuts int) bool {
	if len(nodes) == 0 {
		return len(dst) > 0 && len(other) > 0
	}

	node := nodes[0]
	remaining := nodes[1:]

	if c := countCuts(node, connections, other); c <= cuts {
		dst[node] = struct{}{}
		if assign(remaining, connections, dst, other, cuts-c) {
			return true
		}
		delete(dst, node)
	}

	if c := countCuts(node, connections, dst); c <= cuts {
		other[node] = struct{}{}
		if assign(remaining, connections, dst, other, cuts-c) {
			return true
		}
		delete(other, node)
	}

	return false
}

func countCuts(node string, connections map[string][]string, group map[string]struct{}) int {
	conns := connections[node]
	cuts := 0
	for _, c := range conns {
		if _, ok := group[c]; ok {
			cuts++
		}
	}

	return cuts
}
