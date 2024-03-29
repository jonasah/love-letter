package deck

import (
	"github.com/jonasah/love-letter/card"
)

type Deck struct {
	cards []card.Card
}

func New() *Deck {
	var cards []card.Card
	for card, count := range cardCount {
		for range count {
			cards = append(cards, card)
		}
	}

	return &Deck{cards}
}

func (d *Deck) Empty() bool {
	return len(d.cards) == 0
}

func (d *Deck) Draw() card.Card {
	card := d.cards[0]
	d.cards = d.cards[1:]
	return card
}

func (d *Deck) Shuffle() {
	// TODO: rand.Shuffle()
}

func (d *Deck) Append(cards []card.Card) {
	d.cards = append(d.cards, cards...)
}

var cardCount = map[card.Card]int{
	card.Spy:      2,
	card.Guard:    6,
	card.Priest:   2,
	card.Baron:    2,
	card.Handmaid: 2,
	card.Prince:   2,
	card.King:     1,
	card.Countess: 1,
	card.Princess: 1,
}
