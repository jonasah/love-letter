package main

import (
	"fmt"
	"slices"

	"github.com/jonasah/love-letter/controller"
	"github.com/jonasah/love-letter/deck"
	"github.com/jonasah/love-letter/player"
	"github.com/jonasah/love-letter/randomizer"
)

func main() {
	// TODO: more players
	var players []*player.Player
	for i := range 2 {
		p := player.New(fmt.Sprintf("Player%d", i+1), controller.NewAI(randomizer.NewMathRand()))
		players = append(players, p)
	}

	winningPoints := 5

	for {
		playRound(players)

		fmt.Println(players[0], players[1])

		winnerIdx := slices.IndexFunc(players, func(p *player.Player) bool { return p.Points >= winningPoints })
		if winnerIdx != -1 {
			fmt.Println("WINNER", players[winnerIdx])
			break
		}
	}
}

func playRound(players []*player.Player) {
	fmt.Println("==== NEW ROUND ====")

	deck := deck.New()
	deck.Shuffle()
	removedCard := deck.Draw()

	fmt.Println(*deck)

	for _, p := range players {
		p.Deal(deck.Draw())
	}

	// var out []player.Player

	for !deck.Empty() && len(players) > 1 {
		players[0].Play(players[1], deck)

		players = append(players[1:], players[0])
		players = slices.DeleteFunc(players, func(p *player.Player) bool { return p.IsOut() })
	}

	fmt.Println("REMOVED CARD", removedCard)

	if len(players) == 1 {
		fmt.Println("LAST MAN STANDING", *players[0])
		players[0].Points++
		if players[0].PlayedSpy() {
			fmt.Println("EXTRA POINT SPY")
			players[0].Points++
		}
		return
	}

	p1Card := players[0].Hand()
	p2Card := players[1].Hand()
	if p1Card > p2Card {
		fmt.Println("WIN", players[0].Name, p1Card, ">", p2Card)
		players[0].Points++
	} else if p1Card < p2Card {
		fmt.Println("WIN", players[1].Name, p2Card, ">", p1Card)
		players[1].Points++
	} else {
		fmt.Println("TIE", p1Card, "==", p2Card)
		players[0].Points++
		players[1].Points++
	}

	p1Spy := players[0].PlayedSpy()
	p2Spy := players[1].PlayedSpy()
	if p1Spy == p2Spy {
		fmt.Println("NO POINT FOR SPY", p1Spy, p2Spy)
	} else if p1Spy {
		fmt.Println("SPY", players[0].Name)
		players[0].Points++
	} else {
		fmt.Println("SPY", players[1].Name)
		players[1].Points++
	}
}
