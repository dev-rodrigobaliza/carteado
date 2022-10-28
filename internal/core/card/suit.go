package card

import (
	"github.com/dev-rodrigobaliza/carteado/errors"
)

type Suit struct {
	Symbol Symbol
	Value  int
}

func NewSuit(suit string, value int, joker bool) (*Suit, error) {
	var symbol Symbol

	if suit == "" {
		symbol = RandomSymbol(joker)
	} else {
		symbol = NewSymbol(suit, joker)
		if symbol == SymbolUnknown {
			return nil, errors.ErrInvalidCardSuit
		}
	}

	if value == 0 {
		value = int(symbol)
	}

	s := &Suit{
		Symbol: symbol,
		Value:  value,
	}

	return s, nil
}

func (s *Suit) Graphic() string {
	return s.Symbol.Graphic()
}

func (s *Suit) String() string {
	return s.Symbol.String()
}
