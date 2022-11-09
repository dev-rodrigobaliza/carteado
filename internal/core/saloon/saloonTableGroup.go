package saloon

import (
	"fmt"

	"github.com/dev-rodrigobaliza/carteado/domain/request"
	"github.com/dev-rodrigobaliza/carteado/domain/response"
	"github.com/dev-rodrigobaliza/carteado/internal/core/player"
	"github.com/dev-rodrigobaliza/carteado/internal/core/table"
)

func (s *Saloon) actionTableGroupBot(pl *player.Player, req *request.WSRequest, tb *table.Table, quantity int) map[string]interface{} {
	// onwnership validation
	if !pl.User.IsAdmin && tb.GetOwner() != pl.UUID && tb.GetOwner() != "" {
		s.sendResponseError(pl, req, "table not owned by player", nil)

		return nil
	}
	// create bots
	var botsCreated int
	for i := 0; i < quantity; i++ {
		bot := player.NewBot()
		err := tb.AddPlayerBot(bot)
		if err != nil {
			// ToDo (@dev-rodrigobaliza) log this error
			break
		}
		botsCreated++
	}

	res := s.getTableStatus(tb, pl.User.IsAdmin)
	res["bot"] = response.NewBot(quantity, botsCreated)

	return res
}

func (s *Saloon) actionTableGroupEnter(pl *player.Player, req *request.WSRequest, tb *table.Table, groupID int) map[string]interface{} {
	err := tb.AddGroupPlayer(groupID, pl)
	if err != nil {
		s.sendResponseError(pl, req, fmt.Sprintf("enter group failed: %s", err.Error()), nil)
		return nil
	} else if tb == nil {
		s.sendResponseError(pl, req, "table id invalid", nil)
		return nil
	}

	return s.getTableStatus(tb, pl.User.IsAdmin)
}

func (s *Saloon) actionTableGroupLeave(pl *player.Player, req *request.WSRequest, tb *table.Table, groupID int) map[string]interface{} {
	err := tb.DelGroupPlayer(groupID, pl.UUID)
	if err != nil {
		s.sendResponseError(pl, req, fmt.Sprintf("leave group failed: %s", err.Error()), nil)
		return nil
	} else if tb == nil {
		s.sendResponseError(pl, req, "table id invalid", nil)
		return nil
	}

	return s.getTableStatus(tb, pl.User.IsAdmin)
}
