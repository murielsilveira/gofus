package column

import (
	"github.com/gofiber/fiber/v3"
	"github.com/murielsilveira/gofus/internal/platform/errs"
	"github.com/murielsilveira/gofus/internal/platform/httpx"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

type createRequest struct {
	Name     string `json:"name"`
	Position *int32 `json:"position"`
}

type updateRequest struct {
	Name     *string `json:"name"`
	Position *int32  `json:"position"`
}

func (h *Handler) Create(c fiber.Ctx) error {
	boardID, err := httpx.ParseUUID(c, "boardID")
	if err != nil {
		return httpx.Error(c, err)
	}

	var req createRequest
	if err := c.Bind().Body(&req); err != nil {
		return httpx.Error(c, errs.ErrBadRequest)
	}
	if req.Name == "" {
		return httpx.Error(c, errs.ErrBadRequest)
	}

	position := int32(0)
	if req.Position != nil {
		position = *req.Position
	}

	column, err := h.service.Create(c, CreateInput{
		BoardID:  boardID,
		Name:     req.Name,
		Position: position,
	})
	if err != nil {
		return httpx.Error(c, err)
	}

	return httpx.JSON(c, fiber.StatusCreated, column)
}

func (h *Handler) List(c fiber.Ctx) error {
	boardID, err := httpx.ParseUUID(c, "boardID")
	if err != nil {
		return httpx.Error(c, err)
	}

	columns, err := h.service.ListByBoard(c, boardID)
	if err != nil {
		return httpx.Error(c, err)
	}

	return httpx.JSON(c, fiber.StatusOK, columns)
}

func (h *Handler) Get(c fiber.Ctx) error {
	id, err := httpx.ParseUUID(c, "id")
	if err != nil {
		return httpx.Error(c, err)
	}

	column, err := h.service.Get(c, id)
	if err != nil {
		return httpx.Error(c, err)
	}

	return httpx.JSON(c, fiber.StatusOK, column)
}

func (h *Handler) Update(c fiber.Ctx) error {
	id, err := httpx.ParseUUID(c, "id")
	if err != nil {
		return httpx.Error(c, err)
	}

	var req updateRequest
	if err := c.Bind().Body(&req); err != nil {
		return httpx.Error(c, errs.ErrBadRequest)
	}
	if req.Name == nil && req.Position == nil {
		return httpx.Error(c, errs.ErrBadRequest)
	}

	column, err := h.service.Update(c, id, UpdateInput{
		Name:     req.Name,
		Position: req.Position,
	})
	if err != nil {
		return httpx.Error(c, err)
	}

	return httpx.JSON(c, fiber.StatusOK, column)
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
