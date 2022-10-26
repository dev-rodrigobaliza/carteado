package websocket

import (
	"log"
	"sync"

	"github.com/dev-rodrigobaliza/carteado/domain/config"
	"github.com/dev-rodrigobaliza/carteado/domain/request"
)

const BUFFER_SIZE = 1024

type Hub struct {
	cfg           *config.App
	playersLock   *sync.RWMutex
	players       map[*Player]bool
	gameProcessor *GameProcessor

	broadcast     chan []byte
	wsMessage     chan *WSMessage
	add           chan *Player
	remove        chan *Player
	done          chan struct{}
}

func NewHub(cfg *config.App) *Hub {
	hub := &Hub{
		cfg:           cfg,
		playersLock:   &sync.RWMutex{},
		players:       make(map[*Player]bool),
		gameProcessor: NewGameProcessor(cfg),

		broadcast:     make(chan []byte),
		wsMessage:     make(chan *WSMessage, BUFFER_SIZE),
		add:           make(chan *Player),
		remove:        make(chan *Player),
		done:          make(chan struct{}),
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
			go h.gameProcessor.sendResponseSuccess(player, nil, "welcome player", h.welcomeMesssage())
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
	h.playersLock.Lock()
	defer h.playersLock.Unlock()

	if !h.players[player] {
		h.players[player] = true
	}
}

func (h *Hub) debug(format string, v ...any) {
	if h.cfg.Debug {
		log.Printf(format, v...)
	}
}

func (h *Hub) delPlayer(player *Player, closePlayer bool) {
	if closePlayer {
		close(player.send)
	}

	h.playersLock.Lock()
	defer h.playersLock.Unlock()

	if h.players[player] {
		delete(h.players, player)
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
			h.gameProcessor.sendResponseError(player, nil, "invalid message")
		} else {
			h.debug("--- message [%s] from [%s]", wsMessage.Data, player)
			go h.gameProcessor.Run(player, message)
		}
	}
}

func (h *Hub) sendAll(message []byte) {
	h.playersLock.RLock()
	defer h.playersLock.RUnlock()

	for player := range h.players {
		if player.auth {
			select {
			case player.send <- message:
			default:
				h.delPlayer(player, true)
			}
		}
	}
}

func (h *Hub) sendOne(wsMessage *WSMessage) {
	h.playersLock.RLock()
	defer h.playersLock.RUnlock()

	for player := range h.players {
		if wsMessage.Player == player {
			select {
			case player.send <- wsMessage.Data:
			default:
				h.delPlayer(player, true)
			}
		}
	}
}

func (h *Hub) welcomeMesssage() map[string]interface{} {
	var players []string
	h.playersLock.RLock()
	playersCount := len(h.players)
	for player := range h.players {
		players = append(players, player.id)
	}
	h.playersLock.RUnlock()

	message := make(map[string]interface{})
	message["server"] = h.cfg.Name
	message["version"] = h.cfg.Version
	message["players_count"] = playersCount
	if h.cfg.Debug {
		message["players"] = players
	}

	return message
}
