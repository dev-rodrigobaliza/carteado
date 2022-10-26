package core

import (
	"errors"

	"github.com/dev-rodrigobaliza/carteado/utils"
)

type BlackJack struct {
	id        string
	owner     string
	winners   []string
	table     *Table
	gameState GameState
	gameRound uint64
}

// This line is for get feedback in case we are not implementing the interface correctly
var _ IGame = (*BlackJack)(nil)

func NewBlackJack(owner string, minPlayers, maxPlayers int, allowBots bool) *BlackJack {
	id := "gid-" + utils.NewUUID()
	table := NewTable(owner, minPlayers, maxPlayers, allowBots)

	return &BlackJack{
		id:        id,
		owner:     owner,
		winners:   make([]string, 0),
		table:     table,
		gameState: GameStateStarting,
		gameRound: 0,
	}
}

func (g *BlackJack) GetStatus() *GameStatus {
	return &GameStatus{
		ID:          g.id,
		GameType:    GameTypeBlackJack,
		GameState:   g.gameState,
		GameRound:   g.gameRound,
		Owner:       g.owner,
		Winners:     g.winners,
		PlayerCount: g.table.GetPlayersCount(),
		MinPlayers:  g.table.minPlayers,
		MaxPlayers:  g.table.maxPlayers,
		AllowBots:   g.table.allowBots,
	}
}

func (g *BlackJack) GetTable() *Table {
	return g.table
}

func (g *BlackJack) Start() error {
	// TODO (@dev-rodrigobaliza) set players order, set game state
	return errors.New("not implemented")
}
