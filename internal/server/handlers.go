package server

import (
	"github.com/dev-rodrigobaliza/carteado/utils"
	"github.com/gofiber/fiber/v2"
)

func (s *Server) getHealth(c *fiber.Ctx) error {
	return utils.SendResponseSuccess(c, "healthy", nil)
}
