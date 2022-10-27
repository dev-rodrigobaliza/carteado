package saloon

import (
	"log"

	tbl "github.com/dev-rodrigobaliza/carteado/domain/core/table"
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
	// get players
	players := make([]*response.Player, 0)
	for _, player := range s.players.GetAllValues() {
		if !authenticatedOnly || player.User != nil {
			players = append(players, player.ToResponse())
		}
	}
	// get tables
	tables := make([]*response.Table, 0)
	for _, table := range s.tables.GetAllValues() {
		tables = append(tables, table.ToResponse())
	}

	response := make(map[string]interface{})
	response["server"] = s.cfg.Name
	response["version"] = s.cfg.Version
	response["started_at"] = s.cfg.StartedAt
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

func (s *Saloon) getTableID(message *request.WSRequest) string {
	tableID, ok := message.Data["table_id"].(string)
	if !ok {
		return ""
	}

	return tableID
}

func (s *Saloon) getTableStatusResponse(table *table.Table) map[string]interface{} {
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
	if s.cfg.Debug {
		response["players"] = table.GetPlayers()
	}

	return response
}

func (s *Saloon) greetingMesssage(player *pl.Player, message string) {
	response := make(map[string]interface{})
	response["server"] = s.cfg.Name
	response["version"] = s.cfg.Version

	s.sendResponseSuccess(player, nil, message, response)
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

func (s *Saloon) sendResponse(player *pl.Player, request *request.WSRequest, status, message string, data map[string]interface{}) {
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

func (s *Saloon) sendResponseError(player *pl.Player, request *request.WSRequest, message string, err error) {
	var data map[string]interface{}
	if err != nil && s.cfg.Debug {
		data = make(map[string]interface{})
		data["error"] = err.Error()
	}

	s.sendResponse(player, request, "error", message, data)
}

func (s *Saloon) sendResponseSuccess(player *pl.Player, request *request.WSRequest, message string, data map[string]interface{}) {
	s.sendResponse(player, request, "success", message, data)
}
