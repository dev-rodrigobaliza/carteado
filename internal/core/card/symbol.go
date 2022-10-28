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

func (s Symbol) Graphic() string {
	switch s {
	case SymbolHearts:
		return "♥"

	case SymbolDiamonds:
		return "♦"

	case SymbolClubs:
		return "♣"

	case SymbolSpades:
		return "♠"

	case SymbolJoker:
		return "!"

	default:
		return "?"
	}
}

func (s Symbol) String() string {
	switch s {
	case SymbolHearts:
		return "hearts"

	case SymbolDiamonds:
		return "diamonds"

	case SymbolClubs:
		return "clubs"

	case SymbolSpades:
		return "spades"

	case SymbolJoker:
		return "joker"

	default:
		return "unknown"
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