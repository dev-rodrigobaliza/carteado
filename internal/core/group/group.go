package group

import (
	"fmt"
	"time"

	"github.com/dev-rodrigobaliza/carteado/domain/core/group"
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
	State      group.State
	NextPlayer int
	createdAt  time.Time
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
	if player != "" && !g.decks.HasKey(player) {
		return errors.ErrNotFoundPlayer
	}
	// empty player means first one (probably unique player)
	if player == "" {
		var err error
		player, err = g.players.GetFirstKey()
		if err != nil {
			return err
		}
	}

	deck, _ := g.decks.GetOneValue(player, false)
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
	// check if its time to populate createdAt
	if g.GetPlayersCount() == 0 {
		g.createdAt = time.Now()
	}
	// add player
	g.players.Insert(player.UUID, player)
	player.GroupID = g.id
	// add empty deck
	g.decks.Insert(player.UUID, deck.New())

	return nil
}

func (g *Group) DelCard(player, card string) error {
	// basic validation
	if !g.decks.HasKey(player) {
		return errors.ErrNotFoundPlayer
	}
	deck, _ := g.decks.GetOneValue(player, false)
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
	p, _ := g.players.GetOneValue(player, false)
	err := g.players.DeleteKey(player)
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

func (g *Group) GetGroupScore() int {
	score := 0

	for _, d := range g.decks.GetAllValues() {
		score += d.GetScore(true)
	}

	return score
}

func (g *Group) GetID() int {
	return g.id
}

func (g *Group) GetMaxPlayers() int {
	return g.maxPlayers
}

func (g *Group) GetNextDeck() (*deck.Deck, error) {
	d, err := g.decks.GetIndexedValue(g.NextPlayer, false)
	if err != nil {
		return nil, errors.ErrNotFoundDeck
	}

	return d, nil
}

func (g *Group) GetNextPlayer() (*player.Player, error) {
	p, err := g.players.GetIndexedValue(g.NextPlayer, false)
	if err != nil {
		return nil, errors.ErrNotFoundPlayer
	}

	return p, nil
}

func (g *Group) GetPlayerCards(player string) ([]*card.Card, error) {
	// basic validation
	if player == "" || !g.players.HasKey(player) {
		return nil, errors.ErrNotFoundPlayer
	}
	// get player deck
	d, err := g.decks.GetOneValue(player, false)
	if err != nil {
		return nil, err
	}
	// get player cards from deck
	var cards []*card.Card
	for _, c := range d.GetAllCards() {
		cards = append(cards, c)
	}

	return cards, nil
}

func (g *Group) GetPlayerDeck(player string) (*deck.Deck, error) {
	// basic validation
	if player == "" || !g.players.HasKey(player) {
		return nil, errors.ErrNotFoundPlayer
	}
	// get player deck
	d, err := g.decks.GetOneValue(player, false)
	if err != nil {
		return nil, err
	}

	return d, nil
}

func (g *Group) GetPlayers() []string {
	return g.players.GetAllKeys()
}

func (g *Group) GetPlayerScore(player string) (int, error) {
	// basic validation
	if player == "" || !g.players.HasKey(player) {
		return 0, errors.ErrNotFoundPlayer
	}
	// get player deck
	d, err := g.decks.GetOneValue(player, false)
	if err != nil {
		return 0, err
	}
	// get player cards from deck
	score := 0
	for _, c := range d.GetAllCards() {
		score += c.Value(true)
	}

	return score, nil
}

func (g *Group) GetPlayersCount() int {
	return g.players.Size()
}

func (g *Group) HasPlayer(playerID string) bool {
	return g.players.HasKey(playerID)
}

func (g *Group) ToResponse(full, admin bool) *response.Group {
	pls := make([]*response.Player, 0)
	players := g.players.GetAllValues()
	for _, player := range players {
		pls = append(pls, player.ToResponse(full, admin))
	}

	var created string
	if g.GetPlayersCount() > 0 {
		created = fmt.Sprintf("%d", g.createdAt.UnixMilli())
	}

	gr := response.NewGroup(g.id, g.minPlayers, g.maxPlayers, pls, created)

	return gr
}
