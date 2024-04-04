package deck

import (
	"github.com/jonasah/love-letter/card"
)

type Shuffler interface {
	Shuffle(cards []card.Card)
}

type Deck struct {
	cards       []card.Card
	removedCard card.Card
}

func New(shuffler Shuffler) *Deck {
	var cards []card.Card
	for card, count := range cardCount {
		for range count {
			cards = append(cards, card)
		}
	}

	shuffler.Shuffle(cards)

	return &Deck{cards: cards, removedCard: -1}
}

func (d *Deck) Empty() bool {
	return len(d.cards) == 0
}

func (d *Deck) Draw() card.Card {
	card := d.cards[0]
	d.cards = d.cards[1:]
	return card
}

func (d *Deck) Append(cards []card.Card) {
	d.cards = append(d.cards, cards...)
}

func (d *Deck) RemoveCard() {
	d.removedCard = d.Draw()
}

func (d *Deck) RemovedCard() card.Card {
	return d.removedCard
}

var cardCount = map[card.Card]int{
	card.Spy:        2,
	card.Guard:      6,
	card.Priest:     2,
	card.Baron:      2,
	card.Handmaid:   2,
	card.Prince:     2,
	card.Chancellor: 2,
	card.King:       1,
	card.Countess:   1,
	card.Princess:   1,
}
