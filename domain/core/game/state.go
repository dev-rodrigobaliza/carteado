package game

type State int8

func (g State) String() string {
	return [...]string{"unknown", "betting", "dealing", "playing", "waiting"}[g]
}

func StringToState(gameState string) State {
	switch gameState {
	case "betting":
		return StateBetting

	case "dealing":
		return StateDealing

	case "playing":
		return StatePlaying

	case "waiting":
		return StateWaiting

	default:
		return StateUnknown
	}
}

const (
	StateUnknown State = iota
	StateBetting
	StateDealing
	StatePlaying
	StateWaiting
)
