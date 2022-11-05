package saloon

import (
	"log"

	"github.com/dev-rodrigobaliza/carteado/consts"
	"github.com/dev-rodrigobaliza/carteado/domain/config"
	pl "github.com/dev-rodrigobaliza/carteado/domain/core/player"
	"github.com/dev-rodrigobaliza/carteado/domain/request"
	"github.com/dev-rodrigobaliza/carteado/internal/api/v1/adaptors/services"
	"github.com/dev-rodrigobaliza/carteado/internal/core/player"
	"github.com/dev-rodrigobaliza/carteado/internal/core/table"
	"github.com/dev-rodrigobaliza/carteado/internal/security/paseto"
	"github.com/dev-rodrigobaliza/carteado/pkg/safemap"
	"github.com/dev-rodrigobaliza/carteado/utils"
	"github.com/gofiber/websocket/v2"
)

var (
	Security *paseto.PasetoMaker
)

type Saloon struct {
	cfg        *config.App
	appService *services.AppService
	tables     *safemap.SafeMap[string, *table.Table]
	players    *safemap.SafeMap[string, *player.Player]
	boardChan  chan pl.Message[*player.Player] // messages received from players
	addChan    chan *player.Player
	delChan    chan *player.Player
	doneChan   chan struct{}
}

func New(cfg *config.App, appService *services.AppService) *Saloon {
	saloon := &Saloon{
		cfg:        cfg,
		appService: appService,
		tables:     safemap.New[string, *table.Table](),
		players:    safemap.New[string, *player.Player](),
		boardChan:  make(chan pl.Message[*player.Player], consts.TABLEMANAGER_MESSAGE_STACK_SIZE),
		addChan:    make(chan *player.Player, consts.TABLEMANAGER_MAX_PLAYERS),
		delChan:    make(chan *player.Player, consts.TABLEMANAGER_MAX_PLAYERS),
		doneChan:   make(chan struct{}),
	}

	go saloon.Start()

	return saloon
}

func (s *Saloon) NewPlayer(conn *websocket.Conn) {
	p := player.New(conn, s.boardChan, s.delChan)
	s.greetingMesssage(p, "welcome stranger")
	s.addChan <- p
	p.Listen()
}

func (s *Saloon) ProcessPlayerMessage(player *player.Player, message request.WSRequest) {
	if message.Service == "auth" {
		s.serviceAuth(player, &message)
		return
	}
	if player.User == nil {
		s.sendResponseError(player, &message, "player unauthenticated", nil)
		return
	}
	switch message.Service {
	case "admin":
		s.serviceAdmin(player, &message)

	case "table":
		s.serviceTable(player, &message)

	default:
		s.sendResponseError(player, &message, "service not found", nil)
	}
}

func (s *Saloon) Start() {
	log.Printf("*** saloon is opened")
	go s.processMessages()

	for {
		select {
		case player := <-s.addChan:
			s.players.Insert(player.UUID, player)
			s.debug("+++ hub - new player connected\t%s", player)
			s.debug("*** server status: %v", s.getServerStatusResponse(false))

		case player := <-s.delChan:
			// TODO (@dev-rodrigobaliza) remove from table, group and other stuff
			err := s.players.DeleteKey(player.UUID)
			if err != nil {
				s.debug("error removing player: %s", err.Error())
			} else {
				s.debug("--- hub - player disconnected\t%s", player)
			}
			s.debug("*** server status: %v", s.getServerStatusResponse(false))

		case <-s.doneChan:
			return
		}
	}
}

func (s *Saloon) Stop() {
	s.doneChan <- struct{}{}

	close(s.addChan)
	close(s.delChan)
	close(s.doneChan)

	log.Printf("*** saloon is closed")
}

func (s *Saloon) addTable(table *table.Table) {
	s.tables.Insert(table.GetID(), table)
}

func (s *Saloon) delTable(table *table.Table) error {
	return s.tables.DeleteKey(table.GetID())
}

func (s *Saloon) getTable(tableID string) (*table.Table, error) {
	t, err := s.tables.GetOneValue(tableID, false)
	if err != nil {
		return nil, err
	}

	return t, nil
}

func (s *Saloon) processMessages() {
	for {
		message := <-s.boardChan

		var wsMessage request.WSRequest
		err := wsMessage.FromBytes(message.Data)
		if err != nil {
			log.Printf("!!! message from [%s] - error parsing player websocket [%s]", message.Player, message.Data)
			s.sendResponseError(message.Player, nil, "invalid message", err)
		} else {
			s.debug("--- message from [%s] - [%s]", message.Player, utils.CompactJson(message.Data))
			s.ProcessPlayerMessage(message.Player, wsMessage)
		}
	}
}
