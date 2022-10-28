package player

import (
	"fmt"
	"time"

	"github.com/dev-rodrigobaliza/carteado/consts"
	pl "github.com/dev-rodrigobaliza/carteado/domain/core/player"
	"github.com/dev-rodrigobaliza/carteado/domain/response"
	"github.com/dev-rodrigobaliza/carteado/utils"
	"github.com/gofiber/websocket/v2"
)

type Player struct {
	// monted on create
	UUID      string
	Addr      string
	conn      *websocket.Conn
	send      chan []byte
	boardChan chan pl.Message[*Player]
	delChan   chan *Player
	// mounted after user login
	User     *response.User
	loggedAt time.Time
	// mounted after enter table
	TableID string
	// mounted after enter group
	GroupID int
}

func New(conn *websocket.Conn, boardChan chan pl.Message[*Player], delChan chan *Player) *Player {
	player := &Player{
		UUID:      utils.NewUUID(consts.PLAYER_PREFIX_ID),
		Addr:      conn.RemoteAddr().String(),
		conn:      conn,
		send:      make(chan []byte, consts.PLAYER_MESSAGE_STACK_SIZE),
		boardChan: boardChan,
		delChan:   delChan,
	}

	go player.write()

	return player
}

func (p *Player) Greeting() string {
	return fmt.Sprintf("hello %s", p.User.Name)
}

func (p *Player) Listen() {
	defer func() {
		p.delChan <- p
		if p.conn.Conn != nil {
			_ = p.conn.Close()
		}
	}()
	for {
		_, data, err := p.conn.ReadMessage()
		if err != nil {
			break
		}
		if len(data) > 0 {
			message := pl.Message[*Player]{
				Player: p,
				Data:   data,
			}

			p.boardChan <- message
		}
	}
}

func (p *Player) Login(user *response.User) bool {
	userOut := (p.User == nil || p.User.ID != user.ID)
	if userOut {
		p.User = user
		p.loggedAt = time.Now()
	}

	return userOut
}

func (p *Player) Send(data []byte) {
	p.send <- data
}

func (p *Player) String() string {
	return p.Addr
}

func (p *Player) ToResponse() *response.Player {
	var name string
	if p.User == nil {
		name = "# unauthenticated #"
	} else {
		name = p.User.Name
	}
	player := response.NewPlayer(p.UUID, name)

	return player
}

func (p *Player) write() {
	defer func() {
		p.delChan <- p
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
