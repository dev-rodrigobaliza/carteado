package deck

import (
	"github.com/dev-rodrigobaliza/carteado/errors"
	"github.com/dev-rodrigobaliza/carteado/internal/core/card"
	"github.com/dev-rodrigobaliza/carteado/vars"
)

type Deck struct {
	Cards []*card.Card
}

func New(jokers int) (*Deck, error) {
	cards := make([]*card.Card, 0)

	for i, face := range vars.CardFaces {
		for j, suit := range vars.CardSuits {
			c, err := card.New(face, suit, (i + 1), (j + 1), false)
			if err != nil {
				return nil, err
			}
			cards = append(cards, c)
		}
	}

	for i := 0; i < jokers; i++ {
		c, err := card.New("j", "j", 0, 0, false)
		if err != nil {
			return nil, err
		}
		cards = append(cards, c)
	}

	deck := &Deck{
		Cards: cards,
	}

	return deck, nil
}

func NewCustom(cards []string, values []int) (*Deck, error) {
	if len(cards) == 0 {
		return nil, errors.ErrInvalidCard
	}
	if len(values) == 0 {
		return nil, errors.ErrInvalidCardValue
	}
	if len(cards) != len(values) {
		return nil, errors.ErrInvalidCardDeck
	}

	cas := make([]*card.Card, 0)

	// 1d 4g 6g k!
	for i, c := range cards {
		ca, err := card.New(c[:1], c[1:], values[i], 1, false)
		if err != nil {
			return nil, err
		}

		cas = append(cas, ca)
	}

	deck := &Deck{
		Cards: cas,
	}

	return deck, nil
}
