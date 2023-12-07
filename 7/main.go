package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

const JokersWild = true

type HandType int

const (
	HIGH_CARD HandType = iota
	ONE_PAIR
	TWO_PAIR
	THREE_OF_A_KIND
	FULL_HOUSE
	FOUR_OF_A_KIND
	FIVE_OF_A_KIND
)

func rankValue(c rune) int {
	switch c {
	case '2':
		return 0
	case '3':
		return 1
	case '4':
		return 2
	case '5':
		return 3
	case '6':
		return 4
	case '7':
		return 5
	case '8':
		return 6
	case '9':
		return 7
	case 'T':
		return 8
	case 'J':
		if JokersWild {
			return -1
		}
		return 9
	case 'Q':
		return 10
	case 'K':
		return 11
	case 'A':
		return 12
	}

	panic(c)
}

type Hand struct {
	cards string
	bid   int
}

func (h *Hand) GetType() HandType {
	wilds := 0

	ranks := make([]int, 13)
	for _, c := range h.cards {
		val := rankValue(c)
		if val < 0 {
			wilds++
		} else {
			ranks[val] += 1
		}
	}

	fours := 0
	trips := 0
	pairs := 0
	for _, count := range ranks {
		switch count {
		case 5:
			return FIVE_OF_A_KIND
		case 4:
			fours++
		case 3:
			trips++
		case 2:
			pairs++
		}
	}

	if fours > 0 {
		if wilds > 0 {
			return FIVE_OF_A_KIND
		}
		return FOUR_OF_A_KIND
	}

	if trips > 0 {
		if wilds == 2 {
			return FIVE_OF_A_KIND
		}
		if wilds == 1 {
			return FOUR_OF_A_KIND
		}

		if pairs == 1 {
			return FULL_HOUSE
		}

		return THREE_OF_A_KIND
	}

	if pairs == 2 {
		if wilds == 3 {
			return FIVE_OF_A_KIND
		}
		if wilds == 2 {
			return FOUR_OF_A_KIND
		}
		if wilds == 1 {
			return FULL_HOUSE
		}

		return TWO_PAIR
	}

	if pairs == 1 {
		switch wilds {
		case 3:
			return FIVE_OF_A_KIND
		case 2:
			return FOUR_OF_A_KIND
		case 1:
			return THREE_OF_A_KIND
		}
		return ONE_PAIR
	}

	switch wilds {
	case 4, 5:
		return FIVE_OF_A_KIND
	case 3:
		return FOUR_OF_A_KIND
	case 2:
		return THREE_OF_A_KIND
	case 1:
		return ONE_PAIR
	}

	return HIGH_CARD
}

func main() {
	f, err := os.Open("input")
	if err != nil {
		panic(err)
	}

	scanner := bufio.NewScanner(f)

	var hands []Hand

	for scanner.Scan() {
		line := scanner.Text()
		var hand string
		var bid int
		fmt.Sscanf(line, "%s %d", &hand, &bid)

		hands = append(hands, Hand{hand, bid})
	}

	sort.SliceStable(hands, func(i, j int) bool {
		h1 := hands[i].GetType()
		h2 := hands[j].GetType()

		if h1 < h2 {
			return true
		}

		if h2 < h1 {
			return false
		}

		for n, c := range hands[i].cards {
			if rankValue(c) < rankValue(rune(hands[j].cards[n])) {
				return true
			}

			if rankValue(c) > rankValue(rune(hands[j].cards[n])) {
				return false
			}
		}

		panic("same hand")
	})

	winnings := 0
	for i, h := range hands {
		fmt.Println(h.cards, h.GetType(), i+1, h.bid)
		winnings += (i + 1) * h.bid
	}

	fmt.Println(winnings)
}
