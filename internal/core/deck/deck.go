package deck

import (
	"github.com/dev-rodrigobaliza/carteado/errors"
	"github.com/dev-rodrigobaliza/carteado/internal/core/card"
	"github.com/dev-rodrigobaliza/carteado/pkg/safemap"
	"github.com/dev-rodrigobaliza/carteado/vars"
)

type Deck struct {
	cards *safemap.SafeMap[string, *card.Card]
}

func New(jokers int) (*Deck, error) {
	cards := safemap.New[string, *card.Card]()

	for i, face := range vars.CardFaces {
		for j, suit := range vars.CardSuits {
			c, err := card.New(face, suit, (i + 1), (j + 1), false)
			if err != nil {
				return nil, err
			}
			cards.Insert(c.GetID(), c)
		}
	}

	for i := 0; i < jokers; i++ {
		c, err := card.New("j", "j", 0, 0, false)
		if err != nil {
			return nil, err
		}
		cards.Insert(c.GetID(), c)
	}

	deck := &Deck{
		cards: cards,
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

	cas := safemap.New[string, *card.Card]()

	// 1d 4g 6g k!
	for i, c := range cards {
		ca, err := card.New(c[:1], c[1:], values[i], 1, false)
		if err != nil {
			return nil, err
		}

		cas.Insert(ca.GetID(), ca)
	}

	deck := &Deck{
		cards: cas,
	}

	return deck, nil
}

func (d *Deck) AddCard(card *card.Card) {
	d.cards.Insert(card.GetID(), card)
}

func (d *Deck) DelCard(card string) error {
	return d.cards.Delete(card)
}

func (d *Deck) HasCard(card string) bool {
	return d.cards.HasKey(card)
}
