package player

import (
	"fmt"
	"slices"

	"github.com/jonasah/love-letter/card"
	"github.com/jonasah/love-letter/deck"
)

type Controller interface {
	SelectCardToPlay(card1, card2 card.Card) (card.Card, card.Card)

	// Guard
	GuessCard() card.Card

	// Priest, Baron, King
	SelectPlayer(card card.Card, opponent *Player) *Player

	// Prince
	SelectPlayerToRedraw(self, opponent *Player) *Player

	// Chancellor
	SelectCardToKeep(card1, card2, card3 card.Card) (card.Card, []card.Card)
}

type Player struct {
	Name   string
	Points int

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

func (p *Player) Play(opponent *Player, deck *deck.Deck) {
	var cardToPlay card.Card
	cardToPlay, p.hand = p.controller.SelectCardToPlay(p.hand, deck.Draw())

	p.pile = append(p.pile, cardToPlay)

	p.playCard(cardToPlay, opponent, deck)
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

	// TODO: handle empty deck
	p.hand = deck.Draw()
}

func (p *Player) trade(opponent *Player) {
	p.hand, opponent.hand = opponent.hand, p.hand
}

func (p *Player) playCard(c card.Card, opponent *Player, deck *deck.Deck) {
	fmt.Println(p.Name, "PLAYING", c, "HAND", p.hand)

	cannotPlayIfProtected := []card.Card{card.Guard, card.Priest, card.Baron, card.King}
	if opponent.IsProtected() && slices.Contains(cannotPlayIfProtected, c) {
		fmt.Println("PROTECTED", c)
		return
	}

	actionFunc := cardActions[c]
	if actionFunc != nil {
		actionFunc(p, opponent, deck)
	}
}
