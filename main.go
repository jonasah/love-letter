package main

import (
	"fmt"
	"log"
	"slices"
	"strings"

	"github.com/jonasah/love-letter/controller"
	"github.com/jonasah/love-letter/deck"
	"github.com/jonasah/love-letter/lib"
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

	log.Printf("Start game with %d players. %d tokens needed to win.", len(players), tokensToWin)

	for {
		winner := playRound(slices.Clone(players), rnd)

		gameOver := slices.ContainsFunc(players, func(p *player.Player) bool { return p.Tokens >= tokensToWin })
		if gameOver {
			break
		}

		// winner starts next round
		for players[0] != winner {
			players = append(players[1:], players[0])
		}
	}

	log.Print("=== END OF GAME ===")
	slices.SortFunc(players, func(a, b *player.Player) int { return b.Tokens - a.Tokens })

	winners := lib.Filter(players, func(p *player.Player) bool { return p.Tokens == players[0].Tokens })
	names := lib.Transform(winners, func(p *player.Player) string { return p.Name })
	log.Printf("%s wins the game!", strings.Join(names, " and "))

	log.Print("Results:")
	for _, p := range players {
		log.Printf("- %s: %d", p.Name, p.Tokens)
	}
}

func playRound(players []*player.Player, shuffler deck.Shuffler) *player.Player {
	log.Print("=== START ROUND ===")

	deck := deck.New(shuffler)
	deck.RemoveCard()

	for _, p := range players {
		p.Deal(deck.Draw())
	}

	for !deck.Empty() && len(players) > 1 {
		log.Printf("-- %s's turn --", players[0].Name)

		players[0].Play(players[1:], deck)

		players = append(players[1:], players[0])
		players = slices.DeleteFunc(players, func(p *player.Player) bool { return p.IsOut() })
	}

	log.Print("--- End of round ---")

	log.Printf("Removed card: %s", deck.RemovedCard())

	if len(players) == 1 {
		log.Printf("%s wins the round as the last player standing", players[0].Name)
		players[0].Tokens++

		if players[0].PlayedSpy() {
			log.Print("Extra token awarded for Spy")
			players[0].Tokens++
		}

		return players[0]
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

	return players[0]
}
