package ports

import (
	"github.com/dev-rodrigobaliza/carteado/domain/model"
	"github.com/dev-rodrigobaliza/carteado/domain/request"
	"github.com/dev-rodrigobaliza/carteado/domain/response"
	"github.com/gofiber/fiber/v2"
)

// handlers to api
type IUserHandler interface {
	AddRoutes(r fiber.Router)
	Add(c *fiber.Ctx) error
	Delete(c *fiber.Ctx) error
	Get(c *fiber.Ctx) error
	Update(c *fiber.Ctx) error
}

// all services will consume this interface
type IUserService interface {
	Add(user *request.User) (*response.User, []*response.ErrorValidation, error)
	Delete() error
	Get() error
	Update() error
}

// data manipulation (read/write) in cache and database
type IUserRepository interface {
	Add(user *model.User) error
	Delete() error
	FindByEmail(email string) (*model.User, error)
	FindByID(id uint64) (*model.User, error)
	Update() error	
}
