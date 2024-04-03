package card

import (
	"cmp"
	"fmt"
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
}

func (c Card) String() string {
	if c == -1 {
		return "[OUT]"
	}

	return fmt.Sprintf("%s(%d)", cardToStringMap[c], c)
}

func (c Card) Compare(other Card) int {
	return cmp.Compare(c, other)
}
