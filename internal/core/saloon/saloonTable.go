package saloon

import (
	"time"

	"github.com/dev-rodrigobaliza/carteado/consts"
	"github.com/dev-rodrigobaliza/carteado/domain/core/game"
	"github.com/dev-rodrigobaliza/carteado/domain/request"
	"github.com/dev-rodrigobaliza/carteado/errors"
	pl "github.com/dev-rodrigobaliza/carteado/internal/core/player"
	"github.com/dev-rodrigobaliza/carteado/internal/core/table"
)

func (s *Saloon) resourceTableAddPlayer(player *pl.Player, message *request.WSRequest) {
	// input validation
	tableID := s.getTableID(message)
	if tableID == "" {
		s.sendResponseError(player, message, "table id invalid", nil)
		return
	}
	secret, ok := message.Data["secret"].(string)
	if !ok {
		secret = ""
	}
	// database validation
	// TODO (@dev-rodrigobaliza) database validation ???
	// table validation
	table, err := s.getTable(tableID)
	if err != nil || table == nil {
		s.sendResponseError(player, message, "table id invalid", nil)
		return
	}
	// enter table (sit maybe?)
	err = table.AddPlayer(player, secret)
	if err != nil {
		s.sendResponseError(player, message, "enter table failed", err)
		return
	}
	// TODO (@dev-rodrigobaliza) check if player was in another table and remove it from there
	player.TableID = table.GetID()
	// make reponse
	response := s.getTableStatusResponse(table)
	// send response
	s.sendResponseSuccess(player, message, "enter table", response)
	// debug log
	s.debug("=== table enter %v", response)
}

func (s *Saloon) resourceTableCreate(player *pl.Player, message *request.WSRequest) {
	// input validation
	strGameMode, ok := message.Data["game_mode"].(string)
	if !ok {
		s.sendResponseError(player, message, "game mode invalid", nil)
		return
	}
	minPlayers, ok := message.Data["min_players"].(float64)
	if !ok {
		s.sendResponseError(player, message, "min players invalid", nil)
		return
	}
	maxPlayers, ok := message.Data["max_players"].(float64)
	if !ok || int(maxPlayers) > consts.TABLE_MAX_PLAYERS {
		s.sendResponseError(player, message, "max players invalid", nil)
		return
	}
	allowBots, ok := message.Data["allow_bots"].(bool)
	if !ok {
		s.sendResponseError(player, message, "allow bots invalid", nil)
		return
	}
	// TODO (@dev-rodrigobaliza) anyone can create private table ???
	secret, ok := message.Data["secret"].(string)
	if !ok {
		secret = ""
	}
	// table validation
	gameMode := game.StringToMode(strGameMode)
	if gameMode == game.ModeUnknown {
		s.sendResponseError(player, message, "game mode invalid", nil)
		return
	}
	if maxPlayers <= 0 {
		s.sendResponseError(player, message, "min players invalid", nil)
		return
	}
	if maxPlayers <= 0 {
		s.sendResponseError(player, message, "max players invalid", nil)
		return
	}
	// database validation
	// TODO (@dev-rodrigobaliza) database validation ???
	// create new table
	table, err := table.NewTable(player, secret, int(minPlayers), int(maxPlayers), allowBots, gameMode)
	if err != nil {
		s.sendResponseError(player, message, "max players invalid", err)
		return
	}
	// add table to table list
	s.addTable(table)
	// if player is registered in other table, remove from there and register here
	if player.TableID != "" {
		s.resourceTableRemovePlayer(player, message, player.TableID)
	}
	player.TableID = table.GetID()
	// make response
	response := s.getTableStatusResponse(table)
	// send response
	s.sendResponseSuccess(player, message, "table created", response)
	// debug log
	s.debug("=== table create %v", response)
}

func (s *Saloon) resourceTableDelete(player *pl.Player, message *request.WSRequest) {
	// input validation
	tableID := s.getTableID(message)
	if tableID == "" {
		s.sendResponseError(player, message, "table id invalid, nil", nil)
		return
	}
	if tableID != player.TableID {
		s.sendResponseError(player, message, "table id invalid", nil)
		return
	}
	// database validation
	// TODO (@dev-rodrigobaliza) database validation ???
	// table validation
	table, err := s.getTable(tableID)
	if err != nil || table == nil {
		s.sendResponseError(player, message, "table id invalid", nil)
		return
	}
	// onwnership validation
	if !player.User.IsAdmin || table.GetStatus().Owner != player.UUID || table.GetStatus().Owner != "" {
		s.sendResponseError(player, message, "table id not owned by player", nil)
		return
	}
	// remove table
	s.delTable(table)
	// make response
	response := make(map[string]interface{})
	response["table_id"] = player.TableID
	// send response
	s.sendResponseSuccess(player, message, "table removed", response)
	// debug log
	s.debug("=== table remove %v", response)
}

func (s *Saloon) resourceTableRemovePlayer(player *pl.Player, message *request.WSRequest, tableID string) {
	// empty tableID means request is external (client)
	// otherwise request is internal, (server, probably player leaving one table to create another)
	internal := (tableID != "")

	if !internal {
		// input validation
		tableID = s.getTableID(message)
		if tableID == "" {
			s.sendResponseError(player, message, "table id invalid", nil)
			return
		}
		if tableID != player.TableID {
			s.sendResponseError(player, message, "table id invalid", nil)
			return
		}
	}
	// database validation
	// TODO (@dev-rodrigobaliza) database validation ???
	// table validation
	table, err := s.getTable(tableID)
	if err != nil || table == nil {
		if !internal {
			s.sendResponseError(player, message, "table id invalid", nil)
		}
		return
	}
	// leave table
	err = table.DelPlayer(player.UUID)
	if err != nil {
		if err != errors.ErrMinPlayers && err != errors.ErrEmptyTable && !internal {
			s.sendResponseError(player, message, "leave table failed", err)
			return
		}
		// table empty or min players reached
		s.resourceTableRemoveForced(table, true, err)
	}
	// send response
	data := make(map[string]interface{})
	data["table_id"] = tableID
	s.sendResponseSuccess(player, message, "leave table", data)
	// debug log
	s.debug("=== table leave %s [internal: %t]", tableID, internal)
}

func (s *Saloon) resourceTableRemoveForced(table *table.Table, force bool, err error) {
	if err == errors.ErrMinPlayers {
		players := s.players.GetAllValues()
		// remove all players from this table
		for _, player := range players {
			if player.TableID == table.GetID() {
				player.TableID = ""
				// send the bad news
				s.sendResponseError(player, nil, "table removed", err)
			}
		}
		// wait some time
		time.Sleep(time.Second * 3)
	}
	// set table state
	table.Stop(force)
	// remove the table
	err = s.delTable(table)
	if err != nil {
		s.debug("error while deleting table (forced remove): %s", err.Error())
	}
}

func (s *Saloon) resourceTableStartGame(player *pl.Player, message *request.WSRequest) {
	// input validation
	tableID := s.getTableID(message)
	if tableID == "" {
		s.sendResponseError(player, message, "table id invalid, nil", nil)
		return
	}
	// database validation
	// TODO (@dev-rodrigobaliza) database validation ???
	// table validation
	table, err := s.getTable(tableID)
	if err != nil || table == nil {
		s.sendResponseError(player, message, "table id invalid", nil)
		return
	}
	// onwnership validation
	if !player.User.IsAdmin || table.GetStatus().Owner != player.UUID || table.GetStatus().Owner != "" {
		s.sendResponseError(player, message, "table id not owned by player", nil)
		return
	}
	// start table
	err = table.Play()
	if err != nil {
		s.sendResponseError(player, message, "table not started", err)
		return
	}
	// make response
	response := make(map[string]interface{})
	response["game_id"] = player.TableID
	// send response
	s.sendResponseSuccess(player, message, "table started", response)
	// debug log
	s.debug("=== table start %v", response)
}

func (s *Saloon) resourceTableStatus(player *pl.Player, message *request.WSRequest) {
	// input validation
	tableID := s.getTableID(message)
	if tableID == "" {
		s.sendResponseError(player, message, "table id invalid", nil)
		return
	}
	// database validation
	// TODO (@dev-rodrigobaliza) database validation ???
	// table validation
	table, err := s.getTable(tableID)
	if err != nil || table == nil {
		s.sendResponseError(player, message, "table id invalid", nil)
		return
	}
	// make reponse
	response := s.getTableStatusResponse(table)
	// send response
	s.sendResponseSuccess(player, message, "table status", response)
	// debug log
	s.debug("=== table status %v", response)
}

func (s *Saloon) serviceTable(player *pl.Player, message *request.WSRequest) {
	switch message.Resource {
	case "create":
		s.resourceTableCreate(player, message)

	case "enter":
		s.resourceTableAddPlayer(player, message)

	case "leave":
		s.resourceTableRemovePlayer(player, message, "")

	case "remove":
		s.resourceTableDelete(player, message)

	case "start":
		s.resourceTableStartGame(player, message)

	case "status":
		s.resourceTableStatus(player, message)

	default:
		s.sendResponseError(player, message, "table resource not found", nil)
	}
}
