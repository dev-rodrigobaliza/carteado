package saloon

import (
	"github.com/dev-rodrigobaliza/carteado/domain/request"
	"github.com/dev-rodrigobaliza/carteado/internal/core/player"
	"github.com/dev-rodrigobaliza/carteado/utils"
)

func (s *Saloon) resourceAuthLogin(player *player.Player, message *request.WSRequest) {
	token, ok := message.Data["token"].(string)
	if !ok {
		s.sendResponseError(player, message, "token invalid", nil)
	}
	// token validation
	id, err := Security.VerifyToken(token)
	if err != nil {
		s.sendResponseError(player, message, "token invalid", err)
		return
	}
	// database validation
	err = s.appService.AuthService.VerifyToken(id, token)
	if err != nil {
		s.sendResponseError(player, message, "token invalid", err)
		return
	}
	// get user from database
	userID, _ := utils.StringToUint64(id)
	u := &request.GetUser{
		ID: userID,
	}
	user, _, err := s.appService.UserService.Get(u)
	if err != nil || user == nil {
		s.sendResponseError(player, message, "user not found", err)
		return
	}
	// set player-user information
	firstLogin := player.Login(user)
	s.loginPlayer(player)
	// make response
	response := make(map[string]interface{})
	response["first_login"] = firstLogin
	response["player"] = player.ToResponse()
	// send response
	s.sendResponseSuccess(player, message, "authenticated", response)
	// debug log
	s.debug("=== auth login %v", response)
}

func (s *Saloon) serviceAuth(player *player.Player, message *request.WSRequest) {
	switch message.Resource {
	case "login":
		s.resourceAuthLogin(player, message)

	default:
		s.sendResponseError(player, message, "auth resource not found", nil)
	}
}
