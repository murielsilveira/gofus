package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
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

	engine := html.New("./templates", ".go.html")
	app := fiber.New(fiber.Config{
		Views: engine,
	})

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Error opening database: %q", err)
	}

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, world!")
	})

	app.Get("/app.html", func(c *fiber.Ctx) error {
		return c.Render("app", fiber.Map{"Some": "Var"})
	})

	app.Get("/db", func(c *fiber.Ctx) error {
		if _, err := db.Exec("SELECT 1"); err != nil {
			return c.SendString(fmt.Sprintf("Error querying database: %q", err))
		} else {
			return c.SendString("WORKED!!")
		}
	})

	app.Listen(":" + port)
}
