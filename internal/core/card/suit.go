package card

import (
	"github.com/dev-rodrigobaliza/carteado/errors"
)

type Suit struct {
	Symbol Symbol
	Value  int
}

func NewSuit(suit string) (*Suit, error) {
	if suit == "" {
		return nil, errors.ErrInvalidCardSuit
	}

	symbol := NewSymbol(suit, false)
	if symbol == SymbolUnknown {
		return nil, errors.ErrInvalidCardSuit
	}

	s := &Suit{
		Symbol: symbol,
		Value:  int(symbol),
	}

	return s, nil
}

func NewSuitCustom(suit string, value int) (*Suit, error) {
	s, err := NewSuit(suit)
	if err != nil {
		return nil, err
	}

	s.Value = value

	return s, nil
}

func NewSuitRandom(joker bool) *Suit {
	symbol := RandomSymbol(joker)

	suit := &Suit{
		Symbol: symbol,
		Value: int(symbol),
	}

	return suit
}

func (s *Suit) String() string {
	return s.Symbol.String()
}
