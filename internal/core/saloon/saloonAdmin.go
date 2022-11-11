package saloon

import (
	"github.com/dev-rodrigobaliza/carteado/domain/request"
	"github.com/dev-rodrigobaliza/carteado/errors"
	"github.com/dev-rodrigobaliza/carteado/internal/core/player"
)

func (s *Saloon) resourceAdminStatus(pl *player.Player, req *request.WSRequest) {
	response := s.getServerStatusResponse(pl.User.IsAdmin)

	s.sendResponseSuccess(pl, req, "status server", response)
}

func (s *Saloon) serviceAdmin(pl *player.Player, req *request.WSRequest) {
	// basic validation (admins only)
	if !pl.User.IsAdmin {
		s.sendResponseError(pl, req, "unauthorized", errors.ErrUnauthorized)
		return
	}

	switch req.Resource {
	case "status":
		s.resourceAdminStatus(pl, req)

	default:
		s.sendResponseError(pl, req, "auth resource not found", nil)
	}
}
