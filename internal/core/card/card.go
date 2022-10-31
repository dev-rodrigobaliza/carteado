package card

import (
	"github.com/dev-rodrigobaliza/carteado/consts"
	"github.com/dev-rodrigobaliza/carteado/errors"
	"github.com/dev-rodrigobaliza/carteado/utils"
)

type Card struct {
	id      string
	face    *Face
	suit    *Suit
	Visible bool
}

func New(face, suit string, visible bool) (*Card, error) {
	f, err := NewFace(face)
	if err != nil {
		return nil, errors.ErrInvalidCardFace
	}

	s, err := NewSuit(suit)
	if err != nil {
		return nil, errors.ErrInvalidCardSuit
	}

	card := &Card{
		id:      utils.NewUUID(consts.CARD_PREFIX_ID),
		face:    f,
		suit:    s,
		Visible: visible,
	}

	return card, nil
}

func NewCustom(face, suit string, faceValue, suitValue int, joker, visible bool) (*Card, error) {
	f, err := NewFaceCustom(face, faceValue)
	if err != nil {
		return nil, errors.ErrInvalidCardFace
	}

	s, err := NewSuitCustom(suit, suitValue)
	if err != nil {
		return nil, errors.ErrInvalidCardSuit
	}

	card := &Card{
		id:      utils.NewUUID(consts.CARD_PREFIX_ID),
		face:    f,
		suit:    s,
		Visible: visible,
	}

	return card, nil
}

func NewRandom(joker, visible bool) *Card {
	card := &Card{
		id:      utils.NewUUID(consts.CARD_PREFIX_ID),
		face:    NewFaceRandom(joker),
		suit:    NewSuitRandom(joker),
		Visible: visible,
	}

	return card
}

func (c *Card) GetID() string {
	return c.id
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

func (c *Card) Value(suit bool) int {
	value := c.face.Value

	if suit {
		value = value * c.suit.Value
	}

	return value
}
