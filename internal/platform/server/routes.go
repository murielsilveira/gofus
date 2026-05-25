package server

import (
	"fmt"

	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/murielsilveira/gofus/internal/board"
	"github.com/murielsilveira/gofus/internal/column"
	"github.com/murielsilveira/gofus/internal/task"
)

type Handlers struct {
	Board  *board.Handler
	Column *column.Handler
	Task   *task.Handler
}

func registerRoutes(app *fiber.App, pool *pgxpool.Pool, h Handlers) {
	// Demo
	app.Get("/", func(c fiber.Ctx) error {
		return c.SendString("Hello, world!")
	})
	app.Get("/app.html", func(c fiber.Ctx) error {
		return c.Render("app", fiber.Map{"Some": "Var"})
	})
	app.Get("/db", func(c fiber.Ctx) error {
		if err := pool.Ping(c); err != nil {
			return c.SendString(fmt.Sprintf("Error querying database: %q", err))
		}
		return c.SendString("WORKED!!")
	})

	// Boards
	app.Post("/api/v1/boards", h.Board.Create)
	app.Get("/api/v1/boards", h.Board.List)
	app.Get("/api/v1/boards/:id", h.Board.Get)
	app.Patch("/api/v1/boards/:id", h.Board.Update)
	app.Delete("/api/v1/boards/:id", h.Board.Delete)

	// Columns
	app.Post("/api/v1/boards/:boardID/columns", h.Column.Create)
	app.Get("/api/v1/boards/:boardID/columns", h.Column.List)
	app.Get("/api/v1/columns/:id", h.Column.Get)
	app.Patch("/api/v1/columns/:id", h.Column.Update)
	app.Delete("/api/v1/columns/:id", h.Column.Delete)

	// Tasks
	app.Post("/api/v1/columns/:columnID/tasks", h.Task.Create)
	app.Get("/api/v1/columns/:columnID/tasks", h.Task.List)
	app.Get("/api/v1/tasks/:id", h.Task.Get)
	app.Patch("/api/v1/tasks/:id", h.Task.Update)
	app.Delete("/api/v1/tasks/:id", h.Task.Delete)
}
