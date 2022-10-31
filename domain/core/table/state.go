package table

type State int8

func (s State) String() string {
	return [...]string{"unknown", "start", "play", "finish"}[s]
}

func StringToState(gameState string) State {
	switch gameState {
	case "start":
		return StateStart

	case "play":
		return StatePlay

	case "finish":
		return StateFinish

	default:
		return StateUnknown
	}
}

const (
	StateUnknown State = iota
	StateStart
	StatePlay
	StateFinish
)
