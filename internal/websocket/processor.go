package websocket

import (
	"time"

	"github.com/dev-rodrigobaliza/carteado/domain/request"
	"github.com/dev-rodrigobaliza/carteado/domain/response"
	"github.com/dev-rodrigobaliza/carteado/internal/api/v1/ports"
	"github.com/dev-rodrigobaliza/carteado/internal/security/paseto"
	"github.com/dev-rodrigobaliza/carteado/utils"
)

var (
	Security    *paseto.PasetoMaker
	AuthService ports.IAuthService
)

type Processor struct{}

func NewProcessor() *Processor {
	return &Processor{}
}

func (p *Processor) Run(player *Player, message request.WSRequest) {
	switch message.Service {
	case "auth":
		p.serviceAuth(player, &message)

	default:
		if !player.auth {
			p.sendResponse(player, nil, "error", "player unauthenticated", nil)
			return
		}

		p.sendResponse(player, nil, "error", "service not found", nil)
	}
}

func (p *Processor) resourceLogin(player *Player, message *request.WSRequest) {
	token, ok := message.Data["token"].(string)
	if !ok {
		p.sendResponse(player, message, "error", "token invalid", nil)
	}
	// token validation
	id, err := Security.VerifyToken(token)
	if err != nil {
		p.sendResponse(player, message, "error", "token invalid", nil)
		return
	}
	// database validation
	err = AuthService.VerifyToken(id, token)
	if err != nil {
		p.sendResponse(player, message, "error", "token invalid", nil)
		return
	}

	if !player.auth {
		player.auth = true
		player.id = utils.NewUUID()
	}
	player.since = time.Now()
	
	response := make(map[string]interface{})
	response["id"] = player.id

	p.sendResponse(player, message, "success", "authenticated", response)
}

func (p *Processor) serviceAuth(player *Player, message *request.WSRequest) {
	switch message.Resource {
	case "login":
		p.resourceLogin(player, message)

	default:
		p.sendResponse(player, nil, "error", "resource not found", nil)
	}
}

func (p *Processor) sendResponse(player *Player, request *request.WSRequest, status, message string, data map[string]interface{}) {
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
