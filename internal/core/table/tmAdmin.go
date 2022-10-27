package table

import (
	"github.com/dev-rodrigobaliza/carteado/domain/request"
	"github.com/dev-rodrigobaliza/carteado/errors"
)

func (t *TableManager) resourceAdminStatus(player *Player, message *request.WSRequest) {
	authenticatedOnly, ok := message.Data["authenticated_only"].(bool)
	if !ok {
		t.sendResponseError(player, message, "authenticated only invalid", nil)
	}	
	response := t.getServerStatusResponse(authenticatedOnly)

	t.sendResponseSuccess(player, message, "server status", response)
}

func (t *TableManager) serviceAdmin(player *Player, message *request.WSRequest) {
	// basic validation (admins only)
	if !player.user.IsAdmin {
		t.sendResponseError(player, message, "unauthorized", errors.ErrUnauthorized)
		return
	}

	switch message.Resource {
	case "status":
		t.resourceAdminStatus(player, message)

	default:
		t.sendResponseError(player, message, "auth resource not found", nil)
	}
}