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
	rnd := randomizer.NewMathRand()

	numPlayers := rnd.IntIncl(2, 6)
	var players []*player.Player
	for i := range numPlayers {
		p := player.New(fmt.Sprintf("Player%d", i+1), controller.NewAI(rnd))
		players = append(players, p)
	}

	tokensToWin := map[int]int{2: 6, 3: 5, 4: 4, 5: 3, 6: 3}[numPlayers]

	for {
		playRound(players, rnd)

		for _, p := range players {
			fmt.Println(p)
		}

		winnerIdx := slices.IndexFunc(players, func(p *player.Player) bool { return p.Tokens >= tokensToWin })
		if winnerIdx != -1 {
			fmt.Println("WINNER", players[winnerIdx])
			break
		}
	}
}

func playRound(players []*player.Player, shuffler deck.Shuffler) {
	fmt.Println("==== NEW ROUND ====")

	deck := deck.New(shuffler)
	removedCard := deck.Draw()

	fmt.Println(*deck)

	for _, p := range players {
		p.Deal(deck.Draw())
	}

	// TODO: winner of last round begins

	for !deck.Empty() && len(players) > 1 {
		players[0].Play(players[1:], deck)

		players = append(players[1:], players[0])
		players = slices.DeleteFunc(players, func(p *player.Player) bool { return p.IsOut() })
	}

	fmt.Println("REMOVED CARD", removedCard)

	if len(players) == 1 {
		fmt.Println("LAST MAN STANDING", *players[0])
		players[0].Tokens++
		if players[0].PlayedSpy() {
			fmt.Println("EXTRA POINT SPY")
			players[0].Tokens++
		}
		return
	}

	slices.SortFunc(players, func(a, b *player.Player) int { return -a.Hand().Compare(b.Hand()) })

	if players[0].Hand() > players[1].Hand() {
		fmt.Println("WIN", players[0].Name, players[0].Hand(), ">", players[1].Hand(), players[1].Name)
		players[0].Tokens++
	} else if players[0].Hand() == players[1].Hand() {
		fmt.Println("TIE", players[0].Name, players[0].Hand(), "==", players[1].Hand(), players[1].Name)
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
		fmt.Println("SPY", playersWithSpy[0].Name)
		playersWithSpy[0].Tokens++
	} else {
		fmt.Println("NO POINT FOR SPY", len(playersWithSpy))
	}
}
