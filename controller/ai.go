package controller

import (
	"github.com/jonasah/love-letter/card"
	"github.com/jonasah/love-letter/player"
)

type Randomizer interface {
	Bool() bool
	Card() card.Card
}

type ai struct {
	randomizer Randomizer
}

var _ = player.Controller(ai{})

func NewAI(randomizer Randomizer) ai {
	return ai{randomizer}
}

func (a ai) SelectCardToPlay(card1, card2 card.Card) (card.Card, card.Card) {
	low, high := sort(card1, card2)

	// never play princess
	if high == card.Princess {
		return low, high
	}

	// play countess if required
	if high == card.Countess && (low == card.Prince || low == card.King) {
		return high, low
	}

	if a.randomizer.Bool() {
		return card1, card2
	}

	return card2, card1
}

func (a ai) GuessCard() card.Card {
	guess := a.randomizer.Card()
	for guess == card.Guard {
		guess = a.randomizer.Card()
	}

	return guess
}

func (a ai) SelectPlayer(card card.Card, opponent *player.Player) *player.Player {
	return opponent
}

func (a ai) SelectPlayerToRedraw(self, opponent *player.Player) *player.Player {
	if opponent.IsProtected() {
		return self
	}

	if a.randomizer.Bool() {
		return opponent
	}

	return self
}

func (a ai) SelectCardToKeep(card1, card2, card3 card.Card) (card.Card, []card.Card) {
	// TODO: random
	// idxToKeep := 0 //rand.Intn(len(cards))
	// return cards[idxToKeep], cards[1:]

	return card1, []card.Card{card2, card3}
}

func sort(c1, c2 card.Card) (card.Card, card.Card) {
	if c1 < c2 {
		return c1, c2
	}

	return c2, c1
}
