package core

import (
	"errors"

	"github.com/dev-rodrigobaliza/carteado/pkg/safemap"
)

type Table struct {
	secret     string
	minPlayers int
	maxPlayers int
	allowBots  bool
	players    *safemap.SafeMap[string, bool]
}

func NewTable(owner, secret string, minPlayers, maxPlayers int, allowBots bool) *Table {
	table := &Table{
		secret:     secret,
		minPlayers: minPlayers,
		maxPlayers: maxPlayers,
		allowBots:  allowBots,
		players:    safemap.New[string, bool](),
	}
	table.AddPlayer(owner)

	return table
}

func (t *Table) AddPlayer(playerID string) error {
	if t.HasPlayer(playerID) {
		return errors.New("player already exists")
	}

	if t.players.Size() >= t.maxPlayers {
		return errors.New("table full")
	}

	t.players.Insert(playerID, true)

	return nil
}

func (t *Table) CheckSecret(secret string) bool {
	return t.secret == "" || t.secret == secret
}

func (t *Table) DelPlayer(playerID string) error {
	err := t.players.Delete(playerID)
	if err != nil {
		return errors.New("player not found")
	}

	if t.players.Size() < t.minPlayers {
		return errors.New("players count above minimum")
	}

	return nil
}

func (t *Table) GetAllowBots() bool {
	return t.allowBots
}

func (t *Table) GetMaxPlayers() int {
	return t.maxPlayers
}

func (t *Table) GetMinPlayers() int {
	return t.minPlayers
}

func (t *Table) GetPlayers() []string {
	return t.players.GetAllKeys()
}

func (t *Table) GetPlayersCount() int {
	return t.players.Size()
}

func (t *Table) HasPlayer(playerID string) bool {
	return t.players.HasKey(playerID)
}

func (t *Table) IsPrivate() bool {
	return t.secret != ""
}
