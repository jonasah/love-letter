package card

import (
	"cmp"
)

type Card int

const (
	Spy Card = iota
	Guard
	Priest
	Baron
	Handmaid
	Prince
	Chancellor
	King
	Countess
	Princess
)

var cardToStringMap = map[Card]string{
	Spy:        "Spy",
	Guard:      "Guard",
	Priest:     "Priest",
	Baron:      "Baron",
	Handmaid:   "Handmaid",
	Prince:     "Prince",
	Chancellor: "Chancellor",
	King:       "King",
	Countess:   "Countess",
	Princess:   "Princess",
	-1:         "[OUT]",
}

func (c Card) String() string {
	return cardToStringMap[c]
}

func (c Card) Compare(other Card) int {
	return cmp.Compare(c, other)
}
