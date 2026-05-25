package httpx

import (
	"errors"
	"log"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"

	"github.com/murielsilveira/gofus/internal/platform/errs"
)

type errorResponse struct {
	Error string `json:"error"`
}

func JSON(c fiber.Ctx, status int, body any) error {
	return c.Status(status).JSON(body)
}

func Error(c fiber.Ctx, err error) error {
	switch {
	case errors.Is(err, errs.ErrNotFound):
		return JSON(c, fiber.StatusNotFound, errorResponse{Error: err.Error()})
	case errors.Is(err, errs.ErrBadRequest):
		return JSON(c, fiber.StatusBadRequest, errorResponse{Error: err.Error()})
	default:
		log.Printf("request error: %v", err)
		return JSON(c, fiber.StatusInternalServerError, errorResponse{Error: "internal server error"})
	}
}

func ParseUUID(c fiber.Ctx, param string) (uuid.UUID, error) {
	value := c.Params(param)
	if value == "" {
		return uuid.Nil, errs.ErrBadRequest
	}

	id, err := uuid.Parse(value)
	if err != nil {
		return uuid.Nil, errs.ErrBadRequest
	}

	return id, nil
}
