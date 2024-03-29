package player

import (
	"fmt"

	"github.com/jonasah/love-letter/card"
	"github.com/jonasah/love-letter/deck"
)

type cardActionFn = func(self, opponent *Player, deck *deck.Deck)

var cardActions = map[card.Card]cardActionFn{
	card.Guard:      playGuard,
	card.Priest:     playPriest,
	card.Baron:      playBaron,
	card.Prince:     playPrince,
	card.Chancellor: playChancellor,
	card.King:       playKing,
	card.Princess:   playPrincess,
}

func playGuard(self, opponent *Player, deck *deck.Deck) {
	// TODO: select player to guess on

	guess := self.controller.GuessCard()
	if opponent.Hand() == guess {
		fmt.Println("GUARD", "CORRECT", opponent.Hand())
		opponent.discardHand()
	} else {
		fmt.Println("GUARD", "INCORRECT", guess, "ACTUAL", opponent.Hand())
	}
}

func playPriest(self, opponent *Player, deck *deck.Deck) {
	// TODO: select player to look at

	fmt.Println("PRIEST", opponent.Hand())
}

func playBaron(self, opponent *Player, deck *deck.Deck) {
	// TODO: select player to compare with

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

func playPrince(self, opponent *Player, deck *deck.Deck) {
	playerToDiscard := self.controller.SelectPlayerToRedraw(self, opponent)
	if playerToDiscard == self || !playerToDiscard.IsProtected() {
		before := playerToDiscard.hand
		playerToDiscard.redrawHand(deck)
		fmt.Println("PRINCE", before, "->", playerToDiscard.hand)
	} else {
		fmt.Println("PRINCE", "PROTECTED")
	}
}

func playChancellor(self, opponent *Player, deck *deck.Deck) {
	var drawnCards []card.Card
	for range 2 {
		if !deck.Empty() {
			drawnCards = append(drawnCards, deck.Draw())
		}
	}

	var returnCards []card.Card
	self.hand, returnCards = self.controller.SelectCardToKeep(self.hand, drawnCards...)
	deck.Append(returnCards)
}

func playKing(self, opponent *Player, deck *deck.Deck) {
	// TODO: select player to trade with
	fmt.Println("KING", self.hand, "<->", opponent.hand)
	self.trade(opponent)
}

func playPrincess(self, opponent *Player, deck *deck.Deck) {
	self.discardHand()
}
