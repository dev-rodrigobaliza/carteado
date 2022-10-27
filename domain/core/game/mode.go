package game

type Mode int8

func (g Mode) String() string {
	return [...]string{"unknown", "blackjack"}[g]
}

func StringToMode(gameType string) Mode {
	switch gameType {
	case "blackjack":
		return ModeBlackJack

	default:
		return ModeUnknown
	}
}

const (
	ModeUnknown Mode = iota
	ModeBlackJack
)
