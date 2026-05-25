package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/murielsilveira/gofus/internal/platform/db"
	"github.com/murielsilveira/gofus/internal/platform/server"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "9000"
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	pool, err := db.Connect(ctx, db.DatabaseURL())
	if err != nil {
		log.Fatalf("database connection failed: %v", err)
	}
	defer pool.Close()

	app := server.NewWithPool(pool)

	go func() {
		<-ctx.Done()
		if err := app.Shutdown(); err != nil {
			log.Printf("server shutdown error: %v", err)
		}
	}()

	log.Printf("listening on :%s", port)
	if err := app.Listen(":" + port); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
