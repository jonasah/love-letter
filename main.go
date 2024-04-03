package main

import (
	"fmt"
	"log"
	"slices"

	"github.com/jonasah/love-letter/controller"
	"github.com/jonasah/love-letter/deck"
	"github.com/jonasah/love-letter/player"
	"github.com/jonasah/love-letter/randomizer"
)

func main() {
	log.SetFlags(0)

	rnd := randomizer.NewMathRand()

	numPlayers := rnd.IntIncl(2, 6)
	var players []*player.Player
	for i := range numPlayers {
		p := player.New(fmt.Sprintf("Player%d", i+1), controller.NewAI(rnd))
		players = append(players, p)
	}

	tokensToWin := map[int]int{2: 6, 3: 5, 4: 4, 5: 3, 6: 3}[numPlayers]

	log.Printf("Start game with %d players", len(players))

	for {
		playRound(slices.Clone(players), rnd)

		winnerIdx := slices.IndexFunc(players, func(p *player.Player) bool { return p.Tokens >= tokensToWin })
		if winnerIdx != -1 {
			log.Print("=== END OF GAME ===")
			log.Printf("%s wins the game with %d tokens", players[winnerIdx].Name, players[winnerIdx].Tokens)
			for _, p := range players {
				log.Printf("- %s: %d", p.Name, p.Tokens)
			}
			break
		}
	}
}

func playRound(players []*player.Player, shuffler deck.Shuffler) {
	log.Print("=== START ROUND ===")

	deck := deck.New(shuffler)
	removedCard := deck.Draw()

	for _, p := range players {
		p.Deal(deck.Draw())
	}

	// TODO: winner of last round begins

	for !deck.Empty() && len(players) > 1 {
		log.Printf("-- %s's turn --", players[0].Name)

		players[0].Play(players[1:], deck)

		players = append(players[1:], players[0])
		players = slices.DeleteFunc(players, func(p *player.Player) bool { return p.IsOut() })
	}

	log.Print("--- End of round ---")

	log.Printf("Removed card: %s", removedCard)

	if len(players) == 1 {
		log.Printf("%s wins the round as the last player standing", players[0].Name)
		players[0].Tokens++

		if players[0].PlayedSpy() {
			log.Print("Extra token awarded for Spy")
			players[0].Tokens++
		}

		return
	}

	slices.SortFunc(players, func(a, b *player.Player) int { return -a.Hand().Compare(b.Hand()) })

	log.Printf("Deck is empty, %d players left:", len(players))
	for _, p := range players {
		log.Printf("- %s: %s", p.Name, p.Hand())
	}

	if players[0].Hand() > players[1].Hand() {
		log.Printf("%s wins the round with the highest card", players[0].Name)
		players[0].Tokens++
	} else if players[0].Hand() == players[1].Hand() {
		log.Printf("%s and %s tie the round with the highest card", players[0].Name, players[1].Name)
		players[0].Tokens++
		players[1].Tokens++
	}

	var playersWithSpy []*player.Player
	for _, p := range players {
		if p.PlayedSpy() {
			playersWithSpy = append(playersWithSpy, p)
		}
	}

	if len(playersWithSpy) == 1 {
		log.Printf("%s is awarded an extra token for Spy", playersWithSpy[0].Name)
		playersWithSpy[0].Tokens++
	} else {
		log.Printf("No extra token for Spy. %d players played a Spy", len(playersWithSpy))
	}
}
