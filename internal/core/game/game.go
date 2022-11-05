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
	Loop() (bool, error)
	SetState(game.State)
	Start([]*group.Group) error
	ToResponse() *response.Game
	Stop() error
	bet() (bool, error)
	deal() (bool, error)
	play() (bool, error)
	wait() (bool, error)
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