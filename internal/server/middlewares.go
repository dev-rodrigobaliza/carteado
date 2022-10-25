package server

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

func (s *Server) error404() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "error",
			"message": "request not found",
		})
	}
}

func (s *Server) limiterNext(c *fiber.Ctx) bool {
	// skip requests from these addresses
	skip := false
	for _, addr := range s.config.HTTP.Limiter.AllowedIPs {
		if c.IP() == addr {
			skip = true
			break
		}
	}

	return skip
}

func (s *Server) limiterLimitReached(c *fiber.Ctx) error {
	return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
		"status":  "error",
		"message": "too many requests",
	})
}

func (s *Server) timing() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// calculate time spent
		start := time.Now()
		defer func() {
			c.Set("X-Response-Time", time.Since(start).String())
		}()
		// maybe use this information somewhere
		c.Locals("start", start)
		// do the job (request / response)
		return c.Next()
	}
}

func (s *Server) upgradeHandler(c *fiber.Ctx) error {
	if websocket.IsWebSocketUpgrade(c) {
		return c.Next()
	}
	return c.Status(fiber.StatusUpgradeRequired).JSON(fiber.Map{
		"status":  "error",
		"message": "upgrade required",
	})
}

func (s *Server) versioning() fiber.Handler {
	return func(c *fiber.Ctx) error {
		defer func() {
			c.Set("X-Date", s.config.Date)
			c.Set("X-Version", s.config.Version)
		}()
		// do the job (request / response)
		return c.Next()
	}
}
