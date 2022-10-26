package websocket

import (
	"github.com/dev-rodrigobaliza/carteado/domain/config"
	"github.com/dev-rodrigobaliza/carteado/domain/request"
	"github.com/dev-rodrigobaliza/carteado/internal/api/v1/adaptors/services"
	"github.com/dev-rodrigobaliza/carteado/internal/core"
	"github.com/dev-rodrigobaliza/carteado/internal/security/paseto"
	"github.com/dev-rodrigobaliza/carteado/pkg/safemap"
)

var (
	Security *paseto.PasetoMaker
)

type GameProcessor struct {
	cfg        *config.App
	games      *safemap.SafeMap[string, core.IGame]
	players    *safemap.SafeMap[*Player, bool]
	appService *services.AppService
}

func NewGameProcessor(cfg *config.App, players *safemap.SafeMap[*Player, bool], appService *services.AppService) *GameProcessor {
	return &GameProcessor{
		cfg:        cfg,
		games:      safemap.New[string, core.IGame](),
		players:    players,
		appService: appService,
	}
}

func (g *GameProcessor) ProcessPlayerMessage(player *Player, message request.WSRequest) {
	if message.Service == "auth" {
		g.serviceAuth(player, &message)
		return
	}

	if player.user == nil {
		g.sendResponseError(player, &message, "player unauthenticated", nil)
		return
	}

	switch message.Service {
	case "game":
		g.serviceGame(player, &message)

	default:
		g.sendResponseError(player, &message, "service not found", nil)
	}
}

func (g *GameProcessor) addGame(game core.IGame) {
	g.games.Insert(game.GetID(), game)
}

func (g *GameProcessor) delGame(game core.IGame) error {
	return g.games.Delete(game.GetID())
}

func (g *GameProcessor) getGame(id string) (core.IGame, error) {
	return g.games.GetOneValue(id)
}
