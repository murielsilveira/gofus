package task

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
	Title       string `json:"title"`
	Description string `json:"description"`
	Position    *int32 `json:"position"`
}

type updateRequest struct {
	ColumnID    *int32  `json:"column_id"`
	Title       *string    `json:"title"`
	Description *string    `json:"description"`
	Position    *int32     `json:"position"`
}

func (h *Handler) Create(c fiber.Ctx) error {
	columnID, err := httpx.ParseID(c, "columnID")
	if err != nil {
		return httpx.Error(c, err)
	}

	var req createRequest
	if err := c.Bind().Body(&req); err != nil {
		return httpx.Error(c, errs.ErrBadRequest)
	}
	if req.Title == "" {
		return httpx.Error(c, errs.ErrBadRequest)
	}

	position := int32(0)
	if req.Position != nil {
		position = *req.Position
	}

	task, err := h.service.Create(c, CreateInput{
		ColumnID:    columnID,
		Title:       req.Title,
		Description: req.Description,
		Position:    position,
	})
	if err != nil {
		return httpx.Error(c, err)
	}

	return httpx.JSON(c, fiber.StatusCreated, task)
}

func (h *Handler) List(c fiber.Ctx) error {
	columnID, err := httpx.ParseID(c, "columnID")
	if err != nil {
		return httpx.Error(c, err)
	}

	tasks, err := h.service.ListByColumn(c, columnID)
	if err != nil {
		return httpx.Error(c, err)
	}

	return httpx.JSON(c, fiber.StatusOK, tasks)
}

func (h *Handler) Get(c fiber.Ctx) error {
	id, err := httpx.ParseID(c, "id")
	if err != nil {
		return httpx.Error(c, err)
	}

	task, err := h.service.Get(c, id)
	if err != nil {
		return httpx.Error(c, err)
	}

	return httpx.JSON(c, fiber.StatusOK, task)
}

func (h *Handler) Update(c fiber.Ctx) error {
	id, err := httpx.ParseID(c, "id")
	if err != nil {
		return httpx.Error(c, err)
	}

	var req updateRequest
	if err := c.Bind().Body(&req); err != nil {
		return httpx.Error(c, errs.ErrBadRequest)
	}
	if req.ColumnID == nil && req.Title == nil && req.Description == nil && req.Position == nil {
		return httpx.Error(c, errs.ErrBadRequest)
	}

	task, err := h.service.Update(c, id, UpdateInput{
		ColumnID:    req.ColumnID,
		Title:       req.Title,
		Description: req.Description,
		Position:    req.Position,
	})
	if err != nil {
		return httpx.Error(c, err)
	}

	return httpx.JSON(c, fiber.StatusOK, task)
}

func (h *Handler) Delete(c fiber.Ctx) error {
	id, err := httpx.ParseID(c, "id")
	if err != nil {
		return httpx.Error(c, err)
	}

	if err := h.service.Delete(c, id); err != nil {
		return httpx.Error(c, err)
	}

	return c.SendStatus(fiber.StatusNoContent)
}
