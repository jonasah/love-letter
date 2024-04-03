package player

import (
	"log"

	"github.com/jonasah/love-letter/card"
	"github.com/jonasah/love-letter/deck"
)

type cardEffectFn = func(self *Player, opponents []*Player, deck *deck.Deck)

var cardEffects = map[card.Card]cardEffectFn{
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
		log.Printf("> No effect, %s is protected", opponent.Name)
		return
	}

	if opponent.Hand() == guess {
		log.Printf("> Guess %s on %s: CORRECT", guess, opponent.Name)
		opponent.discardHand()
		return
	}

	log.Printf("> Guess %s on %s: INCORRECT", guess, opponent.Name)
}

func playPriest(self *Player, opponents []*Player, deck *deck.Deck) {
	opponent := self.controller.SelectPlayer(card.Priest, opponents)

	if opponent.IsProtected() {
		log.Printf("> No effect, %s is protected", opponent.Name)
		return
	}

	log.Printf("> Looks at %s's hand", opponent.Name)
	self.controller.LookAt(opponent.Hand())
}

func playBaron(self *Player, opponents []*Player, deck *deck.Deck) {
	opponent := self.controller.SelectPlayer(card.Baron, opponents)

	if opponent.IsProtected() {
		log.Printf("> No effect, %s is protected", opponent.Name)
		return
	}

	switch self.Hand().Compare(opponent.Hand()) {
	case 1:
		log.Printf("> Beats %s who discards %s", opponent.Name, opponent.Hand())
		opponent.discardHand()
	case -1:
		log.Printf("> Loses to %s, discards %s", opponent.Name, self.Hand())
		self.discardHand()
	default:
		log.Printf("> Ties with %s", opponent.Name)
	}
}

func playPrince(self *Player, opponents []*Player, deck *deck.Deck) {
	playerToDiscard := self.controller.SelectPlayerToRedraw(self, opponents)

	if playerToDiscard != self && playerToDiscard.IsProtected() {
		log.Printf("> No effect, %s is protected", playerToDiscard.Name)
		return
	}

	discardedCard := playerToDiscard.discardHand()
	log.Printf("> %s discards %s", playerToDiscard.Name, discardedCard)
	if discardedCard == card.Princess {
		log.Printf("> %s is OUT", playerToDiscard.Name)
		return
	}

	// TODO: pick up removed card if deck is empty
	if !deck.Empty() {
		playerToDiscard.hand = deck.Draw()
	}
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
}

func playKing(self *Player, opponents []*Player, deck *deck.Deck) {
	opponent := self.controller.SelectPlayer(card.King, opponents)

	if opponent.IsProtected() {
		log.Printf("> No effect, %s is protected", opponent.Name)
		return
	}

	log.Printf("> Trades with %s", opponent.Name)
	self.hand, opponent.hand = opponent.hand, self.hand
}

func playPrincess(self *Player, opponents []*Player, deck *deck.Deck) {
	self.discardHand()
}
