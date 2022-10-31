package game

import (
	"github.com/dev-rodrigobaliza/carteado/domain/core/game"
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
	Stop() error
	bet() (bool, error)
	deal() (bool, error)
	play() (bool, error)
	wait() (bool, error)
}
