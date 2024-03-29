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

	switch c {
	case card.Spy:
	case card.Guard:
		// TODO: select player to guess on

		guess := p.controller.GuessCard()
		if opponent.Hand() == guess {
			fmt.Println("GUARD", "CORRECT", opponent.Hand())
			opponent.discardHand()
		} else {
			fmt.Println("GUARD", "INCORRECT", guess, "ACTUAL", opponent.Hand())
		}
	case card.Priest:
		fmt.Println("PRIEST", opponent.Hand())
	case card.Baron:
		// TODO: select player to compare with

		switch c.Compare(opponent.Hand()) {
		case 1:
			fmt.Println("BARON", c, "BEATS", opponent.Hand())
			opponent.discardHand()
		case -1:
			fmt.Println("BARON", c, "LOSES", opponent.Hand())
			p.discardHand()
		default:
			fmt.Println("BARON", c, "EQUALS", opponent.Hand())
		}
	case card.Handmaid:
	case card.Prince:
		playerToDiscard := p.controller.SelectPlayerToRedraw(p, opponent)
		if playerToDiscard == p || !playerToDiscard.IsProtected() {
			playerToDiscard.redrawHand(deck)
		}
	case card.Chancellor:
		var returnCards []card.Card
		p.hand, returnCards = p.controller.SelectCardToKeep(p.hand, deck.Draw(), deck.Draw())
		deck.Append(returnCards)
	case card.King:
		// TODO: select player to trade with
		p.trade(opponent)
	case card.Countess:
	case card.Princess:
		p.discardHand()
	}
}
