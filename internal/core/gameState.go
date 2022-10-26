package core

type GameState int8

func (g GameState) String() string {
	return [...]string{"unknown", "start", "play", "finish"}[g]
}

func StringToGameState(gameState string) GameState {
	switch gameState {
	case "start":
		return GameStateStart

	case "play":
		return GameStatePlay

	case "finish":
		return GameStateFinish

	default:
		return GameStateUnknown
	}
}

const (
	GameStateUnknown GameState = iota
	GameStateStart
	GameStatePlay
	GameStateFinish
)
