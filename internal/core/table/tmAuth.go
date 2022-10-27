package table

import (
	"github.com/dev-rodrigobaliza/carteado/domain/request"
	"github.com/dev-rodrigobaliza/carteado/utils"
)

func (t *TableManager) resourceAuthLogin(player *Player, message *request.WSRequest) {
	token, ok := message.Data["token"].(string)
	if !ok {
		t.sendResponseError(player, message, "token invalid", nil)
	}
	// token validation
	id, err := Security.VerifyToken(token)
	if err != nil {
		t.sendResponseError(player, message, "token invalid", err)
		return
	}
	// database validation
	err = t.appService.AuthService.VerifyToken(id, token)
	if err != nil {
		t.sendResponseError(player, message, "token invalid", err)
		return
	}
	// get user from database
	userID, _ := utils.StringToUint64(id)
	u := &request.GetUser{
		ID: userID,
	}
	user, _, err := t.appService.UserService.Get(u)
	if err != nil || user == nil {
		t.sendResponseError(player, message, "user not found", err)
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
	t.sendResponseSuccess(player, message, "authenticated", response)
	// debug log
	t.debug("=== auth login %v", response)
}

func (t *TableManager) serviceAuth(player *Player, message *request.WSRequest) {
	switch message.Resource {
	case "login":
		t.resourceAuthLogin(player, message)

	default:
		t.sendResponseError(player, message, "auth resource not found", nil)
	}
}
