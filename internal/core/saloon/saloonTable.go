package saloon

import (
	"fmt"
	"time"

	"github.com/dev-rodrigobaliza/carteado/consts"
	"github.com/dev-rodrigobaliza/carteado/domain/core/game"
	"github.com/dev-rodrigobaliza/carteado/domain/request"
	"github.com/dev-rodrigobaliza/carteado/errors"
	"github.com/dev-rodrigobaliza/carteado/internal/core/player"
	"github.com/dev-rodrigobaliza/carteado/internal/core/table"
)

func (s *Saloon) resourceTableCreate(pl *player.Player, req *request.WSRequest) {
	// input validation
	strGameMode, ok := req.Data["game_mode"].(string)
	if !ok {
		s.sendResponseError(pl, req, "game mode invalid", nil)
		return
	}
	// TODO (@dev-rodrigobaliza) use default values from game mode
	minPlayers, ok := req.Data["min_players"].(float64)
	if !ok {
		s.sendResponseError(pl, req, "min players invalid", nil)
		return
	}
	maxPlayers, ok := req.Data["max_players"].(float64)
	if !ok || int(maxPlayers) > consts.TABLE_MAX_PLAYERS {
		s.sendResponseError(pl, req, "max players invalid", nil)
		return
	}
	allowBots, ok := req.Data["allow_bots"].(bool)
	if !ok {
		s.sendResponseError(pl, req, "allow bots invalid", nil)
		return
	}
	// TODO (@dev-rodrigobaliza) anyone can create private table ???
	secret, ok := req.Data["secret"].(string)
	if !ok {
		secret = ""
	}
	// table validation
	gameMode := game.StringToMode(strGameMode)
	if gameMode == game.ModeUnknown {
		s.sendResponseError(pl, req, "game mode invalid", nil)
		return
	}
	if maxPlayers <= 0 {
		s.sendResponseError(pl, req, "min players invalid", nil)
		return
	}
	if maxPlayers <= 0 {
		s.sendResponseError(pl, req, "max players invalid", nil)
		return
	}
	// database validation
	// TODO (@dev-rodrigobaliza) database validation ???
	// if player is registered in other table, remove from there and register here
	if pl.TableID != "" {
		s.resourceTableLeave(pl, req, pl.TableID)
	}
	// create new table
	table, err := table.New(pl, secret, int(minPlayers), int(maxPlayers), allowBots, gameMode)
	if err != nil {
		s.sendResponseError(pl, req, "max players invalid", err)
		return
	}
	// add table to table list
	s.addTable(table)
	// make response
	response := s.getTableStatus(table, pl.User.IsAdmin)
	// send response
	s.sendResponseSuccess(pl, req, "table create", response)
	// debug log
	s.debug("=== table create %v", response)
}

func (s *Saloon) resourceTableDelete(pl *player.Player, req *request.WSRequest) {
	// input validation
	tableID := s.getTableID(req)
	if tableID == "" {
		s.sendResponseError(pl, req, "table id invalid, nil", nil)
		return
	}
	if tableID != pl.TableID {
		s.sendResponseError(pl, req, "table id invalid", nil)
		return
	}
	// database validation
	// TODO (@dev-rodrigobaliza) database validation ???
	// table validation
	table, err := s.getTable(tableID)
	if err != nil || table == nil {
		s.sendResponseError(pl, req, "table id invalid", nil)
		return
	}
	// onwnership validation
	if !pl.User.IsAdmin || table.GetOwner() != pl.UUID || table.GetOwner() != "" {
		s.sendResponseError(pl, req, "table not owned by player", nil)
		return
	}
	// remove table
	s.delTable(table)
	// make response
	response := make(map[string]interface{})
	response["table_id"] = pl.TableID
	// send response
	s.sendResponseSuccess(pl, req, "table remove", response)
	// debug log
	s.debug("=== table remove %v", response)
}

func (s *Saloon) resourceTableEnter(pl *player.Player, req *request.WSRequest) {
	// input validation
	tableID := s.getTableID(req)
	if tableID == "" {
		s.sendResponseError(pl, req, "table id invalid", nil)
		return
	}
	secret, ok := req.Data["secret"].(string)
	if !ok {
		secret = ""
	}
	// database validation
	// TODO (@dev-rodrigobaliza) database validation ???
	// table validation
	table, err := s.getTable(tableID)
	if err != nil || table == nil {
		s.sendResponseError(pl, req, "table id invalid", nil)
		return
	}
	// if player is registered in other table, remove from there and register here
	if pl.TableID != "" {
		s.resourceTableLeave(pl, req, pl.TableID)
	}
	// enter table (sit maybe?)
	err = table.AddPlayer(pl, secret)
	if err != nil {
		s.sendResponseError(pl, req, "enter table failed", err)
		return
	}
	// make reponse
	response := s.getTableStatus(table, pl.User.IsAdmin)
	// send response
	s.sendResponseSuccess(pl, req, "enter table", response)
	// debug log
	s.debug("=== table enter %v", response)
}

func (s *Saloon) resourceTableGame(pl *player.Player, req *request.WSRequest) {
	// input validation
	tableID := s.getTableID(req)
	if tableID == "" {
		s.sendResponseError(pl, req, "table id invalid, nil", nil)
		return
	}
	if tableID != pl.TableID {
		s.sendResponseError(pl, req, "table id invalid", nil)
		return
	}
	action := s.getAction(req)
	if tableID == "" {
		s.sendResponseError(pl, req, "action invalid, nil", nil)
		return
	}
	// database validation
	// TODO (@dev-rodrigobaliza) database validation ???
	// table validation
	table, err := s.getTable(tableID)
	if err != nil || table == nil {
		s.sendResponseError(pl, req, "table id invalid", nil)
		return
	}
	// create response
	var response map[string]interface{}
	switch action {
	case "continue", "discontinue":
		pl.Action = action

	case "status":
		response = s.getTableGameStatus(table)

	case "start":
		response = s.actionTableGameStart(pl, req, table)

	case "stop":
		response = s.actionTableGameStop(pl, req, table)

	default:
		s.sendResponseError(pl, req, "action invalid", nil)
		return
	}
	// send response
	if response != nil {
		s.sendResponseSuccess(pl, req, fmt.Sprintf("%s table game", action), response)
		// debug log
		s.debug("=== table group %s %v", action, response)
	}
}

func (s *Saloon) resourceTableGroup(pl *player.Player, req *request.WSRequest) {
	// input validation
	tableID := s.getTableID(req)
	if tableID == "" {
		s.sendResponseError(pl, req, "table id invalid", nil)
		return
	}
	action := s.getAction(req)
	if action == "" {
		s.sendResponseError(pl, req, "action invalid", nil)
		return
	}
	var groupID int
	if action != "bot" {
		groupID = s.getGroupID(req)
		if groupID == 0 {
			s.sendResponseError(pl, req, "group id invalid", nil)
			return
		}
	}
	var quantity int
	if action == "bot" {
		quantity = s.getQuantity(req)
		if quantity == 0 {
			s.sendResponseError(pl, req, "bot quantity invalid", nil)
			return
		}
	}
	// database validation
	// TODO (@dev-rodrigobaliza) database validation ???
	// table validation
	table, err := s.getTable(tableID)
	if err != nil || table == nil {
		s.sendResponseError(pl, req, "table id invalid", nil)
		return
	}
	if action != "bot" && !table.HasGroup(groupID) {
		s.sendResponseError(pl, req, "group id invalid", nil)
		return
	}
	// create response
	var response map[string]interface{}
	switch action {
	case "enter":
		response = s.actionTableGroupEnter(pl, req, table, groupID)

	case "leave":
		response = s.actionTableGroupLeave(pl, req, table, groupID)

	case "bot":
		response = s.actionTableGroupBot(pl, req, table, quantity)

	case "status":
		response = s.getTableGroupStatus(table, groupID, pl.User.IsAdmin)

	default:
		s.sendResponseError(pl, req, "action invalid", nil)
		return
	}
	// send response
	if response != nil {
		s.sendResponseSuccess(pl, req, fmt.Sprintf("%s table group", action), response)
		// debug log
		s.debug("=== table group %s %v", action, response)
	}
}

func (s *Saloon) resourceTableLeave(pl *player.Player, req *request.WSRequest, tableID string) {
	// empty tableID means request is external (client)
	// otherwise request is internal, (server, probably player leaving one table to create another)
	internal := (tableID != "")

	if !internal {
		// input validation
		tableID = s.getTableID(req)
		if tableID == "" {
			s.sendResponseError(pl, req, "table id invalid", nil)
			return
		}
		if tableID != pl.TableID {
			s.sendResponseError(pl, req, "table id invalid", nil)
			return
		}
	}
	// database validation
	// TODO (@dev-rodrigobaliza) database validation ???
	// table validation
	table, err := s.getTable(tableID)
	if err != nil || table == nil {
		if !internal {
			s.sendResponseError(pl, req, "table id invalid", nil)
		}
		return
	}
	// leave table
	err = table.DelPlayer(pl.UUID)
	if err != nil {
		if err != errors.ErrMinPlayers && err != errors.ErrEmptyTable && !internal {
			s.sendResponseError(pl, req, "leave table failed", err)
			return
		}
		// table empty or min players reached
		s.resourceTableRemoveForced(table, true, err)
	}
	// send response
	data := make(map[string]interface{})
	data["table_id"] = tableID
	s.sendResponseSuccess(pl, req, "leave table", data)
	// debug log
	s.debug("=== table leave %s [internal: %t]", tableID, internal)
}

func (s *Saloon) resourceTableRemoveForced(tb *table.Table, force bool, err error) {
	if err == errors.ErrMinPlayers {
		players := s.players.GetAllValues()
		// remove all players from this table
		for _, player := range players {
			if player.TableID == tb.GetID() {
				player.TableID = ""
				// send the bad news
				s.sendResponseError(player, nil, "table removed", err)
			}
		}
		// wait some time
		time.Sleep(time.Second * 3)
	}
	// set table state
	tb.Stop(force)
	// remove the table
	err = s.delTable(tb)
	if err != nil {
		s.debug("error while deleting table (forced remove): %s", err.Error())
	}
}

func (s *Saloon) resourceTableStatus(pl *player.Player, req *request.WSRequest) {
	// input validation
	tableID := s.getTableID(req)
	if tableID == "" {
		s.sendResponseError(pl, req, "table id invalid", nil)
		return
	}
	// database validation
	// TODO (@dev-rodrigobaliza) database validation ???
	// table validation
	table, err := s.getTable(tableID)
	if err != nil || table == nil {
		s.sendResponseError(pl, req, "table id invalid", nil)
		return
	}
	// send response
	response := s.getTableStatus(table, pl.User.IsAdmin)
	s.sendResponseSuccess(pl, req, "status table", response)
	// debug log
	s.debug("=== table status %v", response)
}

func (s *Saloon) serviceTable(pl *player.Player, req *request.WSRequest) {
	switch req.Resource {
	case "create":
		s.resourceTableCreate(pl, req)

	case "enter":
		s.resourceTableEnter(pl, req)

	case "game":
		s.resourceTableGame(pl, req)

	case "group":
		s.resourceTableGroup(pl, req)

	case "leave":
		s.resourceTableLeave(pl, req, "")

	case "remove":
		s.resourceTableDelete(pl, req)

	case "status":
		s.resourceTableStatus(pl, req)

	default:
		s.sendResponseError(pl, req, "table resource not found", nil)
	}
}
