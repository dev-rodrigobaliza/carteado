package core

type GameState int8

func (g GameState) String() string {
	return [...]string{"unknown", "starting", "betting", "playing", "finished"}[g]
}

func StringToGameState(gameState string) GameState {
	switch gameState {
	case "starting":
		return GameStateStarting

	case "betting":
		return GameStateBetting

	case "playing":
		return GameStatePlaying

	case "finished":
		return GameStateFinished

	default:
		return GameStateUnknown
	}
}

const (
	GameStateUnknown GameState = iota
	GameStateStarting
	GameStateBetting
	GameStatePlaying
	GameStateFinished
)
