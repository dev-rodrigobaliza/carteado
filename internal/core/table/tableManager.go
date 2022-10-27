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

func (t *TableManager) ProcessPlayerMessage(player *Player, message request.WSRequest) {
	if message.Service == "auth" {
		t.serviceAuth(player, &message)
		return
	}

	if player.user == nil {
		t.sendResponseError(player, &message, "player unauthenticated", nil)
		return
	}

	switch message.Service {
	case "admin":
		t.serviceAdmin(player, &message)

	case "table":
		t.serviceTable(player, &message)

	default:
		t.sendResponseError(player, &message, "service not found", nil)
	}
}

func (t *TableManager) addTable(table *Table) {
	t.tables.Insert(table.GetID(), table)
}

func (t *TableManager) delTable(table *Table) error {
	return t.tables.Delete(table.GetID())
}

func (t *TableManager) getTable(id string) (*Table, error) {
	return t.tables.GetOneValue(id)
}
