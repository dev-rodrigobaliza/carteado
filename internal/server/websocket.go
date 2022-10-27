package server

import (
	"github.com/dev-rodrigobaliza/carteado/internal/api/v1/adaptors/services"
	"github.com/dev-rodrigobaliza/carteado/internal/core/table"
	"github.com/gofiber/fiber/v2"
	fws "github.com/gofiber/websocket/v2"
)

func (s *Server) hubLoad(r fiber.Router, appService *services.AppService) {
	s.hub = table.NewHub(s.config, appService)

	r.Use(s.upgradeHandler)
	// websocket connection handler
	r.Get("", fws.New(func(c *fws.Conn) {
		defer func() {
			c.Close()
		}()

		table.NewPlayer(s.hub, c)
	}))
}
