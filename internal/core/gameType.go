package core

type GameType int8

func (g GameType) String() string {
	return [...]string{"unknown", "blackjack"}[g]
}

func StringToGametype(gameType string) GameType {
	switch gameType {
	case "blackjack":
		return GameTypeBlackJack
		
	default:
		return GameTypeUnknown
	}
}

const (
	GameTypeUnknown GameType = iota
	GameTypeBlackJack
)
