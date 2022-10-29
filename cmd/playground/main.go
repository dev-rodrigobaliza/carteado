package main

import (
	"github.com/dev-rodrigobaliza/carteado/internal/core/card"
)

func main() {
	cardJokerForcedFace, _ := card.New("1", "!", 0, 10, false)
	println(cardJokerForcedFace.Value(true, true))
}
