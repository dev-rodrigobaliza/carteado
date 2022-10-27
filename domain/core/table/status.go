package table

import "github.com/dev-rodrigobaliza/carteado/domain/core/game"

type Status struct {
	ID          string
	Owner       string
	Winners     []string
	MinPlayers  int
	MaxPlayers  int
	PlayerCount int
	GameRound   uint64
	AllowBots   bool
	Private     bool
	State       State
	GameMode    game.Mode
	GameState   game.State
}
