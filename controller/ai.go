package controller

import (
	"github.com/jonasah/love-letter/card"
	"github.com/jonasah/love-letter/player"
)

type Randomizer interface {
	Bool() bool
	Intn(n int) int

	Card() card.Card
	Shuffle(cards []card.Card)
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

func (a ai) GuessCard(opponents []*player.Player) (*player.Player, card.Card) {
	// TODO: do not guess on protected opponent
	opponent := opponents[a.randomizer.Intn(len(opponents))]

	guess := a.randomizer.Card()
	for guess == card.Guard {
		guess = a.randomizer.Card()
	}

	return opponent, guess
}

func (a ai) SelectPlayer(card card.Card, opponents []*player.Player) *player.Player {
	// TODO: do not guess on protected opponent
	opponent := opponents[a.randomizer.Intn(len(opponents))]
	return opponent
}

func (a ai) LookAt(card card.Card) {}

func (a ai) SelectPlayerToRedraw(self *player.Player, opponents []*player.Player) *player.Player {
	// TODO: do not guess on protected opponent
	opponent := opponents[a.randomizer.Intn(len(opponents))]

	if opponent.IsProtected() {
		return self
	}

	if a.randomizer.Bool() {
		return opponent
	}

	return self
}

func (a ai) SelectCardToKeep(card1 card.Card, rest ...card.Card) (card.Card, []card.Card) {
	cards := append([]card.Card{card1}, rest...)
	a.randomizer.Shuffle(cards)

	return cards[0], cards[1:]
}

func sort(c1, c2 card.Card) (card.Card, card.Card) {
	if c1 < c2 {
		return c1, c2
	}

	return c2, c1
}
