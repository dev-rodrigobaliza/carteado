package core

type GameStatus struct {
	ID           string
	Owner        string
	Winners      []string
	MinPlayers   int
	MaxPlayers   int
	PlayerCount  int
	GameRound    uint64
	AllowBots    bool
	Private      bool
	GameMomentum GameMomentum
	GameState    GameState
	GameType     GameType
}
