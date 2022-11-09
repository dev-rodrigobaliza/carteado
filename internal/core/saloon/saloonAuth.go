package saloon

import (
	"github.com/dev-rodrigobaliza/carteado/domain/request"
	"github.com/dev-rodrigobaliza/carteado/internal/core/player"
	"github.com/dev-rodrigobaliza/carteado/utils"
)

func (s *Saloon) resourceAuthLogin(pl *player.Player, req *request.WSRequest) {
	token, ok := req.Data["token"].(string)
	if !ok {
		s.sendResponseError(pl, req, "token invalid", nil)
	}
	// token validation
	id, err := Security.VerifyToken(token)
	if err != nil {
		s.sendResponseError(pl, req, "token invalid", err)
		return
	}
	// database validation
	err = s.appService.AuthService.VerifyToken(id, token)
	if err != nil {
		s.sendResponseError(pl, req, "token invalid", err)
		return
	}
	// get user from database
	userID, _ := utils.StringToUint64(id)
	u := &request.GetUser{
		ID: userID,
	}
	user, _, err := s.appService.UserService.Get(u)
	if err != nil || user == nil {
		s.sendResponseError(pl, req, "user not found", err)
		return
	}
	// set player-user information
	firstLogin := pl.Login(user)
	s.loginPlayer(pl)
	// make response
	response := make(map[string]interface{})
	response["first_login"] = firstLogin
	response["player"] = pl.ToResponse(true, pl.User.IsAdmin)
	// send response
	s.sendResponseSuccess(pl, req, pl.Greeting(), response)
	if pl.User.IsAdmin {
		pl.SendResponse(nil, "info", "server status", s.getServerStatusResponse(false))
	}
	// debug log
	s.debug("=== auth login %v", response)
}

func (s *Saloon) serviceAuth(pl *player.Player, req *request.WSRequest) {
	switch req.Resource {
	case "login":
		s.resourceAuthLogin(pl, req)

	default:
		s.sendResponseError(pl, req, "auth resource not found", nil)
	}
}
