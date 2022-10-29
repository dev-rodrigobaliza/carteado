package card

import (
	"github.com/dev-rodrigobaliza/carteado/consts"
	"github.com/dev-rodrigobaliza/carteado/errors"
	"github.com/dev-rodrigobaliza/carteado/utils"
)

type Card struct {
	id   string
	face *Face
	suit *Suit
}

func New(face, suit string, faceValue, suitValue int, joker bool) (*Card, error) {
	f, err := NewFace(face, faceValue, joker)
	if err != nil {
		return nil, errors.ErrInvalidCardFace
	}

	s, err := NewSuit(suit, suitValue, joker)
	if err != nil {
		return nil, errors.ErrInvalidCardSuit
	}

	if f.Height == HeightJoker || s.Symbol == SymbolJoker {
		f.Height = HeightJoker
		s.Symbol = SymbolJoker
	}

	card := &Card{
		id:   utils.NewUUID(consts.CARD_PREFIX_ID),
		face: f,
		suit: s,
	}

	return card, nil
}

func (c *Card) GetID() string {
	return c.id
}

func (c *Card) Graphic(face bool) string {
	str := ""

	if face {
		str += c.face.String()
	}

	str += c.suit.Graphic()

	return str
}

func (c *Card) String(face, suit bool) string {
	str := ""

	if face {
		str += c.face.String()
	}

	if suit {
		str += c.suit.String()
	}

	return str
}

func (c *Card) Value(face, suit bool) int {
	value := 0

	if face {
		value += c.face.Value
	}

	if suit {
		value += c.suit.Value
	}

	return value
}
