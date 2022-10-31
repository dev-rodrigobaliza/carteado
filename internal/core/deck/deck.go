package deck

import (
	"math/rand"

	"github.com/dev-rodrigobaliza/carteado/domain/response"
	"github.com/dev-rodrigobaliza/carteado/errors"
	"github.com/dev-rodrigobaliza/carteado/internal/core/card"
	"github.com/dev-rodrigobaliza/carteado/pkg/safemap"
	"github.com/dev-rodrigobaliza/carteado/vars"
)

type Deck struct {
	cards *safemap.SafeMap[string, *card.Card]
}

func New() *Deck {
	deck := &Deck{
		cards: safemap.New[string, *card.Card](),
	}

	return deck
}

func NewCustom(cards []string, values []int, visible bool) (*Deck, error) {
	if len(cards) == 0 {
		return nil, errors.ErrInvalidCard
	}
	if len(values) == 0 {
		return nil, errors.ErrInvalidCardValue
	}
	if len(cards) != len(values) {
		return nil, errors.ErrInvalidCardDeck
	}

	deck := New()

	for i, c := range cards {
		ca, err := card.NewCustom(c[:1], c[1:], values[i], 1, false, visible)
		if err != nil {
			return nil, err
		}

		deck.AddCard(ca)
	}

	return deck, nil
}

func NewDefault(jokers int, visible bool) (*Deck, error) {
	deck := New()

	for _, face := range vars.CardFaces {
		for _, suit := range vars.CardSuits {
			c, err := card.New(face, suit, visible)
			if err != nil {
				return nil, err
			}

			deck.AddCard(c)
		}
	}

	for i := 0; i < jokers; i++ {
		c, err := card.New("j", "j", visible)
		if err != nil {
			return nil, err
		}

		deck.AddCard(c)
	}

	return deck, nil
}

func (d *Deck) AddCard(card *card.Card) {
	d.cards.Insert(card.GetID(), card)
}

func (d *Deck) AddCardRandom(joker, visible bool) {
	card := card.NewRandom(joker, visible)
	d.cards.Insert(card.GetID(), card)
}

func (d *Deck) DelCard(card string) error {
	return d.cards.DeleteKey(card)
}

func (d *Deck) GetAllCards() []*card.Card {
	return d.cards.GetAllValues()
}

func (d *Deck) GetFirstCard(remove bool) (*card.Card, error) {
	return d.cards.GetFirstValue(remove)
}

func (d *Deck) GetLastCard(remove bool) (*card.Card, error) {
	return d.cards.GetLastValue(remove)
}

func (d *Deck) GetRandomCard(remove bool) (*card.Card, error) {
	if d.cards.Size() == 0 {
		return nil, errors.ErrEmptyDeck
	}

	var c *card.Card
	index := rand.Intn(d.cards.Size())
	for i, card := range d.cards.GetAllValues() {
		if i == index {
			c = card
		}
	}

	return c, nil
}

func (d *Deck) GetScore(suit bool) int {
	score := 0

	for _, c := range d.cards.GetAllValues() {
		score += c.Value(suit)
	}

	return score
}

func (d *Deck) HasCard(card string) bool {
	return d.cards.HasKey(card)
}

func (d *Deck) Shuffle(times int) {
	cards := d.cards.GetAllValues()

	for i := 0; i < times; i++ {
		for j := 0; j < len(cards); j++ {
			r := rand.Intn(j + 1)
			if r != j {
				cards[j], cards[r] = cards[r], cards[j]
			}
		}
	}

	d.cards.Clear()
	for _, card := range cards {
		d.cards.Insert(card.GetID(), card)
	}
}

func (d *Deck) ToResponse(tableID string, visibleOnly bool) *response.Deck {
	cards := make([]*response.Card, 0)
	for _, ca := range d.cards.GetAllValues() {
		var cardStr string
		if ca.Visible {
			cardStr = ca.String(true, true)
		} else {
			cardStr = ""
		}
		car := response.NewCard(cardStr, ca.Value(true))
		cards = append(cards, car)
	}

	return response.NewDeck(tableID, cards)
}
