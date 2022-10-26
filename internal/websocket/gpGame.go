package websocket

import (
	"github.com/dev-rodrigobaliza/carteado/domain/request"
	"github.com/dev-rodrigobaliza/carteado/internal/core"
)

func (g *GameProcessor) resourceGameCreate(player *Player, message *request.WSRequest) {
	// input validation
	strGameType, ok := message.Data["game_type"].(string)
	if !ok {
		g.sendResponseError(player, message, "game type invalid")
	}
	minPlayers, ok := message.Data["min_players"].(float64)
	if !ok {
		g.sendResponseError(player, message, "min players invalid")
	}
	maxPlayers, ok := message.Data["max_players"].(float64)
	if !ok {
		g.sendResponseError(player, message, "max players invalid")
	}
	allowBots, ok := message.Data["allow_bots"].(bool)
	if !ok {
		g.sendResponseError(player, message, "allow bots invalid")
	}
	// game validation
	gameType := core.StringToGametype(strGameType)
	if gameType == core.GameTypeUnknown {
		g.sendResponseError(player, message, "game type invalid")
		return
	}
	if maxPlayers <= 0 {
		g.sendResponseError(player, message, "min players invalid")
	}
	if maxPlayers <= 0 {
		g.sendResponseError(player, message, "max players invalid")
	}
	// database validation
	// TODO (@dev-rodrigobaliza) database validation ???
	// create new game
	var game core.IGame
	switch gameType {
	case core.GameTypeBlackJack:
		game = core.NewBlackJack(player.id, int(minPlayers), int(maxPlayers), allowBots)
	}
	// TODO (@dev-rodrigobaliza) should locate other games owned by this player and remove them???
	player.gameID = game.GetStatus().ID
	// make response
	response := make(map[string]interface{})
	response["game_id"] = player.gameID

	g.addGame(game)
	g.sendResponseSuccess(player, message, "game created", response)
}

func (g *GameProcessor) resourceGameRemove(player *Player, message *request.WSRequest) {
	// input validation
	gameID := g.getGameID(message)
	if gameID == "" {
		g.sendResponseError(player, message, "game id invalid")
		return
	}
	// database validation
	// TODO (@dev-rodrigobaliza) database validation ???
	// game validation
	game := g.getGame(gameID)
	if game == nil {
		g.sendResponseError(player, message, "game id invalid")
		return
	}
	// onwnership validation
	if game.GetStatus().Owner != player.id {
		g.sendResponseError(player, message, "game id not owned by player")
		return
	}
	// remove game
	g.delGame(game)
	// make response
	response := make(map[string]interface{})
	response["game_id"] = player.gameID

	g.sendResponseSuccess(player, message, "game removed", response)
}

func (g *GameProcessor) resourceGameStatus(player *Player, message *request.WSRequest) {
	// input validation
	gameID := g.getGameID(message)
	if gameID == "" {
		g.sendResponseError(player, message, "game id invalid")
		return
	}
	// database validation
	// TODO (@dev-rodrigobaliza) database validation ???
	// game validation
	game := g.getGame(gameID)
	if game == nil {
		g.sendResponseError(player, message, "game id invalid")
		return
	}
	// get game status
	gameStatus := game.GetStatus()
	gameState := gameStatus.GameState.String()
	table := game.GetTable()
	// make reponse
	response := make(map[string]interface{})
	response["game_id"] = gameStatus.ID
	response["game_type"] = gameStatus.GameType.String()
	response["game_state"] = gameState
	if gameState != "starting" {
		response["round"] = gameStatus.GameRound
	}
	response["min_players"] = table.GetMinPlayers()
	response["max_players"] = table.GetMaxPlayers()
	response["allow_bots"] = table.GetAllowBots()
	response["player_count"] = gameStatus.PlayerCount
	if g.cfg.Debug {
		response["players"] = table.GetPlayers()
	}
	if gameState == "finished" {
		response["winners"] = gameStatus.Winners
	}

	g.sendResponseSuccess(player, message, "game status", response)
}

func (g *GameProcessor) serviceGame(player *Player, message *request.WSRequest) {
	switch message.Resource {
	case "create":
		g.resourceGameCreate(player, message)

	case "remove":
		g.resourceGameRemove(player, message)

	case "status":
		g.resourceGameStatus(player, message)

	default:
		g.sendResponseError(player, message, "game resource not found")
	}
}
