package table

import (
	"time"

	"github.com/dev-rodrigobaliza/carteado/domain/core/game"
	"github.com/dev-rodrigobaliza/carteado/domain/request"
	"github.com/dev-rodrigobaliza/carteado/errors"
)

func (t *TableManager) resourceTableAddPlayer(player *Player, message *request.WSRequest) {
	// input validation
	tableID := t.getTableID(message)
	if tableID == "" {
		t.sendResponseError(player, message, "table id invalid", nil)
		return
	}
	secret, ok := message.Data["secret"].(string)
	if !ok {
		secret = ""
	}
	// database validation
	// TODO (@dev-rodrigobaliza) database validation ???
	// table validation
	table, err := t.getTable(tableID)
	if err != nil || table == nil {
		t.sendResponseError(player, message, "table id invalid", nil)
		return
	}
	// enter table (sit maybe?)
	err = table.AddPlayer(player, secret)
	if err != nil {
		t.sendResponseError(player, message, "enter table failed", err)
		return
	}
	// TODO (@dev-rodrigobaliza) check if player was in another table and remove it from there
	player.tableID = table.GetID()
	// make reponse
	response := t.getTableStatusResponse(table)
	// send response
	t.sendResponseSuccess(player, message, "enter table", response)
	// debug log
	t.debug("=== table enter %v", response)
}

func (t *TableManager) resourceTableCreate(player *Player, message *request.WSRequest) {
	// input validation
	strGameMode, ok := message.Data["game_mode"].(string)
	if !ok {
		t.sendResponseError(player, message, "game mode invalid", nil)
		return
	}
	minPlayers, ok := message.Data["min_players"].(float64)
	if !ok {
		t.sendResponseError(player, message, "min players invalid", nil)
		return
	}
	maxPlayers, ok := message.Data["max_players"].(float64)
	if !ok {
		t.sendResponseError(player, message, "max players invalid", nil)
		return
	}
	allowBots, ok := message.Data["allow_bots"].(bool)
	if !ok {
		t.sendResponseError(player, message, "allow bots invalid", nil)
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
		t.sendResponseError(player, message, "game mode invalid", nil)
		return
	}
	if maxPlayers <= 0 {
		t.sendResponseError(player, message, "min players invalid", nil)
		return
	}
	if maxPlayers <= 0 {
		t.sendResponseError(player, message, "max players invalid", nil)
		return
	}
	// database validation
	// TODO (@dev-rodrigobaliza) database validation ???
	// create new table
	table, err := NewTable(player, secret, int(minPlayers), int(maxPlayers), allowBots, gameMode)
	if err != nil {
		t.sendResponseError(player, message, "max players invalid", err)
		return
	}
	// add table to table list
	t.addTable(table)
	// if player is registered in other table, remove from there and register here
	if player.tableID != "" {
		t.resourceTableRemovePlayer(player, message, player.tableID)
	}
	player.tableID = table.GetID()
	// make response
	response := t.getTableStatusResponse(table)
	// send response
	t.sendResponseSuccess(player, message, "table created", response)
	// debug log
	t.debug("=== table create %v", response)
}

func (t *TableManager) resourceTableDelete(player *Player, message *request.WSRequest) {
	// input validation
	tableID := t.getTableID(message)
	if tableID == "" {
		t.sendResponseError(player, message, "table id invalid, nil", nil)
		return
	}
	if tableID != player.tableID {
		t.sendResponseError(player, message, "table id invalid", nil)
		return
	}
	// database validation
	// TODO (@dev-rodrigobaliza) database validation ???
	// table validation
	table, err := t.getTable(tableID)
	if err != nil || table == nil {
		t.sendResponseError(player, message, "table id invalid", nil)
		return
	}
	// onwnership validation
	if !player.user.IsAdmin || table.GetStatus().Owner != player.uuid || table.GetStatus().Owner != "" {
		t.sendResponseError(player, message, "table id not owned by player", nil)
		return
	}
	// remove table
	t.delTable(table)
	// make response
	response := make(map[string]interface{})
	response["table_id"] = player.tableID
	// send response
	t.sendResponseSuccess(player, message, "table removed", response)
	// debug log
	t.debug("=== table remove %v", response)
}

func (t *TableManager) resourceTableRemovePlayer(player *Player, message *request.WSRequest, tableID string) {
	// empty tableID means request is external (client)
	// otherwise request is internal, (server, probably player leaving one table to create another)
	internal := (tableID != "")

	if !internal {
		// input validation
		tableID = t.getTableID(message)
		if tableID == "" {
			t.sendResponseError(player, message, "table id invalid", nil)
			return
		}
		if tableID != player.tableID {
			t.sendResponseError(player, message, "table id invalid", nil)
			return
		}
	}
	// database validation
	// TODO (@dev-rodrigobaliza) database validation ???
	// table validation
	table, err := t.getTable(tableID)
	if err != nil || table == nil {
		if !internal {
			t.sendResponseError(player, message, "table id invalid", nil)
		}
		return
	}
	// leave table
	err = table.DelPlayer(player.uuid)
	if err != nil {
		if err != errors.ErrMinPlayers && err != errors.ErrEmptyTable && !internal {
			t.sendResponseError(player, message, "leave table failed", err)
			return
		}
		// table empty or min players reached
		t.resourceTableRemoveForced(table, true, err)
	}
	// send response
	data := make(map[string]interface{})
	data["table_id"] = tableID
	t.sendResponseSuccess(player, message, "leave table", data)
	// debug log
	t.debug("=== table leave %s [internal: %t]", tableID, internal)
}

func (t *TableManager) resourceTableRemoveForced(table *Table, force bool, err error) {
	if err == errors.ErrMinPlayers {
		players := t.players.GetAllValues()
		// remove all players from this table
		for _, player := range players {
			if player.tableID == table.GetID() {
				player.tableID = ""
				// send the bad news
				t.sendResponseError(player, nil, "table removed", err)
			}
		}
		// wait some time
		time.Sleep(time.Second * 3)
	}
	// set table state
	table.Stop(force)
	// remove the table
	err = t.delTable(table)
	if err != nil {
		t.debug("error while deleting table (forced remove): %s", err.Error())
	}
}

func (t *TableManager) resourceTableStartGame(player *Player, message *request.WSRequest) {
	// input validation
	tableID := t.getTableID(message)
	if tableID == "" {
		t.sendResponseError(player, message, "table id invalid, nil", nil)
		return
	}
	// database validation
	// TODO (@dev-rodrigobaliza) database validation ???
	// table validation
	table, err := t.getTable(tableID)
	if err != nil || table == nil {
		t.sendResponseError(player, message, "table id invalid", nil)
		return
	}
	// onwnership validation
	if !player.user.IsAdmin || table.GetStatus().Owner != player.uuid || table.GetStatus().Owner != "" {
		t.sendResponseError(player, message, "table id not owned by player", nil)
		return
	}
	// start table
	err = table.Play()
	if err != nil {
		t.sendResponseError(player, message, "table not started", err)
		return
	}
	// make response
	response := make(map[string]interface{})
	response["game_id"] = player.tableID
	// send response
	t.sendResponseSuccess(player, message, "table started", response)
	// debug log
	t.debug("=== table start %v", response)
}

func (t *TableManager) resourceTableStatus(player *Player, message *request.WSRequest) {
	// input validation
	tableID := t.getTableID(message)
	if tableID == "" {
		t.sendResponseError(player, message, "table id invalid", nil)
		return
	}
	// database validation
	// TODO (@dev-rodrigobaliza) database validation ???
	// table validation
	table, err := t.getTable(tableID)
	if err != nil || table == nil {
		t.sendResponseError(player, message, "table id invalid", nil)
		return
	}
	// make reponse
	response := t.getTableStatusResponse(table)
	// send response
	t.sendResponseSuccess(player, message, "table status", response)
	// debug log
	t.debug("=== table status %v", response)
}

func (t *TableManager) serviceTable(player *Player, message *request.WSRequest) {
	switch message.Resource {
	case "create":
		t.resourceTableCreate(player, message)

	case "enter":
		t.resourceTableAddPlayer(player, message)

	case "leave":
		t.resourceTableRemovePlayer(player, message, "")

	case "remove":
		t.resourceTableDelete(player, message)

	case "start":
		t.resourceTableStartGame(player, message)

	case "status":
		t.resourceTableStatus(player, message)

	default:
		t.sendResponseError(player, message, "table resource not found", nil)
	}
}
