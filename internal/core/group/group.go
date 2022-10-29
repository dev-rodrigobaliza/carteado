package group

import (
	"github.com/dev-rodrigobaliza/carteado/domain/response"
	"github.com/dev-rodrigobaliza/carteado/errors"
	"github.com/dev-rodrigobaliza/carteado/internal/core/card"
	"github.com/dev-rodrigobaliza/carteado/internal/core/deck"
	"github.com/dev-rodrigobaliza/carteado/internal/core/player"
	"github.com/dev-rodrigobaliza/carteado/pkg/safemap"
)

type Group struct {
	id         int
	minPlayers int
	maxPlayers int
	players    *safemap.SafeMap[string, *player.Player]
	decks      *safemap.SafeMap[string, *deck.Deck]
}

func New(id, minPlayers, maxPlayers int) *Group {
	group := &Group{
		id:         id,
		minPlayers: minPlayers,
		maxPlayers: maxPlayers,
		players:    safemap.New[string, *player.Player](),
		decks:      safemap.New[string, *deck.Deck](),
	}

	return group
}

func (g *Group) AddCard(player string, card *card.Card) error {
	// basic validation
	if !g.decks.HasKey(player) {
		return errors.ErrNotFoundPlayer
	}
	deck, _ := g.decks.GetOneValue(player)
	deck.AddCard(card)

	return nil
}

func (g *Group) AddPlayer(player *player.Player) error {
	// basic validation
	if player == nil {
		return errors.ErrNotFoundPlayer
	}
	// group validation
	if g.HasPlayer(player.UUID) {
		return errors.ErrExistsPlayer
	}
	if g.players.Size() >= g.maxPlayers {
		return errors.ErrMaxPlayers
	}
	// add player
	g.players.Insert(player.UUID, player)
	player.GroupID = g.id

	return nil
}

func (g *Group) DelCard(player, card string) error {
	// basic validation
	if !g.decks.HasKey(player) {
		return errors.ErrNotFoundPlayer
	}
	deck, _ := g.decks.GetOneValue(player)
	if !deck.HasCard(card) {
		return errors.ErrNotFoundCard
	}

	return deck.DelCard(card)
}

func (g *Group) DelPlayer(player string) error {
	// basic validation
	if player == "" || !g.players.HasKey(player) {
		return errors.ErrNotFoundPlayer
	}
	// del player
	p, _ := g.players.GetOneValue(player)
	err := g.players.Delete(player)
	if err != nil {
		return errors.ErrNotFoundPlayer
	}
	p.GroupID = 0
	// group validation
	if g.players.Size() < g.minPlayers {
		return errors.ErrMinPlayers
		// TODO (@dev-rodrigobaliza) should stop the game ???
	}

	return nil
}

func (g *Group) GetID() int {
	return g.id
}

func (g *Group) GetMaxPlayers() int {
	return g.maxPlayers
}

func (g *Group) GetPlayers() []string {
	return g.players.GetAllKeys()
}

func (g *Group) GetPlayersCount() int {
	return g.players.Size()
}

func (g *Group) HasPlayer(playerID string) bool {
	return g.players.HasKey(playerID)
}

func (g *Group) ToResponse() *response.Group {
	pls := make([]*response.Player, 0)
	players := g.players.GetAllValues()
	for _, player := range players {
		pls = append(pls, player.ToResponse())
	}

	gr := response.NewGroup(g.id, g.minPlayers, g.maxPlayers, pls)

	return gr
}
