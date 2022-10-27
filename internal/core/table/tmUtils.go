package table

import (
	"log"

	tbl "github.com/dev-rodrigobaliza/carteado/domain/core/table"
	"github.com/dev-rodrigobaliza/carteado/domain/request"
	"github.com/dev-rodrigobaliza/carteado/domain/response"
)

func (t *TableManager) debug(format string, v ...any) {
	if t.cfg.Debug {
		log.Printf(format, v...)
	}
}

func (t *TableManager) getServerStatusResponse(authenticatedOnly bool) map[string]interface{} {
	// get players
	players := make([]*response.Player, 0)
	allPlayers := t.players.GetAllValues()
	for _, player := range allPlayers {
		if !authenticatedOnly || player.user != nil {
			players = append(players, player.ToResponse())
		}
	}
	// get tables
	tables := make([]*response.Table, 0)
	allTables := t.tables.GetAllValues()
	for _, table := range allTables {
		tables = append(tables, table.ToResponse())
	}	

	response := make(map[string]interface{})
	response["server"] = t.cfg.Name
	response["version"] = t.cfg.Version
	response["started_at"] = t.cfg.StartedAt
	response["players_count"] = len(players)
	if len(players) > 0 {
		response["players"] = players
	}
	response["tables_count"] = len(tables)
	if len(tables) > 0 {
		response["tables"] = tables	
	}

	return response
}

func (t *TableManager) getTableID(message *request.WSRequest) string {
	tableID, ok := message.Data["table_id"].(string)
	if !ok {
		return ""
	}

	return tableID
}

func (t *TableManager) getTableStatusResponse(table *Table) map[string]interface{} {
	// get table status
	tableStatus := table.GetStatus()
	// make reponse
	response := make(map[string]interface{})
	response["table_id"] = tableStatus.ID
	response["table_owner"] = tableStatus.Owner
	if tableStatus.State == tbl.StateFinish {
		response["table_winners"] = tableStatus.Winners
	}
	response["min_players"] = table.GetMinPlayers()
	response["max_players"] = table.GetMaxPlayers()
	response["player_count"] = tableStatus.PlayerCount
	if tableStatus.State != tbl.StateStart {
		response["game_round"] = tableStatus.GameRound
	}
	response["allow_bots"] = table.GetAllowBots()
	response["private"] = table.IsPrivate()
	response["game_mode"] = tableStatus.GameMode.String()
	if tableStatus.State == tbl.StatePlay {
		response["game_state"] = tableStatus.GameState.String()
	}
	if t.cfg.Debug {
		response["players"] = table.GetPlayers()
	}

	return response
}

func (t *TableManager) sendResponse(player *Player, request *request.WSRequest, status, message string, data map[string]interface{}) {
	response := &response.WSResponse{
		Status:  status,
		Message: message,
	}
	if request != nil {
		response.RequestID = request.RequestID
	}
	if len(data) > 0 {
		response.Data = data
	}

	player.Send(response.ToBytes())
}

func (t *TableManager) sendResponseError(player *Player, request *request.WSRequest, message string, err error) {
	var data map[string]interface{}
	if err != nil && t.cfg.Debug {
		data = make(map[string]interface{})
		data["error"] = err.Error()
	}

	t.sendResponse(player, request, "error", message, data)
}

func (t *TableManager) sendResponseSuccess(player *Player, request *request.WSRequest, message string, data map[string]interface{}) {
	t.sendResponse(player, request, "success", message, data)
}
