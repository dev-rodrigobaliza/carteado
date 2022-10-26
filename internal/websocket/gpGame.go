package websocket

import (
	"time"

	"github.com/dev-rodrigobaliza/carteado/domain/request"
	"github.com/dev-rodrigobaliza/carteado/errors"
	"github.com/dev-rodrigobaliza/carteado/internal/core"
)

func (g *GameProcessor) getGameStatusResponse(game core.IGame) map[string]interface{} {
	// get game status
	gameStatus := game.GetStatus()
	gameState := gameStatus.GameState.String()
	table := game.GetTable()
	// make reponse
	response := make(map[string]interface{})
	response["game_id"] = gameStatus.ID
	response["game_type"] = gameStatus.GameType.String()
	response["game_state"] = gameState
	response["game_momentum"] = gameStatus.GameMomentum
	if gameState != "starting" {
		response["round"] = gameStatus.GameRound
	}
	response["min_players"] = table.GetMinPlayers()
	response["max_players"] = table.GetMaxPlayers()
	response["allow_bots"] = table.GetAllowBots()
	response["private"] = table.IsPrivate()
	response["player_count"] = gameStatus.PlayerCount
	if g.cfg.Debug {
		response["players"] = table.GetPlayers()
	}
	if gameState == "finished" {
		response["winners"] = gameStatus.Winners
	}

	return response
}

func (g *GameProcessor) resourceGameCreate(player *Player, message *request.WSRequest) {
	// input validation
	strGameType, ok := message.Data["game_type"].(string)
	if !ok {
		g.sendResponseError(player, message, "game type invalid", nil)
	}
	minPlayers, ok := message.Data["min_players"].(float64)
	if !ok {
		g.sendResponseError(player, message, "min players invalid", nil)
	}
	maxPlayers, ok := message.Data["max_players"].(float64)
	if !ok {
		g.sendResponseError(player, message, "max players invalid", nil)
	}
	allowBots, ok := message.Data["allow_bots"].(bool)
	if !ok {
		g.sendResponseError(player, message, "allow bots invalid", nil)
	}
	// TODO (@dev-rodrigobaliza) anyone can create private game ???
	secret, ok := message.Data["secret"].(string)
	if !ok {
		secret = ""
	}
	// game validation
	gameType := core.StringToGametype(strGameType)
	if gameType == core.GameTypeUnknown {
		g.sendResponseError(player, message, "game type invalid", nil)
		return
	}
	if maxPlayers <= 0 {
		g.sendResponseError(player, message, "min players invalid", nil)
	}
	if maxPlayers <= 0 {
		g.sendResponseError(player, message, "max players invalid", nil)
	}
	// database validation
	// TODO (@dev-rodrigobaliza) database validation ???
	// create new game
	var game core.IGame
	switch gameType {
	case core.GameTypeBlackJack:
		game = core.NewBlackJack(player.uuid, secret, int(minPlayers), int(maxPlayers), allowBots)
	}
	// add game to game list
	g.addGame(game)
	// TODO (@dev-rodrigobaliza) should locate other games owned by this player and remove them???
	player.gameID = game.GetID()
	// make response
	response := g.getGameStatusResponse(game)
	// send response
	g.sendResponseSuccess(player, message, "game created", response)
	// debug log
	g.debug("=== game create %v", response)
}

func (g *GameProcessor) resourceGameEnter(player *Player, message *request.WSRequest) {
	// input validation
	gameID := g.getGameID(message)
	if gameID == "" {
		g.sendResponseError(player, message, "game id invalid", nil)
		return
	}
	secret, ok := message.Data["secret"].(string)
	if !ok {
		secret = ""
	}
	// database validation
	// TODO (@dev-rodrigobaliza) database validation ???
	// game validation
	game, err := g.getGame(gameID)
	if err != nil || game == nil {
		g.sendResponseError(player, message, "game id invalid", nil)
		return
	}
	// enter game
	err = game.Enter(player.uuid, secret)
	if err != nil {
		g.sendResponseError(player, message, "enter game failed", err)
		return
	}
	// TODO (@dev-rodrigobaliza) check if player was in another game and remove it from there
	player.gameID = game.GetID()
	// make reponse
	response := g.getGameStatusResponse(game)
	// send response
	g.sendResponseSuccess(player, message, "enter game", response)
	// debug log
	g.debug("=== game enter %v", response)
}

func (g *GameProcessor) resourceGameFinishGame(game core.IGame, force bool, err error) {
	players := g.players.GetAllKeys()
	// remove all players from this game
	for _, player := range players {
		if player.gameID == game.GetID() {
			player.gameID = ""
			// send the bad news
			g.sendResponseError(player, nil, "game finished", err)
		}
	}
	// wait some time
	time.Sleep(time.Second * 10)
	// set game state
	game.Stop(force)
	// remove the game
	g.delGame(game)
}

func (g *GameProcessor) resourceGameLeave(player *Player, message *request.WSRequest) {
	// input validation
	gameID := g.getGameID(message)
	if gameID == "" {
		g.sendResponseError(player, message, "game id invalid", nil)
		return
	}
	if gameID != player.gameID {
		g.sendResponseError(player, message, "game id invalid", nil)
		return
	}
	// database validation
	// TODO (@dev-rodrigobaliza) database validation ???
	// game validation
	game, err := g.getGame(gameID)
	if err != nil || game == nil {
		g.sendResponseError(player, message, "game id invalid", nil)
		return
	}
	// leave game
	err = game.Leave(player.uuid)
	if err != nil {
		if err != errors.ErrMinPlayers {
			g.sendResponseError(player, message, "leave game failed", err)
			return
		}
		// min players reached
		go g.resourceGameFinishGame(game, true, err)
	}
	// send response
	g.sendResponseSuccess(player, message, "leave game", nil)
	// debug log
	g.debug("=== game leave %s", gameID)
}

func (g *GameProcessor) resourceGameRemove(player *Player, message *request.WSRequest) {
	// input validation
	gameID := g.getGameID(message)
	if gameID == "" {
		g.sendResponseError(player, message, "game id invalid, nil", nil)
		return
	}
	// database validation
	// TODO (@dev-rodrigobaliza) database validation ???
	// game validation
	game, err := g.getGame(gameID)
	if err != nil || game == nil {
		g.sendResponseError(player, message, "game id invalid", nil)
		return
	}
	// onwnership validation
	if !player.user.IsAdmin || game.GetOwner() != player.uuid || game.GetOwner() != "" {
		g.sendResponseError(player, message, "game id not owned by player", nil)
		return
	}
	// remove game
	g.delGame(game)
	// make response
	response := make(map[string]interface{})
	response["game_id"] = player.gameID
	// send response
	g.sendResponseSuccess(player, message, "game removed", response)
	// debug log
	g.debug("=== game remove %v", response)
}

func (g *GameProcessor) resourceGameStart(player *Player, message *request.WSRequest) {
	// input validation
	gameID := g.getGameID(message)
	if gameID == "" {
		g.sendResponseError(player, message, "game id invalid, nil", nil)
		return
	}
	// database validation
	// TODO (@dev-rodrigobaliza) database validation ???
	// game validation
	game, err := g.getGame(gameID)
	if err != nil || game == nil {
		g.sendResponseError(player, message, "game id invalid", nil)
		return
	}
	// onwnership validation
	if !player.user.IsAdmin || game.GetOwner() != player.uuid || game.GetOwner() != "" {
		g.sendResponseError(player, message, "game id not owned by player", nil)
		return
	}
	// start game
	err = game.Play()
	if err != nil {
		g.sendResponseError(player, message, "game not started", err)
		return
	}
	// make response
	response := make(map[string]interface{})
	response["game_id"] = player.gameID
	// send response
	g.sendResponseSuccess(player, message, "game started", response)
	// debug log
	g.debug("=== game start %v", response)
}

func (g *GameProcessor) resourceGameStatus(player *Player, message *request.WSRequest) {
	// input validation
	gameID := g.getGameID(message)
	if gameID == "" {
		g.sendResponseError(player, message, "game id invalid", nil)
		return
	}
	// database validation
	// TODO (@dev-rodrigobaliza) database validation ???
	// game validation
	game, err := g.getGame(gameID)
	if err != nil || game == nil {
		g.sendResponseError(player, message, "game id invalid", nil)
		return
	}
	// make reponse
	response := g.getGameStatusResponse(game)
	// send response
	g.sendResponseSuccess(player, message, "game status", response)
	// debug log
	g.debug("=== game status %v", response)
}

func (g *GameProcessor) serviceGame(player *Player, message *request.WSRequest) {
	switch message.Resource {
	case "create":
		g.resourceGameCreate(player, message)

	case "enter":
		g.resourceGameEnter(player, message)

	case "leave":
		g.resourceGameLeave(player, message)

	case "remove":
		g.resourceGameRemove(player, message)

	case "start":
		g.resourceGameStart(player, message)

	case "status":
		g.resourceGameStatus(player, message)

	default:
		g.sendResponseError(player, message, "game resource not found", nil)
	}
}
