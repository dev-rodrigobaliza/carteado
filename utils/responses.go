package utils

import (
	"github.com/dev-rodrigobaliza/carteado/domain/response"
	"github.com/gofiber/fiber/v2"
)

func SendResponseBadRequest(c *fiber.Ctx, err error) error {
	return SendResponseFull(c, fiber.StatusBadRequest, "error", err.Error(), nil)
}

func SendResponseError(c *fiber.Ctx, statusCode int, message string) error {
	return SendResponseFull(c, statusCode, "error", message, nil)
}

func SendResponseFull(c *fiber.Ctx, statusCode int, status, message string, data interface{}) error {
	response := fiber.Map{
		"message": message,
		"status":  status,
	}

	if data != nil {
		response["data"] = data
	}

	return c.Status(statusCode).JSON(&response)
}

func SendResponseInvalidToken(c *fiber.Ctx) error {
	return SendResponseFull(c, fiber.StatusUnauthorized, "error", "invalid token", nil)
}

func SendResponseNotImplemented(c *fiber.Ctx) error {
	return SendResponseFull(c, fiber.StatusNotImplemented, "error", "not implemented", nil)
}

func SendResponseSuccess(c *fiber.Ctx, message string, data interface{}) error {
	return SendResponseFull(c, fiber.StatusOK, "success", message, data)
}

func SendResponseUnauthorized(c *fiber.Ctx) error {
	return SendResponseFull(c, fiber.StatusUnauthorized, "error", "invalid credentials", nil)
}

func SendResponseUnprocessableEntity(c *fiber.Ctx, message string) error {
	return SendResponseFull(c, fiber.StatusUnprocessableEntity, "error", message, nil)
}

func SendResponseValidationError(c *fiber.Ctx, errors []*response.ErrorValidation) error {
	return SendResponseFull(c, fiber.StatusUnprocessableEntity, "error", "validation error", errors)
}
