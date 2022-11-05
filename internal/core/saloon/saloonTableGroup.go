package saloon

import (
	"fmt"

	"github.com/dev-rodrigobaliza/carteado/domain/request"
	"github.com/dev-rodrigobaliza/carteado/internal/core/player"
	"github.com/dev-rodrigobaliza/carteado/internal/core/table"
)

func (s *Saloon) actionTableGroupEnter(player *player.Player, message *request.WSRequest, table *table.Table, groupID int) map[string]interface{} {
	err := table.AddGroupPlayer(groupID, player)
	if err != nil {
		s.sendResponseError(player, message, fmt.Sprintf("enter group failed: %s", err.Error()), nil)
		return nil
	} else if table == nil {
		s.sendResponseError(player, message, "table id invalid", nil)
		return nil
	}

	return s.getTableStatus(table)
}

func (s *Saloon) actionTableGroupLeave(player *player.Player, message *request.WSRequest, table *table.Table, groupID int) map[string]interface{} {
	err := table.DelGroupPlayer(groupID, player.UUID)
	if err != nil {
		s.sendResponseError(player, message, fmt.Sprintf("leave group failed: %s", err.Error()), nil)
		return nil
	} else if table == nil {
		s.sendResponseError(player, message, "table id invalid", nil)
		return nil
	}

	return s.getTableStatus(table)
}
