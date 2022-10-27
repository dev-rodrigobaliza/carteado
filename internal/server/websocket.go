package server

import (
	"github.com/dev-rodrigobaliza/carteado/internal/api/v1/adaptors/services"
	"github.com/dev-rodrigobaliza/carteado/internal/core/saloon"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

func (s *Server) initSaloon(r fiber.Router, appService *services.AppService) {
	s.saloon = saloon.NewSaloon(s.config, appService)

	r.Use(s.upgradeHandler)
	// websocket connection handler
	r.Get("", websocket.New(func(c *websocket.Conn) {
		defer func() {
			c.Close()
		}()

		s.saloon.NewPlayer(c)
	}))
}
