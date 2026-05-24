package board

import (
	"github.com/gofiber/fiber/v3"
	"github.com/murielsilveira/gofus/internal/platform/httpx"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

type createRequest struct {
	Name string `json:"name"`
}

type updateRequest struct {
	Name *string `json:"name"`
}

func (h *Handler) Create(c fiber.Ctx) error {
	var req createRequest
	if err := c.Bind().Body(&req); err != nil {
		return httpx.Error(c, httpx.ErrBadRequest)
	}
	if req.Name == "" {
		return httpx.Error(c, httpx.ErrBadRequest)
	}

	board, err := h.service.Create(c, req.Name)
	if err != nil {
		return httpx.Error(c, err)
	}

	return httpx.JSON(c, fiber.StatusCreated, board)
}

func (h *Handler) List(c fiber.Ctx) error {
	boards, err := h.service.List(c)
	if err != nil {
		return httpx.Error(c, err)
	}

	return httpx.JSON(c, fiber.StatusOK, boards)
}

func (h *Handler) Get(c fiber.Ctx) error {
	id, err := httpx.ParseUUID(c, "id")
	if err != nil {
		return httpx.Error(c, err)
	}

	board, err := h.service.Get(c, id)
	if err != nil {
		return httpx.Error(c, err)
	}

	return httpx.JSON(c, fiber.StatusOK, board)
}

func (h *Handler) Update(c fiber.Ctx) error {
	id, err := httpx.ParseUUID(c, "id")
	if err != nil {
		return httpx.Error(c, err)
	}

	var req updateRequest
	if err := c.Bind().Body(&req); err != nil {
		return httpx.Error(c, httpx.ErrBadRequest)
	}
	if req.Name == nil {
		return httpx.Error(c, httpx.ErrBadRequest)
	}

	board, err := h.service.Update(c, id, UpdateInput{Name: req.Name})
	if err != nil {
		return httpx.Error(c, err)
	}

	return httpx.JSON(c, fiber.StatusOK, board)
}

func (h *Handler) Delete(c fiber.Ctx) error {
	id, err := httpx.ParseUUID(c, "id")
	if err != nil {
		return httpx.Error(c, err)
	}

	if err := h.service.Delete(c, id); err != nil {
		return httpx.Error(c, err)
	}

	return c.SendStatus(fiber.StatusNoContent)
}
