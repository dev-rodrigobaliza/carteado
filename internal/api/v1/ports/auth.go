package ports

import (
	"github.com/dev-rodrigobaliza/carteado/domain/request"
	"github.com/dev-rodrigobaliza/carteado/domain/response"
	"github.com/gofiber/fiber/v2"
)

// handlers to api
type IAuthHandler interface {
	AddRoutes(r fiber.Router) fiber.Router
	Login(c *fiber.Ctx) error
	Logout(c *fiber.Ctx) error
}

// all services will consume this interface
type IAuthService interface {
	Login(login *request.Login, ip string) (*response.Login, []*response.ErrorValidation, error)
	Logout(userID uint64, accessToken string) error
	VerifyToken(userID, accessToken string) error
}

// data manipulation (read/write) in cache and database
type IAuthRepository interface {
	Login(userID uint64, accessToken string) error
	Logout(userID uint64, accessToken string) error
	VerifyToken(userID uint64, accessToken string) error
}
