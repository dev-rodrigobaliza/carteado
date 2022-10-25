package websocket

import (
	"log"
	"sync/atomic"

	"github.com/dev-rodrigobaliza/carteado/domain/config"
	"github.com/dev-rodrigobaliza/carteado/domain/request"
	"github.com/dev-rodrigobaliza/carteado/domain/response"
)

const BUFFER_SIZE = 1024

type Hub struct {
	cfg          *config.App
	playersCount int64
	broadcast    chan []byte
	wsMessage    chan *WSMessage
	add          chan *Player
	remove       chan *Player
	done         chan struct{}
	players      map[*Player]bool
	processor    *Processor
}

func NewHub(cfg *config.App) *Hub {
	hub := &Hub{
		cfg:          cfg,
		playersCount: 0,
		broadcast:    make(chan []byte),
		wsMessage:    make(chan *WSMessage, BUFFER_SIZE),
		add:          make(chan *Player),
		remove:       make(chan *Player),
		done:         make(chan struct{}),
		players:      make(map[*Player]bool),
		processor:    NewProcessor(),
	}
	go hub.Run()

	return hub
}

func (h *Hub) Run() {
	go h.processMessages()

	for {
		select {
		case player := <-h.add:
			h.debug("+++ new player connected\t%s", player)
			h.players[player] = true
			atomic.AddInt64(&h.playersCount, 1)
			h.sendResponse(player, nil, "success", "welcome player", h.welcomeMesssage())

		case player := <-h.remove:
			if h.players[player] {
				h.debug("--- player disconnected\t%s", player)
				delete(h.players, player)
				atomic.AddInt64(&h.playersCount, -1)
			}

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

func (h *Hub) bufferFull(player *Player) {
	h.debug("!!! buffer full, player disconnected\t%s", player)
	close(player.send)
	delete(h.players, player)
	atomic.AddInt64(&h.playersCount, -1)
}

func (h *Hub) debug(format string, v ...any) {
	if h.cfg.Debug {
		log.Printf(format, v...)
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
			h.sendResponse(player, nil, "error", "invalid message", nil)
		} else {
			h.debug("--- message [%s] from [%s]", wsMessage.Data, player)
			go h.processor.Run(player, message)
		}
	}
}

func (h *Hub) sendAll(message []byte) {
	for player := range h.players {
		if player.auth {
			select {
			case player.send <- message:
			default:
				h.bufferFull(player)
			}
		}
	}
}

func (h *Hub) sendOne(wsMessage *WSMessage) {
	for player := range h.players {
		if wsMessage.Player == player {
			select {
			case player.send <- wsMessage.Data:
			default:
				h.bufferFull(player)
			}
		}
	}
}

func (h *Hub) sendResponse(player *Player, request *request.WSRequest, status, message string, data map[string]interface{}) {
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

func (h *Hub) welcomeMesssage() map[string]interface{} {
	message := make(map[string]interface{})
	message["server"] = h.cfg.Name
	message["version"] = h.cfg.Version
	if h.cfg.Debug {
		message["players"] = h.playersCount
	}

	return message
}
