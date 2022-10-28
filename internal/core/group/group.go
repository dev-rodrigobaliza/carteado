package group

import (
	"github.com/dev-rodrigobaliza/carteado/consts"
	"github.com/dev-rodrigobaliza/carteado/domain/response"
	"github.com/dev-rodrigobaliza/carteado/errors"
	"github.com/dev-rodrigobaliza/carteado/internal/core/player"
	"github.com/dev-rodrigobaliza/carteado/pkg/safemap"
	"github.com/dev-rodrigobaliza/carteado/utils"
)

type Group struct {
	id         string
	players    *safemap.SafeMap[string, *player.Player]
	minPlayers int
	maxPlayers int
}

func New(minPlayers, maxPlayers int) *Group {
	group := &Group{
		id:         utils.NewUUID(consts.GROUP_PREFIX_ID),
		players:    safemap.New[string, *player.Player](),
		minPlayers: minPlayers,
		maxPlayers: maxPlayers,
	}

	return group
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

	return nil
}

func (g *Group) DelPlayer(player string) error {
	// basic validation
	if player == "" {
		return errors.ErrNotFoundPlayer
	}
	// del player
	err := g.players.Delete(player)
	if err != nil {
		return errors.ErrNotFoundPlayer
	}
	// group validation
	if g.players.Size() < g.minPlayers {
		return errors.ErrMinPlayers
		// TODO (@dev-rodrigobaliza) should stop the game ???
	}

	return nil
}

func (g *Group) GetID() string {
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

	gr := response.NewGroup(g.id, pls)

	return gr
}
