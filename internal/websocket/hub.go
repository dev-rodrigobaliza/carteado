package websocket

import (
	"log"

	"github.com/dev-rodrigobaliza/carteado/domain/config"
	"github.com/dev-rodrigobaliza/carteado/domain/request"
	"github.com/dev-rodrigobaliza/carteado/internal/api/v1/adaptors/services"
	"github.com/dev-rodrigobaliza/carteado/pkg/safemap"
)

const BUFFER_SIZE = 1024

type Hub struct {
	cfg           *config.App
	players       *safemap.SafeMap[*Player, bool]
	broadcast     chan []byte
	wsMessage     chan *WSMessage
	add           chan *Player
	remove        chan *Player
	done          chan struct{}
	gameProcessor *GameProcessor
}

func NewHub(cfg *config.App, appService *services.AppService) *Hub {
	players := safemap.New[*Player, bool]()
	gameProcessor := NewGameProcessor(cfg, players, appService)

	hub := &Hub{
		cfg:           cfg,
		players:       players,
		broadcast:     make(chan []byte),
		wsMessage:     make(chan *WSMessage, BUFFER_SIZE),
		add:           make(chan *Player),
		remove:        make(chan *Player),
		done:          make(chan struct{}),
		gameProcessor: gameProcessor,
	}
	go hub.Run()

	return hub
}

func (h *Hub) Run() {
	go h.processMessages()

	for {
		select {
		case player := <-h.add:
			h.addPlayer(player)
			go h.greetingMesssage(player, "welcome player")
			h.debug("+++ new player connected\t%s", player)

		case player := <-h.remove:
			h.delPlayer(player, false)
			h.debug("--- player disconnected\t%s", player)

		case message := <-h.broadcast:
			h.sendAll(message)

		case <-h.done:
			return
		}
	}
}

func (h *Hub) Stop() {
	h.done <- struct{}{}

	close(h.broadcast)
	close(h.wsMessage)
	close(h.add)
	close(h.remove)
}

func (h *Hub) addPlayer(player *Player) {
	h.players.Insert(player, true)
}

func (h *Hub) debug(format string, v ...any) {
	if h.cfg.Debug {
		log.Printf(format, v...)
	}
}

func (h *Hub) delPlayer(player *Player, closePlayer bool) error {
	if closePlayer {
		close(player.send)
	}

	return h.players.Delete(player)
}

func (h *Hub) loginPlayer(player *Player) {
	var welcomePlayer *Player
	// if player has previous login and remove it
	players := h.players.GetAllKeys()
	for _, p := range players {
		if p != player && p.user != nil && p.user.ID == player.user.ID {
			p.user = nil
			welcomePlayer = p
			break
		}
	}

	if welcomePlayer != nil {
		// TODO (@dev-rodrigobaliza) send message about login another place
		h.greetingMesssage(welcomePlayer, "disconnected (using another session)")
	}
}

func (h *Hub) processMessages() {
	for {
		wsMessage := <-h.wsMessage
		player := wsMessage.Player
		var message request.WSRequest
		err := message.FromBytes(wsMessage.Data)
		if err != nil {
			log.Printf("!!! error parsing player websocket message [%s] from %s", wsMessage.Data, player)
			h.gameProcessor.sendResponseError(player, nil, "invalid message", err)
		} else {
			h.debug("--- message received [%s] from [%s]", wsMessage.Data, player)
			h.gameProcessor.ProcessPlayerMessage(player, message)
		}
	}
}

func (h *Hub) sendAll(message []byte) {
	players := h.players.GetAllKeys()
	for _, player := range players {
		if player.user != nil {
			select {
			case player.send <- message:
			default:
				h.delPlayer(player, true)
			}
		}
	}
}

func (h *Hub) sendOne(wsMessage *WSMessage) {
	players := h.players.GetAllKeys()
	for _, player := range players {
		if wsMessage.Player == player {
			select {
			case player.send <- wsMessage.Data:
			default:
				h.delPlayer(player, true)
			}

			return
		}
	}
}

func (h *Hub) greetingMesssage(player *Player, message string) {
	// only authenticated players
	players := make([]string, 0)
	allPlayers := h.players.GetAllKeys()
	for _, player := range allPlayers {
		if player.user != nil {
			players = append(players, player.uuid)
		}
	}

	response := make(map[string]interface{})
	response["server"] = h.cfg.Name
	response["version"] = h.cfg.Version
	response["players_count"] = len(players)
	if h.cfg.Debug {
		response["players"] = players
	}

	h.gameProcessor.sendResponseSuccess(player, nil, message, response)
}
