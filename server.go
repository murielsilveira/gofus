package main

import (
	"os"

	"github.com/gofiber/fiber"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "9000"
	}

	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) {
		c.Send("Hello, world!")
	})

	app.Listen(port)
}
