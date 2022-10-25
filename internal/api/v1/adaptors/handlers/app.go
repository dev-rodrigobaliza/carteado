package handlers

import (
	"github.com/dev-rodrigobaliza/carteado/internal/api/v1/adaptors/services"
	"github.com/dev-rodrigobaliza/carteado/internal/api/v1/ports"
	"github.com/gofiber/fiber/v2"
)

type AppHandler struct {
	AuthHandler ports.IAuthHandler
	UserHandler ports.IUserHandler
}

// Load creates a new AppHandler instance
func Load(appService *services.AppService, api fiber.Router) {
	authHandler := NewAuthHandler(appService.AuthService)
	userHandler := NewUserHandler(appService.UserService)

	a := &AppHandler{
		AuthHandler: authHandler,
		UserHandler: userHandler,
	}
	a.loadRoutesV1(api)
}

func (h *AppHandler) loadRoutesV1(r fiber.Router) {
	v1 := r.Group("/v1")
	
	protectedRoute := h.AuthHandler.AddRoutes(v1)
	h.UserHandler.AddRoutes(protectedRoute)
}
