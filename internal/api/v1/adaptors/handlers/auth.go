package handlers

import (
	"fmt"
	"strings"

	"github.com/dev-rodrigobaliza/carteado/domain/request"
	"github.com/dev-rodrigobaliza/carteado/internal/api/v1/ports"
	"github.com/dev-rodrigobaliza/carteado/internal/security/paseto"
	"github.com/dev-rodrigobaliza/carteado/utils"
	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	authService ports.IAuthService
}

var (
	Security *paseto.PasetoMaker
	// This line is for get feedback in case we are not implementing the interface correctly
	_ ports.IAuthHandler = (*AuthHandler)(nil)
)

// NewAuthHandler creates a new instance of AuthHandler
func NewAuthHandler(authService ports.IAuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// LoadRoutes loads the routes for the auth handlers
func (h *AuthHandler) AddRoutes(r fiber.Router) fiber.Router {
	// ATENTION: only here exists api open routes
	original := r
	// auth routes
	open := original.Group("/auth")
	open.Post("/login", h.Login)
	// protected routes
	protected := open.Use(h.isAuthenticated)
	protected.Post("/logout", h.Logout)

	return original.Use(h.isAuthenticated)
}

// Login handles the login request
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	// extract data from request
	var loginData request.Login
	err := c.BodyParser(&loginData)
	if err != nil {
		return utils.SendResponseUnprocessableEntity(c, "expected json with fields: email, password")
	}
	// let the auth service process request
	response, errors, err := h.authService.Login(&loginData, c.IP())
	if err != nil {
		return utils.SendResponseUnauthorized(c)
	}
	if errors != nil {
		return utils.SendResponseValidationError(c, errors)
	}

	return utils.SendResponseSuccess(c, "access granted", response)
}

// Logout handles the logout request
func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	// get id
	id := c.Locals("user_id")
	if id == nil {
		return utils.SendResponseUnauthorized(c)
	}
	// format user id
	str := fmt.Sprintf("%v", id)
	userID, err := utils.StringToUint64(str)
	if err != nil {
		return utils.SendResponseUnauthorized(c)
	}
	// get access token from context
	token := c.Locals("access_token")
	if token == nil {
		return utils.SendResponseUnauthorized(c)
	}
	accessToken := fmt.Sprintf("%v", token)
	// call the auth service - logout
	err = h.authService.Logout(userID, accessToken)
	if err != nil {
		return utils.SendResponseUnauthorized(c)
	}

	return utils.SendResponseSuccess(c, "access revoked", nil)
}

// isAuthenticated verifies if the user is authenticated with the access token bearer in request header:
// {"Authorization": "Bearer <access_token>"}
func (h *AuthHandler) isAuthenticated(c *fiber.Ctx) error {
	// primary validations
	bearerToken := c.Get("Authorization")
	if bearerToken == "" {
		return utils.SendResponseInvalidToken(c)
	}
	preToken := strings.Split(bearerToken, " ")
	if len(preToken) != 2 {
		return utils.SendResponseInvalidToken(c)
	}
	if preToken[0] != "Bearer" {
		return utils.SendResponseInvalidToken(c)
	}
	// token validation
	id, err := Security.VerifyToken(preToken[1])
	if err != nil {
		return utils.SendResponseInvalidToken(c)
	}
	// database validation
	err = h.authService.VerifyToken(id, preToken[1])
	if err != nil {
		return utils.SendResponseInvalidToken(c)
	}

	// set important values
	c.Locals("access_token", preToken[1])
	c.Locals("user_id", id)

	return c.Next()
}
