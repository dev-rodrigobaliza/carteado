package game

import (
	"github.com/dev-rodrigobaliza/carteado/domain/core/game"
)

type IGame interface {
	GetMaxGroups() int
	GetMaxPlayersGroup() int
	GetMinPlayersGroup() int
	GetRound() uint64
	GetState() game.State
	Loop() (bool, error)
	SetState(game.State)
	Start() error
	Stop() error
	bet() (bool, error)
	deal() (bool, error)
	play() (bool, error)
	wait() (bool, error)
}
