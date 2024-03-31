package player

import (
	"fmt"

	"github.com/jonasah/love-letter/card"
	"github.com/jonasah/love-letter/deck"
)

type cardActionFn = func(self *Player, opponents []*Player, deck *deck.Deck)

var cardActions = map[card.Card]cardActionFn{
	card.Guard:      playGuard,
	card.Priest:     playPriest,
	card.Baron:      playBaron,
	card.Prince:     playPrince,
	card.Chancellor: playChancellor,
	card.King:       playKing,
	card.Princess:   playPrincess,
}

func playGuard(self *Player, opponents []*Player, deck *deck.Deck) {
	opponent, guess := self.controller.GuessCard(opponents)

	if opponent.IsProtected() {
		fmt.Println("GUARD", "PROTECTED")
		return
	}

	if opponent.Hand() == guess {
		fmt.Println("GUARD", "CORRECT", opponent.Hand())
		opponent.discardHand()
		return
	}

	fmt.Println("GUARD", "INCORRECT", guess, "ACTUAL", opponent.Hand())
}

func playPriest(self *Player, opponents []*Player, deck *deck.Deck) {
	opponent := self.controller.SelectPlayer(card.Priest, opponents)

	if opponent.IsProtected() {
		fmt.Println("PRIEST", "PROTECTED")
		return
	}

	fmt.Println("PRIEST", opponent.Hand())
}

func playBaron(self *Player, opponents []*Player, deck *deck.Deck) {
	opponent := self.controller.SelectPlayer(card.Baron, opponents)

	if opponent.IsProtected() {
		fmt.Println("BARON", "PROTECTED")
		return
	}

	switch self.Hand().Compare(opponent.Hand()) {
	case 1:
		fmt.Println("BARON", self.Hand(), ">", opponent.Hand())
		opponent.discardHand()
	case -1:
		fmt.Println("BARON", self.Hand(), "<", opponent.Hand())
		self.discardHand()
	default:
		fmt.Println("BARON", self.Hand(), "==", opponent.Hand())
	}
}

func playPrince(self *Player, opponents []*Player, deck *deck.Deck) {
	playerToDiscard := self.controller.SelectPlayerToRedraw(self, opponents)

	if playerToDiscard != self && playerToDiscard.IsProtected() {
		fmt.Println("PRINCE", "PROTECTED")
		return
	}

	before := playerToDiscard.Hand()
	playerToDiscard.redrawHand(deck)
	fmt.Println("PRINCE", before, "->", playerToDiscard.Hand())
}

func playChancellor(self *Player, opponents []*Player, deck *deck.Deck) {
	var drawnCards []card.Card
	for range 2 {
		if !deck.Empty() {
			drawnCards = append(drawnCards, deck.Draw())
		}
	}

	var returnCards []card.Card
	self.hand, returnCards = self.controller.SelectCardToKeep(self.Hand(), drawnCards...)
	deck.Append(returnCards)

	fmt.Println("CHANCELLOR", self.Hand(), returnCards)
}

func playKing(self *Player, opponents []*Player, deck *deck.Deck) {
	opponent := self.controller.SelectPlayer(card.King, opponents)

	if opponent.IsProtected() {
		fmt.Println("KING", "PROTECTED")
		return
	}

	fmt.Println("KING", self.Hand(), "<->", opponent.Hand())
	self.trade(opponent)
}

func playPrincess(self *Player, opponents []*Player, deck *deck.Deck) {
	self.discardHand()
}
