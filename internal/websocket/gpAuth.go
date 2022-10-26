package websocket

import (
	"strconv"

	"github.com/dev-rodrigobaliza/carteado/domain/request"
)

func (g *GameProcessor) resourceAuthLogin(player *Player, message *request.WSRequest) {
	token, ok := message.Data["token"].(string)
	if !ok {
		g.sendResponseError(player, message, "token invalid", nil)
	}
	// token validation
	id, err := Security.VerifyToken(token)
	if err != nil {
		g.sendResponseError(player, message, "token invalid", err)
		return
	}
	// database validation
	err = g.appService.AuthService.VerifyToken(id, token)
	if err != nil {
		g.sendResponseError(player, message, "token invalid", err)
		return
	}
	// get user from database
	userID, _ := strconv.Atoi(id)
	u := &request.GetUser{
		ID: uint64(userID),
	}
	user, _, err := g.appService.UserService.Get(u)
	if err != nil || user == nil {
		g.sendResponseError(player, message, "user not found", err)
		return
	}
	// set player-user information
	firstLogin := player.Login(user)
	// make response
	response := make(map[string]interface{})
	response["first_login"] = firstLogin
	response["id"] = player.uuid
	response["name"] = player.user.Name
	response["email"] = player.user.Email
	response["is_admin"] = player.user.IsAdmin
	// send response
	g.sendResponseSuccess(player, message, "authenticated", response)
	// debug log
	g.debug("=== auth login %v", response)
}

func (g *GameProcessor) serviceAuth(player *Player, message *request.WSRequest) {
	switch message.Resource {
	case "login":
		g.resourceAuthLogin(player, message)

	default:
		g.sendResponseError(player, message, "auth resource not found", nil)
	}
}
