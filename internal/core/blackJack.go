package core

import (
	"fmt"

	"github.com/dev-rodrigobaliza/carteado/errors"
	"github.com/dev-rodrigobaliza/carteado/utils"
)

type BlackJack struct {
	id           string
	owner        string
	winners      []string
	table        *Table
	gameMomentum GameMomentum
	gameState    GameState
	gameRound    uint64
}

// This line is for get feedback in case we are not implementing the interface correctly
var _ IGame = (*BlackJack)(nil)

func NewBlackJack(owner, secret string, minPlayers, maxPlayers int, allowBots bool) *BlackJack {
	id := "gid-" + utils.NewUUID()
	table := NewTable(owner, secret, minPlayers, maxPlayers, allowBots)

	return &BlackJack{
		id:           id,
		owner:        owner,
		winners:      make([]string, 0),
		table:        table,
		gameMomentum: GameMomentumWaiting,
		gameState:    GameStateStart,
		gameRound:    0,
	}
}

func (g *BlackJack) GetID() string {
	return g.id
}

func (g *BlackJack) GetOwner() string {
	return g.owner
}

func (g *BlackJack) GetStatus() *GameStatus {
	return &GameStatus{
		ID:           g.id,
		GameType:     GameTypeBlackJack,
		GameMomentum: g.gameMomentum,
		GameState:    g.gameState,
		GameRound:    g.gameRound,
		Owner:        g.owner,
		Winners:      g.winners,
		PlayerCount:  g.table.GetPlayersCount(),
		MinPlayers:   g.table.minPlayers,
		MaxPlayers:   g.table.maxPlayers,
		AllowBots:    g.table.allowBots,
		Private:      g.table.IsPrivate(),
	}
}

func (g *BlackJack) GetTable() *Table {
	return g.table
}

func (g *BlackJack) Enter(player, secret string) error {
	// basic validation
	if player == "" {
		return errors.ErrNotFoundPlayer
	}
	if !g.table.CheckSecret(secret) {
		return errors.ErrInvalidPassword
	}
	// add player
	err := g.table.AddPlayer(player)
	if err != nil {
		return err
	}

	return nil
}

func (g *BlackJack) Leave(player string) error {
	// basic validation
	if player == "" {
		return errors.ErrNotFoundPlayer
	}
	// del player
	err := g.table.DelPlayer(player)
	if err != nil {
		return err
	}
	// adjust owner
	if g.owner == player {
		g.owner = ""
	}

	return nil
}

func (g *BlackJack) Play() error {
	playersCount := g.table.GetPlayersCount()
	if playersCount < g.table.minPlayers {
		return fmt.Errorf("%d players not enough, game with %d min players", playersCount, g.table.minPlayers)
	}
	if g.gameState != GameStateStart {
		return fmt.Errorf("game already started")
	}
	// TODO (@dev-rodrigobaliza) more game conditions to start ???
	// TODO (@dev-rodrigobaliza) set players order

	g.gameState = GameStatePlay
	g.gameMomentum = GameMomentumCycling

	return nil
}

func (g *BlackJack) Stop(force bool) error {
	// TODO (@dev-rodrigobaliza) check game conditions to finish if not force
	g.gameState = GameStateFinish

	return nil
}
