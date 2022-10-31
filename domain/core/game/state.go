package game

type State int8

func (s State) String() string {
	return [...]string{"unknown", "betting", "dealing", "playing", "waiting"}[s]
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
