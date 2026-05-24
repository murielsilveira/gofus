package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/murielsilveira/gofus/internal/board"
	"github.com/murielsilveira/gofus/internal/column"
	"github.com/murielsilveira/gofus/internal/db/sqlc"
	platformdb "github.com/murielsilveira/gofus/internal/platform/db"
	"github.com/murielsilveira/gofus/internal/platform/server"
	"github.com/murielsilveira/gofus/internal/task"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "9000"
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	pool, err := platformdb.Connect(ctx, platformdb.DefaultDatabaseURL())
	if err != nil {
		log.Fatalf("database connection failed: %v", err)
	}
	defer pool.Close()

	queries := sqlc.New(pool)

	boardService := board.NewService(queries)
	columnService := column.NewService(queries)
	taskService := task.NewService(queries)

	app := server.New(server.Config{
		Pool: pool,
		Handlers: server.Handlers{
			Board:  board.NewHandler(boardService),
			Column: column.NewHandler(columnService),
			Task:   task.NewHandler(taskService),
		},
	})

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
