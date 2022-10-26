package server

import (
	"github.com/dev-rodrigobaliza/carteado/internal/api/v1/adaptors/services"
	"github.com/dev-rodrigobaliza/carteado/internal/websocket"
	"github.com/gofiber/fiber/v2"
	fws "github.com/gofiber/websocket/v2"
)

func (s *Server) websocketLoad(r fiber.Router, appService *services.AppService) {
	s.ws = websocket.NewHub(s.config, appService)

	r.Use(s.upgradeHandler)
	// websocket connection handler
	r.Get("", fws.New(func(c *fws.Conn) {
		defer func() {
			c.Close()
		}()

		websocket.NewPlayer(s.ws, c)
	}))
}
