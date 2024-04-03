package player

import (
	"log"
	"slices"

	"github.com/jonasah/love-letter/card"
	"github.com/jonasah/love-letter/deck"
)

type Controller interface {
	SelectCardToPlay(card1, card2 card.Card) (card.Card, card.Card)

	// Guard
	GuessCard(opponents []*Player) (*Player, card.Card)

	// Priest, Baron, King
	SelectPlayer(card card.Card, opponents []*Player) *Player

	// Prince
	SelectPlayerToRedraw(self *Player, opponents []*Player) *Player

	// Chancellor
	SelectCardToKeep(card1 card.Card, rest ...card.Card) (card.Card, []card.Card)
}

type Player struct {
	Name   string
	Tokens int

	controller Controller

	hand card.Card
	pile []card.Card
}

func New(name string, controller Controller) *Player {
	return &Player{Name: name, controller: controller}
}

func (p *Player) Deal(hand card.Card) {
	p.hand = hand
	p.pile = nil
}

func (p *Player) Play(opponents []*Player, deck *deck.Deck) {
	var cardToPlay card.Card
	cardToPlay, p.hand = p.controller.SelectCardToPlay(p.hand, deck.Draw())

	p.pile = append(p.pile, cardToPlay)

	log.Printf("Plays %s", cardToPlay)
	effectFunc := cardEffects[cardToPlay]
	if effectFunc != nil {
		effectFunc(p, opponents, deck)
	}
}

func (p *Player) Hand() card.Card {
	return p.hand
}

func (p *Player) IsProtected() bool {
	if len(p.pile) == 0 {
		return false
	}

	return p.pile[len(p.pile)-1] == card.Handmaid
}

func (p *Player) IsOut() bool {
	return p.hand == -1
}

func (p *Player) PlayedSpy() bool {
	return slices.Contains(p.pile, card.Spy)
}

func (p *Player) discardHand() card.Card {
	discardedCard := p.hand
	p.hand = -1
	p.pile = append(p.pile, discardedCard)
	return discardedCard
}

func (p *Player) redrawHand(deck *deck.Deck) {
	discardedCard := p.discardHand()
	if discardedCard == card.Princess {
		return
	}

	if !deck.Empty() {
		p.hand = deck.Draw()
	}
}

func (p *Player) trade(opponent *Player) {
	p.hand, opponent.hand = opponent.hand, p.hand
}
