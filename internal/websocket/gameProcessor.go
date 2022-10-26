package websocket

import (
	"sync"

	"github.com/dev-rodrigobaliza/carteado/domain/config"
	"github.com/dev-rodrigobaliza/carteado/domain/request"
	"github.com/dev-rodrigobaliza/carteado/internal/api/v1/ports"
	"github.com/dev-rodrigobaliza/carteado/internal/core"
	"github.com/dev-rodrigobaliza/carteado/internal/security/paseto"
)

var (
	Security    *paseto.PasetoMaker
	AuthService ports.IAuthService
)

type GameProcessor struct {
	cfg        *config.App
	games      map[string]core.IGame
	gamesCount int64
	gamesLock  *sync.RWMutex
}

func NewGameProcessor(cfg *config.App) *GameProcessor {
	return &GameProcessor{
		cfg:        cfg,
		games:      make(map[string]core.IGame),
		gamesCount: 0,
		gamesLock:  &sync.RWMutex{},
	}
}

func (g *GameProcessor) Run(player *Player, message request.WSRequest) {
	if message.Service == "auth" {
		g.serviceAuth(player, &message)
		return
	}

	if !player.auth {
		g.sendResponseError(player, &message, "player unauthenticated")
		return
	}

	switch message.Service {
	case "game":
		g.serviceGame(player, &message)

	default:
		g.sendResponseError(player, &message, "service not found")
	}
}

func (g *GameProcessor) addGame(game core.IGame) {
	g.gamesLock.Lock()
	defer g.gamesLock.Unlock()

	g.games[game.GetStatus().ID] = game
	g.gamesCount++
}

func (g *GameProcessor) delGame(game core.IGame) {
	// TODO (@dev-rodrigobaliza) notify all players the game will be removed???
	g.gamesLock.Lock()
	defer g.gamesLock.Unlock()

	delete(g.games, game.GetStatus().ID)
	g.gamesCount--
}

func (g *GameProcessor) getGame(id string) core.IGame {
	g.gamesLock.RLock()
	defer g.gamesLock.RUnlock()

	return g.games[id]
}
