package player

import (
	"fmt"
	"strconv"
	"time"

	"github.com/dev-rodrigobaliza/carteado/consts"
	pl "github.com/dev-rodrigobaliza/carteado/domain/core/player"
	"github.com/dev-rodrigobaliza/carteado/domain/request"
	"github.com/dev-rodrigobaliza/carteado/domain/response"
	"github.com/dev-rodrigobaliza/carteado/utils"
	"github.com/gofiber/websocket/v2"
)

type Player struct {
	// monted on create
	UUID      string
	conn      *websocket.Conn
	send      chan []byte
	boardChan chan pl.Message[*Player]
	delChan   chan *Player
	// mounted after user login
	IsBot     bool
	User      *response.User
	createdAt time.Time
	loggedAt  time.Time
	// mounted after enter table
	TableID string
	// mounted after enter group
	GroupID int
	Action  string
}

func New(conn *websocket.Conn, boardChan chan pl.Message[*Player], delChan chan *Player) *Player {
	player := &Player{
		UUID:      utils.NewUUID(consts.PLAYER_PREFIX_ID),
		conn:      conn,
		send:      make(chan []byte, consts.PLAYER_MESSAGE_STACK_SIZE),
		boardChan: boardChan,
		delChan:   delChan,
		IsBot:     false,
		createdAt: time.Now(),
	}

	go player.write()

	return player
}

func NewBot() *Player {
	player := &Player{
		UUID:      utils.NewUUID(consts.PLAYER_PREFIX_ID),
		IsBot:     true,
		createdAt: time.Now(),
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

func (p *Player) SendResponse(request *request.WSRequest, status, message string, data map[string]interface{}) {
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

	p.send <- response.ToBytes()
}

func (p *Player) String() string {
	var name string
	if p.IsBot {
		name = "bot [" + p.UUID + "]"
	} else {
		if p.User == nil {
			name = "# unauthenticated #"
		} else {
			name = p.User.Name
		}
	}

	return name
}

func (p *Player) Response(showAddress, showTableID bool) *response.Player {
	var address string
	var name string
	var tableID string
	var logged string

	if p.IsBot {
		name = "bot [" + p.UUID + "]"
	} else {
		if showAddress {
			address = p.conn.RemoteAddr().String()
		}
		if p.User == nil {
			name = "# unauthenticated #"
		} else {
			name = p.User.Name
			logged = fmt.Sprintf("%d", p.loggedAt.UnixMilli())
		}
	}
	if showTableID {
		tableID = p.TableID
	}
	created := fmt.Sprintf("%d", p.createdAt.UnixMilli())
	groupID := strconv.Itoa(p.GroupID)
	if groupID == "0" {
		groupID = ""
	}

	player := response.NewPlayer(p.UUID, address, name, tableID, groupID, created, logged, p.IsBot)

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
