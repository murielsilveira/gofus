package server

import (
	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/murielsilveira/gofus/internal/board"
	"github.com/murielsilveira/gofus/internal/column"
	"github.com/murielsilveira/gofus/internal/db/sqlc"
	"github.com/murielsilveira/gofus/internal/task"
)

func NewWithPool(pool *pgxpool.Pool) *fiber.App {
	queries := sqlc.New(pool)

	return New(Config{
		Pool: pool,
		Handlers: Handlers{
			Board:  board.NewHandler(board.NewService(queries)),
			Column: column.NewHandler(column.NewService(queries)),
			Task:   task.NewHandler(task.NewService(queries)),
		},
	})
}
