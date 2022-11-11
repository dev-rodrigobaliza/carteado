package saloon

import (
	"fmt"
	"log"

	"github.com/dev-rodrigobaliza/carteado/domain/request"
	"github.com/dev-rodrigobaliza/carteado/domain/response"
	"github.com/dev-rodrigobaliza/carteado/internal/core/player"
	"github.com/dev-rodrigobaliza/carteado/internal/core/table"
)

func (s *Saloon) debug(format string, v ...any) {
	if s.cfg.Debug {
		log.Printf(format, v...)
	}
}

func (s *Saloon) getAction(req *request.WSRequest) string {
	action, ok := req.Data["action"].(string)
	if !ok {
		return ""
	}

	return action
}

func (s *Saloon) getGroupID(req *request.WSRequest) int {
	groupID, ok := req.Data["group_id"].(float64)
	if !ok {
		return 0
	}

	return int(groupID)
}

func (s *Saloon) getQuantity(req *request.WSRequest) int {
	quantity, ok := req.Data["quantity"].(float64)
	if !ok {
		return 0
	}

	return int(quantity)
}

func (s *Saloon) getServerStatusResponse(admin bool) map[string]interface{} {
	// get players out of table
	playersOut := make([]*response.Player, 0)
	for _, pl := range s.players.GetAllValues() {
		if pl.TableID == "" {
			playersOut = append(playersOut, pl.Response(true, false))
		}
	}
	// get tables
	tables := make([]*response.Table, 0)
	for _, tb := range s.tables.GetAllValues() {
		tables = append(tables, tb.Response(true))
	}
	// make reponse
	response := make(map[string]interface{})
	response["server"] = s.cfg.Name
	response["version"] = s.cfg.Version
	response["created_at"] = fmt.Sprintf("%d", s.cfg.CreatedAt.UnixMilli())
	response["players_count"] = s.players.Size()
	response["players_out_table_count"] = len(playersOut)
	if len(playersOut) > 0 {
		response["players_out_table"] = playersOut
	}
	response["tables_count"] = len(tables)
	if len(tables) > 0 {
		response["tables"] = tables
	}

	return response
}

func (s *Saloon) getTableGameStatus(tb *table.Table) map[string]interface{} {
	response := make(map[string]interface{})
	response["game_status"] = tb.GetGameStatus()

	return response
}

func (s *Saloon) getTableGroupStatus(tb *table.Table, groupID int, admin bool) map[string]interface{} {
	group, _ := tb.GetGroup(groupID)
	response := make(map[string]interface{})
	response["table"] = group.Response(admin)

	return response
}

func (s *Saloon) getTableID(req *request.WSRequest) string {
	tableID, ok := req.Data["table_id"].(string)
	if !ok {
		return ""
	}

	return tableID
}

func (s *Saloon) getTableStatus(tb *table.Table, admin bool) map[string]interface{} {
	response := make(map[string]interface{})
	response["table"] = tb.Response(admin)

	return response
}

func (s *Saloon) greetingMesssage(pl *player.Player, message string) {
	response := make(map[string]interface{})
	response["server"] = s.cfg.Name
	response["version"] = s.cfg.Version

	pl.SendResponse(nil, "info", message, response)
}

func (s *Saloon) loginPlayer(pl *player.Player) {
	var welcomePlayer *player.Player
	// if player has previous login and remove it
	players := s.players.GetAllValues()
	for _, p := range players {
		if p != pl && p.User != nil && p.User.ID == pl.User.ID {
			p.User = nil
			welcomePlayer = p
			break
		}
	}

	if welcomePlayer != nil {
		s.greetingMesssage(welcomePlayer, "disconnected (using another session)")
	}
}

func (s *Saloon) sendResponseError(pl *player.Player, req *request.WSRequest, message string, err error) {
	var data map[string]interface{}
	if err != nil && s.cfg.Debug {
		data = make(map[string]interface{})
		data["error"] = err.Error()
	}

	pl.SendResponse(req, "error", message, data)
}

func (s *Saloon) sendResponseSuccess(pl *player.Player, req *request.WSRequest, message string, data map[string]interface{}) {
	pl.SendResponse(req, "success", message, data)
}
