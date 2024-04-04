package controller

import (
	"github.com/jonasah/love-letter/card"
	"github.com/jonasah/love-letter/lib"
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
	selected := a.selectUnprotectedOpponent(opponents)
	if selected == nil {
		return opponents[0], -1 // guess doesn't matter
	}

	guess := a.randomizer.Card()
	for guess == card.Guard {
		guess = a.randomizer.Card()
	}

	return selected, guess
}

func (a ai) SelectOpponentForEffect(card card.Card, opponents []*player.Player) *player.Player {
	selected := a.selectUnprotectedOpponent(opponents)
	if selected == nil {
		return opponents[0]
	}

	return selected
}

func (a ai) LookAt(card card.Card) {}

func (a ai) SelectPlayerToRedraw(self *player.Player, opponents []*player.Player) *player.Player {
	selected := a.selectUnprotectedOpponent(opponents)
	if selected == nil {
		// must select self if all opponents are protected
		return self
	}

	if a.randomizer.Bool() {
		return selected
	}

	return self
}

func (a ai) SelectCardToKeep(card1 card.Card, rest ...card.Card) (card.Card, []card.Card) {
	cards := append([]card.Card{card1}, rest...)
	a.randomizer.Shuffle(cards)

	return cards[0], cards[1:]
}

func (a ai) selectUnprotectedOpponent(opponents []*player.Player) *player.Player {
	selectable := lib.Filter(opponents, func(p *player.Player) bool { return !p.IsProtected() })
	if len(selectable) == 0 {
		// all opponents are protected
		return nil
	}

	return selectable[a.randomizer.Intn(len(selectable))]
}

func sort(c1, c2 card.Card) (card.Card, card.Card) {
	if c1 < c2 {
		return c1, c2
	}

	return c2, c1
}
