package server

import (
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/template/html/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Config struct {
	Pool     *pgxpool.Pool
	Handlers Handlers
}

func New(cfg Config) *fiber.App {
	engine := html.New("./templates", ".go.html")
	app := fiber.New(fiber.Config{
		Views: engine,
	})

	registerRoutes(app, cfg.Pool, cfg.Handlers)

	return app
}
