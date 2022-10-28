package game

import (
	"github.com/dev-rodrigobaliza/carteado/consts"
	"github.com/dev-rodrigobaliza/carteado/domain/core/game"
	"github.com/dev-rodrigobaliza/carteado/errors"
)

type BlackJack struct {
	maxGroups       int
	maxPlayersGroup int
	minPlayersGroup int
	state           game.State
	round           uint64
}

// This line is for get feedback in case we are not implementing the interface correctly
var _ IGame = (*BlackJack)(nil)

func NewBlackJack() *BlackJack {
	return &BlackJack{
		maxGroups:       consts.GAME_BLACKJACK_MAX_GROUPS,
		maxPlayersGroup: consts.GAME_BLACKJACK_MAX_PLAYERS_GROUP,
		minPlayersGroup: consts.GAME_BLACKJACK_MIN_PLAYERS_GROUP,
		state:           game.StateWaiting,
		round:           0,
	}
}

func (g *BlackJack) GetMaxGroups() int {
	return g.maxGroups
}

func (g *BlackJack) GetMaxPlayersGroup() int {
	return g.maxPlayersGroup
}

func (g *BlackJack) GetMinPlayersGroup() int {
	return g.minPlayersGroup
}

func (g *BlackJack) GetRound() uint64 {
	return g.round
}

func (g *BlackJack) GetState() game.State {
	return g.state
}

func (g *BlackJack) Loop() (bool, error) {
	switch g.state {
	case game.StateDealing:
		return g.deal()

	case game.StateBetting:
		return g.bet()

	case game.StatePlaying:
		return g.play()

	case game.StateWaiting:
		return g.wait()
	}

	return false, errors.ErrInvalidGameState
}

func (g *BlackJack) SetState(gameState game.State) {
	g.state = gameState
}

func (g *BlackJack) Start() error {
	if g.state != game.StateWaiting {
		return errors.ErrInvalidGameState
	}
	// TODO (@dev-rodrigobaliza) prepare the game to start

	return nil
}

func (g *BlackJack) Stop() error {
	if g.state == game.StateWaiting {
		return errors.ErrInvalidGameState
	}
	// TODO (@dev-rodrigobaliza) finish all the loose ends

	return nil
}

func (g *BlackJack) bet() (bool, error) {
	return false, errors.ErrNotImplemented
}

func (g *BlackJack) deal() (bool, error) {
	return false, errors.ErrNotImplemented
}

func (g *BlackJack) play() (bool, error) {
	return false, errors.ErrNotImplemented
}

func (g *BlackJack) wait() (bool, error) {
	return false, errors.ErrNotImplemented
}
