package card

import "math/rand"

type Symbol int

const (
	SymbolUnknown Symbol = iota
	SymbolHearts
	SymbolDiamonds
	SymbolClubs
	SymbolSpades
	SymbolJoker
)

func NewSymbol(symbol string, joker bool) Symbol {
	if symbol == "" {
		return RandomSymbol(joker)
	}

	switch symbol {
	case "h":
		return SymbolHearts

	case "d":
		return SymbolDiamonds

	case "c":
		return SymbolClubs

	case "s":
		return SymbolSpades

	case "!":
		return SymbolJoker

	default:
		return SymbolUnknown
	}
}

func (s Symbol) String() string {
	switch s {
	case SymbolHearts:
		return "h"

	case SymbolDiamonds:
		return "d"

	case SymbolClubs:
		return "c"

	case SymbolSpades:
		return "s"

	case SymbolJoker:
		return "!"

	default:
		return "?"
	}
}

func RandomSymbol(joker bool) Symbol {
	var limiter int
	if joker {
		limiter = int(SymbolJoker)
	} else {
		limiter = int(SymbolJoker) - 1
	}

	return Symbol(rand.Intn(limiter) + 1)
}