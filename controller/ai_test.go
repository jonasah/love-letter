package controller_test

import (
	"testing"

	"github.com/jonasah/love-letter/card"
	"github.com/jonasah/love-letter/controller"
)

func TestAI_SelectCardToPlay(t *testing.T) {
	tests := []struct {
		name               string
		playCard1          bool
		card1              card.Card
		card2              card.Card
		expectedCardToPlay card.Card
		expectedCardToKeep card.Card
	}{
		// never play princess
		{
			name:               "never play princess (1)",
			playCard1:          true,
			card1:              card.Princess,
			card2:              card.Baron,
			expectedCardToPlay: card.Baron,
			expectedCardToKeep: card.Princess,
		},
		{
			name:               "never play princess (2)",
			playCard1:          false,
			card1:              card.Baron,
			card2:              card.Princess,
			expectedCardToPlay: card.Baron,
			expectedCardToKeep: card.Princess,
		},
		// play countess if required
		{
			name:               "play countess if prince (1)",
			playCard1:          true,
			card1:              card.Prince,
			card2:              card.Countess,
			expectedCardToPlay: card.Countess,
			expectedCardToKeep: card.Prince,
		},
		{
			name:               "play countess if king (1)",
			playCard1:          true,
			card1:              card.King,
			card2:              card.Countess,
			expectedCardToPlay: card.Countess,
			expectedCardToKeep: card.King,
		},
		{
			name:               "play countess if prince (2)",
			playCard1:          false,
			card1:              card.Countess,
			card2:              card.Prince,
			expectedCardToPlay: card.Countess,
			expectedCardToKeep: card.Prince,
		},
		{
			name:               "play countess if king (2)",
			playCard1:          false,
			card1:              card.Countess,
			card2:              card.King,
			expectedCardToPlay: card.Countess,
			expectedCardToKeep: card.King,
		},
		// play random card
		{
			name:               "play random card (1)",
			playCard1:          true,
			card1:              card.Chancellor,
			card2:              card.Baron,
			expectedCardToPlay: card.Chancellor,
			expectedCardToKeep: card.Baron,
		},
		{
			name:               "play random card (2)",
			playCard1:          false,
			card1:              card.Chancellor,
			card2:              card.Baron,
			expectedCardToPlay: card.Baron,
			expectedCardToKeep: card.Chancellor,
		},
	}

	for _, tc := range tests {
		ai := controller.NewAI(&MockRandomizer{b: tc.playCard1})

		cardToPlay, cardToKeep := ai.SelectCardToPlay(tc.card1, tc.card2)
		if cardToPlay != tc.expectedCardToPlay {
			t.Errorf("%s: expected %v to be played, got %v", tc.name, tc.expectedCardToPlay, cardToPlay)
		}
		if cardToKeep != tc.expectedCardToKeep {
			t.Errorf("%s: expected %v to be kept, got %v", tc.name, tc.expectedCardToKeep, cardToKeep)
		}
	}
}

func TestAI_GuessCard(t *testing.T) {
	expectedGuess := card.Baron
	ai := controller.NewAI(&MockRandomizer{c: []card.Card{card.Guard, card.Guard, expectedGuess}})

	guess := ai.GuessCard()
	if guess != expectedGuess {
		t.Errorf("expected guess to be %v, got %v", expectedGuess, guess)
	}
}

type MockRandomizer struct {
	b bool
	c []card.Card
}

func (m *MockRandomizer) Bool() bool {
	return m.b
}

func (m *MockRandomizer) Card() card.Card {
	c := m.c[0]
	m.c = m.c[1:]
	return c
}

func (m *MockRandomizer) Shuffle(cards []card.Card) {
	panic("unimplemented")
}
