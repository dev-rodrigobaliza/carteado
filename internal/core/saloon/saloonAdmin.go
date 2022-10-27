package saloon

import (
	"github.com/dev-rodrigobaliza/carteado/domain/request"
	"github.com/dev-rodrigobaliza/carteado/errors"
	pl "github.com/dev-rodrigobaliza/carteado/internal/core/player"
)

func (s *Saloon) resourceAdminStatus(player *pl.Player, message *request.WSRequest) {
	authenticatedOnly, ok := message.Data["authenticated_only"].(bool)
	if !ok {
		s.sendResponseError(player, message, "authenticated only invalid", nil)
	}
	response := s.getServerStatusResponse(authenticatedOnly)

	s.sendResponseSuccess(player, message, "server status", response)
}

func (s *Saloon) serviceAdmin(player *pl.Player, message *request.WSRequest) {
	// basic validation (admins only)
	if !player.User.IsAdmin {
		s.sendResponseError(player, message, "unauthorized", errors.ErrUnauthorized)
		return
	}

	switch message.Resource {
	case "status":
		s.resourceAdminStatus(player, message)

	default:
		s.sendResponseError(player, message, "auth resource not found", nil)
	}
}
