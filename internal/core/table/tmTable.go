package table

import (
	"time"

	"github.com/dev-rodrigobaliza/carteado/domain/core/game"
	tbl "github.com/dev-rodrigobaliza/carteado/domain/core/table"
	"github.com/dev-rodrigobaliza/carteado/domain/request"
	"github.com/dev-rodrigobaliza/carteado/errors"
)

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
	// TODO (@dev-rodrigobaliza) should locate other tables owned by this player and remove them???
	player.tableID = table.GetID()
	// make response
	response := t.getTableStatusResponse(table)
	// send response
	t.sendResponseSuccess(player, message, "table created", response)
	// debug log
	t.debug("=== table create %v", response)
}

func (g *TableManager) resourceTableDelete(player *Player, message *request.WSRequest) {
	// input validation
	tableID := g.getTableID(message)
	if tableID == "" {
		g.sendResponseError(player, message, "table id invalid, nil", nil)
		return
	}
	if tableID != player.tableID {
		g.sendResponseError(player, message, "table id invalid", nil)
		return
	}
	// database validation
	// TODO (@dev-rodrigobaliza) database validation ???
	// table validation
	table, err := g.getTable(tableID)
	if err != nil || table == nil {
		g.sendResponseError(player, message, "table id invalid", nil)
		return
	}
	// onwnership validation
	if !player.user.IsAdmin || table.GetStatus().Owner != player.uuid || table.GetStatus().Owner != "" {
		g.sendResponseError(player, message, "table id not owned by player", nil)
		return
	}
	// remove table
	g.delTable(table)
	// make response
	response := make(map[string]interface{})
	response["table_id"] = player.tableID
	// send response
	g.sendResponseSuccess(player, message, "table removed", response)
	// debug log
	g.debug("=== table remove %v", response)
}

func (g *TableManager) resourceTableRemovePlayer(player *Player, message *request.WSRequest) {
	// input validation
	tableID := g.getTableID(message)
	if tableID == "" {
		g.sendResponseError(player, message, "table id invalid", nil)
		return
	}
	if tableID != player.tableID {
		g.sendResponseError(player, message, "table id invalid", nil)
		return
	}
	// database validation
	// TODO (@dev-rodrigobaliza) database validation ???
	// table validation
	table, err := g.getTable(tableID)
	if err != nil || table == nil {
		g.sendResponseError(player, message, "table id invalid", nil)
		return
	}
	// leave table
	err = table.DelPlayer(player.uuid)
	if err != nil {
		if err != errors.ErrMinPlayers {
			g.sendResponseError(player, message, "leave table failed", err)
			return
		}
		// min players reached
		go g.resourceTableRemoveMinReached(table, true, err)
	}
	// send response
	g.sendResponseSuccess(player, message, "leave table", nil)
	// debug log
	g.debug("=== table leave %s", tableID)
}

func (g *TableManager) resourceTableRemoveMinReached(table *Table, force bool, err error) {
	players := g.players.GetAllValues()
	// remove all players from this table
	for _, player := range players {
		if player.tableID == table.GetID() {
			player.tableID = ""
			// send the bad news
			g.sendResponseError(player, nil, "table removed", err)
		}
	}
	// wait some time
	time.Sleep(time.Second * 10)
	// set table state
	table.Stop(force)
	// remove the table
	err = g.delTable(table)
	if err != nil {
		g.debug("error while deleting table (min reached): %s", err.Error())
	}
}

func (g *TableManager) resourceTableStartGame(player *Player, message *request.WSRequest) {
	// input validation
	tableID := g.getTableID(message)
	if tableID == "" {
		g.sendResponseError(player, message, "table id invalid, nil", nil)
		return
	}
	// database validation
	// TODO (@dev-rodrigobaliza) database validation ???
	// table validation
	table, err := g.getTable(tableID)
	if err != nil || table == nil {
		g.sendResponseError(player, message, "table id invalid", nil)
		return
	}
	// onwnership validation
	if !player.user.IsAdmin || table.GetStatus().Owner != player.uuid || table.GetStatus().Owner != "" {
		g.sendResponseError(player, message, "table id not owned by player", nil)
		return
	}
	// start table
	err = table.Play()
	if err != nil {
		g.sendResponseError(player, message, "table not started", err)
		return
	}
	// make response
	response := make(map[string]interface{})
	response["game_id"] = player.tableID
	// send response
	g.sendResponseSuccess(player, message, "table started", response)
	// debug log
	g.debug("=== table start %v", response)
}

func (g *TableManager) resourceTableStatus(player *Player, message *request.WSRequest) {
	// input validation
	tableID := g.getTableID(message)
	if tableID == "" {
		g.sendResponseError(player, message, "table id invalid", nil)
		return
	}
	// database validation
	// TODO (@dev-rodrigobaliza) database validation ???
	// table validation
	table, err := g.getTable(tableID)
	if err != nil || table == nil {
		g.sendResponseError(player, message, "table id invalid", nil)
		return
	}
	// make reponse
	response := g.getTableStatusResponse(table)
	// send response
	g.sendResponseSuccess(player, message, "table status", response)
	// debug log
	g.debug("=== table status %v", response)
}

func (g *TableManager) serviceTable(player *Player, message *request.WSRequest) {
	switch message.Resource {
	case "create":
		g.resourceTableCreate(player, message)

	case "enter":
		g.resourceTableAddPlayer(player, message)

	case "leave":
		g.resourceTableRemovePlayer(player, message)

	case "remove":
		g.resourceTableDelete(player, message)

	case "start":
		g.resourceTableStartGame(player, message)

	case "status":
		g.resourceTableStatus(player, message)

	default:
		g.sendResponseError(player, message, "table resource not found", nil)
	}
}
