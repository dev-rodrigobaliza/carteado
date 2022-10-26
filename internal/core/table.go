package core

import (
	"errors"
	"sync"
)

type Table struct {
	minPlayers  int
	maxPlayers  int
	allowBots   bool
	players     map[string]bool
	playersLock sync.RWMutex
}

func NewTable(owner string, minPlayers, maxPlayers int, allowBots bool) *Table {
	table := &Table{
		allowBots:   allowBots,
		minPlayers:  minPlayers,
		maxPlayers:  maxPlayers,
		players:     make(map[string]bool),
		playersLock: sync.RWMutex{},
	}
	table.AddPlayer(owner)

	return table
}

func (t *Table) AddPlayer(playerID string) error {
	if t.HasPlayer(playerID) {
		return errors.New("player already exists")
	}

	t.playersLock.Lock()
	defer t.playersLock.Unlock()

	if len(t.players) >= t.maxPlayers {
		return errors.New("table full")
	}

	t.players[playerID] = true

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

func (t *Table) DelPlayer(playerID string) error {
	if !t.HasPlayer(playerID) {
		return errors.New("player not found")
	}

	t.playersLock.Lock()
	defer t.playersLock.Unlock()

	delete(t.players, playerID)

	if len(t.players) < t.minPlayers {
		return errors.New("players count above minimum")
	}

	return nil
}

func (t *Table) GetPlayers() []string {
	players := make([]string, 0)

	t.playersLock.RLock()
	defer t.playersLock.RUnlock()

	for player := range t.players {
		players = append(players, player)
	}

	return players
}

func (t *Table) GetPlayersCount() int {
	t.playersLock.RLock()
	defer t.playersLock.RUnlock()

	return len(t.players)
}

func (t *Table) HasPlayer(playerID string) bool {
	t.playersLock.RLock()
	defer t.playersLock.RUnlock()

	for player := range t.players {
		if player == playerID {
			return true
		}
	}

	return false
}
