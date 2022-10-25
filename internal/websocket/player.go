package websocket

import (
	"time"

	"github.com/gofiber/websocket/v2"
)

type Player struct {
	hub   *Hub
	conn  *websocket.Conn
	send  chan []byte
	addr  string
	id    string
	auth  bool
	since time.Time
}

func NewPlayer(hub *Hub, conn *websocket.Conn) {
	player := &Player{
		conn: conn,
		hub:  hub,
		send: make(chan []byte, BUFFER_SIZE),
		addr: conn.RemoteAddr().String(),
		auth: false,
	}
	player.addr = player.String()

	player.hub.add <- player
	go player.write()
	player.read()
}

func (p *Player) Send(data []byte) {
	wsMessage := NewWSMessage(p, data)
	go p.hub.sendOne(wsMessage)
}

func (p *Player) String() string {
	return p.addr
}

func (p *Player) read() {
	defer func() {
		p.hub.remove <- p
		if p.conn.Conn != nil {
			_ = p.conn.Close()
		}
	}()
	for {
		_, message, err := p.conn.ReadMessage()
		if err != nil {
			break
		}
		if len(message) > 0 {
			wsMessage := NewWSMessage(p, message)
			p.hub.wsMessage <- wsMessage
		}
	}
}

func (p *Player) write() {
	defer func() {
		p.hub.remove <- p
		if p.conn.Conn != nil {
			_ = p.conn.Close()
		}
		close(p.send)
	}()
	for message := range p.send {
		err := p.conn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			return
		}
		n := len(p.send)
		for i := 0; i < n; i++ {
			err = p.conn.WriteMessage(websocket.TextMessage, <-p.send)
			if err != nil {
				return
			}
		}
	}
}
