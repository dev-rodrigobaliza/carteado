package table

import (
	"github.com/dev-rodrigobaliza/carteado/domain/config"
	"github.com/dev-rodrigobaliza/carteado/domain/request"
	"github.com/dev-rodrigobaliza/carteado/internal/api/v1/adaptors/services"
	"github.com/dev-rodrigobaliza/carteado/internal/security/paseto"
	"github.com/dev-rodrigobaliza/carteado/pkg/safemap"
)

var (
	Security *paseto.PasetoMaker
)

type TableManager struct {
	cfg        *config.App
	appService *services.AppService
	tables     *safemap.SafeMap[string, *Table]
	players    *safemap.SafeMap[string, *Player]
}

func NewTableManager(cfg *config.App, players *safemap.SafeMap[string, *Player], appService *services.AppService) *TableManager {
	return &TableManager{
		cfg:        cfg,
		tables:     safemap.New[string, *Table](),
		players:    players,
		appService: appService,
	}
}

func (g *TableManager) ProcessPlayerMessage(player *Player, message request.WSRequest) {
	if message.Service == "auth" {
		g.serviceAuth(player, &message)
		return
	}

	if player.user == nil {
		g.sendResponseError(player, &message, "player unauthenticated", nil)
		return
	}

	switch message.Service {
	case "table":
		g.serviceTable(player, &message)

	default:
		g.sendResponseError(player, &message, "service not found", nil)
	}
}

func (g *TableManager) addTable(table *Table) {
	g.tables.Insert(table.GetID(), table)
}

func (g *TableManager) delTable(table *Table) error {
	return g.tables.Delete(table.GetID())
}

func (g *TableManager) getTable(id string) (*Table, error) {
	return g.tables.GetOneValue(id)
}
