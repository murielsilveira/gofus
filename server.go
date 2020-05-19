package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber"
	_ "github.com/lib/pq"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "9000"
	}
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:postgres@localhost/?sslmode=disable"
	}

	app := fiber.New()

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Error opening database: %q", err)
	}

	app.Get("/", func(c *fiber.Ctx) {
		c.Send("Hello, world!")
	})

	app.Get("/db", func(c *fiber.Ctx) {
		if _, err := db.Exec("SELECT 1"); err != nil {
			c.Send(fmt.Sprintf("Error querying database: %q", err))
		} else {
			c.Send("WORKED!!")
		}
	})

	app.Listen(port)
}
