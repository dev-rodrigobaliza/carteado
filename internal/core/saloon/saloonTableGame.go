package saloon

import (
	"github.com/dev-rodrigobaliza/carteado/domain/request"
	"github.com/dev-rodrigobaliza/carteado/internal/core/player"
	"github.com/dev-rodrigobaliza/carteado/internal/core/table"
)

func (s *Saloon) actionTableGameStart(player *player.Player, message *request.WSRequest, table *table.Table) map[string]interface{} {
	// database validation
	// TODO (@dev-rodrigobaliza) database validation ???
	// onwnership validation
	if !player.User.IsAdmin && table.GetOwner() != player.UUID && table.GetOwner() != "" {
		s.sendResponseError(player, message, "table id not owned by player", nil)
		return nil
	}
	// start table
	err := table.Start()
	if err != nil {
		s.sendResponseError(player, message, "table not started", err)
		return nil
	}
	// make response
	response := make(map[string]interface{})
	response["game_id"] = player.TableID

	return response
}

func (s *Saloon) actionTableGameStop(player *player.Player, message *request.WSRequest, table *table.Table) map[string]interface{} {
	// database validation
	// TODO (@dev-rodrigobaliza) database validation ???
	// onwnership validation
	if !player.User.IsAdmin && table.GetOwner() != player.UUID && table.GetOwner() != "" {
		s.sendResponseError(player, message, "table id not owned by player", nil)
		return nil
	}
	// start table
	err := table.Stop(true)
	if err != nil {
		s.sendResponseError(player, message, "table not stopped", err)
		return nil
	}
	// make response
	response := make(map[string]interface{})
	response["game_id"] = player.TableID

	return response
}
