package card

import "math/rand"

type Height int

const (
	HeightUnknown Height = iota
	HeightAce
	Height2
	Height3
	Height4
	Height5
	Height6
	Height7
	Height8
	Height9
	Height10
	HeightJack
	HeightQueen
	HeightKing
	HeightJoker
)

func NewHeight(height string, joker bool) Height {
	if height == "" {
		return RandomHeight(joker)
	}

	switch height {
	case "1":
		return HeightAce

	case "2":
		return Height2

	case "3":
		return Height3

	case "4":
		return Height4

	case "5":
		return Height5

	case "6":
		return Height6

	case "7":
		return Height7

	case "8":
		return Height8

	case "9":
		return Height9

	case "0":
		return Height10

	case "j":
		return HeightJack

	case "q":
		return HeightQueen

	case "k":
		return HeightKing

	case "!":
		return HeightJoker

	default:
		return HeightUnknown
	}
}

func (h Height) String() string {
	switch h {
	case HeightAce:
		return "1"

	case Height2:
		return "2"

	case Height3:
		return "3"

	case Height4:
		return "4"

	case Height5:
		return "5"

	case Height6:
		return "6"

	case Height7:
		return "7"

	case Height8:
		return "8"

	case Height9:
		return "9"

	case Height10:
		return "10"

	case HeightJack:
		return "j"

	case HeightQueen:
		return "q"

	case HeightKing:
		return "k"

	case HeightJoker:
		return "!"

	default:
		return "?"
	}
}

func RandomHeight(joker bool) Height {
	var limiter int
	if joker {
		limiter = int(HeightJoker)
	} else {
		limiter = int(HeightJoker) - 1
	}

	return Height(rand.Intn(limiter) + 1)
}
