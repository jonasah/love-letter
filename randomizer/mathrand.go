package randomizer

import (
	"math/rand"

	"github.com/jonasah/love-letter/card"
	"github.com/jonasah/love-letter/controller"
)

type mathRand struct{}

var _ = controller.Randomizer(mathRand{})

func NewMathRand() mathRand {
	return mathRand{}
}

func (m mathRand) Bool() bool {
	return rand.Intn(2) == 0
}

func (m mathRand) Card() card.Card {
	return card.Card(rand.Intn(int(card.Princess) + 1))
}

func (m mathRand) Shuffle(cards []card.Card) {
	rand.Shuffle(len(cards), func(i, j int) {
		cards[i], cards[j] = cards[j], cards[i]
	})
}
