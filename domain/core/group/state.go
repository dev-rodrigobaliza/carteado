package group

type State int8

func (s State) String() string {
	return [...]string{"unknown", "ready", "card", "action", "stop", "bet", "finish"}[s]
}

func StringToState(groupState string) State {
	switch groupState {
	case "ready":
		return StateReady

	case "card":
		return StateCard

	case "action":
		return StateAction

	case "stop":
		return StateStop

	case "bet":
		return StateBet

	case "finish":
		return StateFinish

	default:
		return StateUnknown
	}
}

const (
	StateUnknown State = iota
	StateReady
	StateCard
	StateAction
	StateStop
	StateBet
	StateFinish
)
