package handlers

import (
	"github.com/dev-rodrigobaliza/carteado/domain/request"
	"github.com/dev-rodrigobaliza/carteado/internal/api/v1/ports"
	"github.com/dev-rodrigobaliza/carteado/utils"
	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	userService ports.IUserService
}

var _ ports.IUserHandler = (*UserHandler)(nil)

// NewAuthHandler creates a new instance of AuthHandler
func NewUserHandler(userService ports.IUserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// LoadRoutes loads the routes for the auth handlers
func (h *UserHandler) AddRoutes(r fiber.Router) {
	user := r.Group("/user")
	user.Post("/add", h.Add)
	user.Post("/delete", h.Delete)
	user.Post("/get", h.Get)
	user.Post("/update", h.Update)
}

func (h *UserHandler) Add(c *fiber.Ctx) error {
	// extract data from request
	var userData request.User
	err := c.BodyParser(&userData)
	if err != nil {
		return utils.SendResponseUnprocessableEntity(c, "expected json with fields: name, email, password")
	}
	// let the user service process request
	response, errors, err := h.userService.Add(&userData)
	if err != nil {
		return utils.SendResponseBadRequest(c)
	}
	if errors != nil {
		return utils.SendResponseValidationError(c, errors)
	}

	return utils.SendResponseSuccess(c, "user created", response)
}

func (h *UserHandler) Delete(c *fiber.Ctx) error {
	return utils.SendResponseNotImplemented(c)
}

func (h *UserHandler) Get(c *fiber.Ctx) error {
	return utils.SendResponseNotImplemented(c)
}

func (h *UserHandler) Update(c *fiber.Ctx) error {
	return utils.SendResponseNotImplemented(c)
}
