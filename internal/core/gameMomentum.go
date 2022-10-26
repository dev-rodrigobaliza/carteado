package core

type GameMomentum int8

func (g GameMomentum) String() string {
	return [...]string{"unknown", "betting", "cycling", "waiting"}[g]
}

func StringToGameMomentum(gameState string) GameMomentum {
	switch gameState {
	case "start":
		return GameMomentumBetting

	case "play":
		return GameMomentumCycling

	case "finish":
		return GameMomentumWaiting

	default:
		return GameMomentumUnknown
	}
}

const (
	GameMomentumUnknown GameMomentum = iota
	GameMomentumBetting
	GameMomentumCycling
	GameMomentumWaiting
)
