package saloon

import (
	"log"

	"github.com/dev-rodrigobaliza/carteado/domain/request"
	"github.com/dev-rodrigobaliza/carteado/domain/response"
	pl "github.com/dev-rodrigobaliza/carteado/internal/core/player"
	"github.com/dev-rodrigobaliza/carteado/internal/core/table"
)

func (s *Saloon) debug(format string, v ...any) {
	if s.cfg.Debug {
		log.Printf(format, v...)
	}
}

func (s *Saloon) getServerStatusResponse(authenticatedOnly bool) map[string]interface{} {
	// get tables
	tables := make([]*response.Table, 0)
	for _, table := range s.tables.GetAllValues() {
		tables = append(tables, table.ToResponse())
	}

	response := make(map[string]interface{})
	response["server"] = s.cfg.Name
	response["version"] = s.cfg.Version
	response["started_at"] = s.cfg.StartedAt
	response["players_count"] = s.players.Size()
	response["tables_count"] = len(tables)
	if len(tables) > 0 {
		response["tables"] = tables
	}

	return response
}

func (s *Saloon) getAction(message *request.WSRequest) string {
	action, ok := message.Data["action"].(string)
	if !ok {
		return ""
	}

	return action
}

func (s *Saloon) getGroupID(message *request.WSRequest) int {
	groupID, ok := message.Data["group_id"].(float64)
	if !ok {
		return 0
	}

	return int(groupID)
}

func (s *Saloon) getTableGroupStatus(table *table.Table, groupID int) map[string]interface{} {
	group, _ := table.GetGroup(groupID)
	response := make(map[string]interface{})
	response["table"] = group.ToResponse(false)

	return response
}

func (s *Saloon) getTableID(message *request.WSRequest) string {
	tableID, ok := message.Data["table_id"].(string)
	if !ok {
		return ""
	}

	return tableID
}

func (s *Saloon) getTableStatus(table *table.Table) map[string]interface{} {
	response := make(map[string]interface{})
	response["table"] = table.ToResponse()

	return response
}

func (s *Saloon) greetingMesssage(player *pl.Player, message string) {
	response := make(map[string]interface{})
	response["server"] = s.cfg.Name
	response["version"] = s.cfg.Version

	player.SendResponse(nil, "info", message, response)
}

func (s *Saloon) loginPlayer(player *pl.Player) {
	var welcomePlayer *pl.Player
	// if player has previous login and remove it
	players := s.players.GetAllValues()
	for _, p := range players {
		if p != player && p.User != nil && p.User.ID == player.User.ID {
			p.User = nil
			welcomePlayer = p
			break
		}
	}

	if welcomePlayer != nil {
		s.greetingMesssage(welcomePlayer, "disconnected (using another session)")
	}
}

func (s *Saloon) sendResponseError(player *pl.Player, request *request.WSRequest, message string, err error) {
	var data map[string]interface{}
	if err != nil && s.cfg.Debug {
		data = make(map[string]interface{})
		data["error"] = err.Error()
	}

	player.SendResponse(request, "error", message, data)
}

func (s *Saloon) sendResponseSuccess(player *pl.Player, request *request.WSRequest, message string, data map[string]interface{}) {
	player.SendResponse(request, "success", message, data)
}
