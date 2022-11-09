package saloon

import (
	"github.com/dev-rodrigobaliza/carteado/domain/request"
	"github.com/dev-rodrigobaliza/carteado/internal/core/player"
	"github.com/dev-rodrigobaliza/carteado/internal/core/table"
)

func (s *Saloon) actionTableGameStart(pl *player.Player, req *request.WSRequest, tb *table.Table) map[string]interface{} {
	// database validation
	// TODO (@dev-rodrigobaliza) database validation ???
	// onwnership validation
	if !pl.User.IsAdmin && tb.GetOwner() != pl.UUID && tb.GetOwner() != "" {
		s.sendResponseError(pl, req, "table not owned by player", nil)
		return nil
	}
	// start table
	err := tb.Start()
	if err != nil {
		s.sendResponseError(pl, req, "table not started", err)
		return nil
	}
	// make response
	response := make(map[string]interface{})
	response["game_id"] = pl.TableID

	return response
}

func (s *Saloon) actionTableGameStop(pl *player.Player, req *request.WSRequest, tb *table.Table) map[string]interface{} {
	// database validation
	// TODO (@dev-rodrigobaliza) database validation ???
	// onwnership validation
	if !pl.User.IsAdmin && tb.GetOwner() != pl.UUID && tb.GetOwner() != "" {
		s.sendResponseError(pl, req, "table not owned by player", nil)
		return nil
	}
	// start table
	err := tb.Stop(true)
	if err != nil {
		s.sendResponseError(pl, req, "table not stopped", err)
		return nil
	}
	// make response
	response := make(map[string]interface{})
	response["game_id"] = pl.TableID

	return response
}
