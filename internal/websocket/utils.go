package websocket

import (
	"log"

	"github.com/dev-rodrigobaliza/carteado/domain/request"
	"github.com/dev-rodrigobaliza/carteado/domain/response"
)

func (g *GameProcessor) debug(format string, v ...any) {
	if g.cfg.Debug {
		log.Printf(format, v...)
	}
}

func (g *GameProcessor) getGameID(message *request.WSRequest) string {
	gameID, ok := message.Data["game_id"].(string)
	if !ok {
		return ""
	}

	return gameID
}


func (g *GameProcessor) sendResponse(player *Player, request *request.WSRequest, status, message string, data map[string]interface{}) {
	response := &response.WSResponse{
		Status:  status,
		Message: message,
	}
	if request != nil {
		response.RequestID = request.RequestID
	}
	if len(data) > 0 {
		response.Data = data
	}

	player.Send(response.ToBytes())
}

func (g *GameProcessor) sendResponseError(player *Player, request *request.WSRequest, message string, err error) {
	var data map[string]interface{}
	if err != nil && g.cfg.Debug {
		data = make(map[string]interface{})
		data["error"] = err.Error()
	}

	g.sendResponse(player, request, "error", message, data)
}

func (g *GameProcessor) sendResponseSuccess(player *Player, request *request.WSRequest, message string, data map[string]interface{}) {
	g.sendResponse(player, request, "success", message, data)
}
