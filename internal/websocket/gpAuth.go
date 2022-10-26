package websocket

import (
	"time"

	"github.com/dev-rodrigobaliza/carteado/domain/request"
	"github.com/dev-rodrigobaliza/carteado/utils"
)

func (g *GameProcessor) resourceAuthLogin(player *Player, message *request.WSRequest) {
	token, ok := message.Data["token"].(string)
	if !ok {
		g.sendResponseError(player, message, "token invalid")
	}
	// token validation
	id, err := Security.VerifyToken(token)
	if err != nil {
		g.sendResponseError(player, message, "token invalid")
		return
	}
	// database validation
	err = AuthService.VerifyToken(id, token)
	if err != nil {
		g.sendResponseError(player, message, "token invalid")
		return
	}

	if !player.auth {
		player.auth = true
		player.id = "pid-" + utils.NewUUID()
	}
	player.since = time.Now()

	response := make(map[string]interface{})
	response["player_id"] = player.id

	g.sendResponseSuccess(player, message, "authenticated", response)
}

func (g *GameProcessor) serviceAuth(player *Player, message *request.WSRequest) {
	switch message.Resource {
	case "login":
		g.resourceAuthLogin(player, message)

	default:
		g.sendResponseError(player, message, "auth resource not found")
	}
}
