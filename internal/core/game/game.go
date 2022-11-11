package game

import (
	"github.com/dev-rodrigobaliza/carteado/domain/core/game"
	"github.com/dev-rodrigobaliza/carteado/domain/response"
	"github.com/dev-rodrigobaliza/carteado/errors"
	"github.com/dev-rodrigobaliza/carteado/internal/core/deck"
	"github.com/dev-rodrigobaliza/carteado/internal/core/group"
	"github.com/dev-rodrigobaliza/carteado/internal/core/player"
)

type IGame interface {
	GetActiveGroup() *group.Group
	GetActivePlayer() (*player.Player, error)
	GetActivePlayerDeck() (*deck.Deck, error)
	GetMaxGroups() int
	GetMaxPlayersGroup() int
	GetMinPlayersGroup() int
	GetRound() uint64
	GetState() game.State
	Loop() (string, bool, error)
	Response() *response.Game
	SetState(game.State)
	Start([]*group.Group) error
	Stop() error
	bet() (string, bool, error)
	deal() (string, bool, error)
	play() (string, bool, error)
	wait() (string, bool, error)
}

func New(gameMode game.Mode) (IGame, error) {
	var g IGame

	switch gameMode {
	case game.ModeBlackJack:
		g = NewBlackJack()

	default:
		return nil, errors.ErrInvalidGameMode
	}

	return g, nil
}