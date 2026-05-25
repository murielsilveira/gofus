package httpx

import (
	"errors"
	"log"
	"strconv"

	"github.com/gofiber/fiber/v3"

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

func ParseID(c fiber.Ctx, param string) (int32, error) {
	n, err := strconv.ParseInt(c.Params(param), 10, 32)
	if err != nil || n <= 0 {
		return 0, errs.ErrBadRequest
	}
	return int32(n), nil
}
