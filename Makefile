DATABASE_URL ?= postgres://postgres:postgres@localhost/gofus_dev?sslmode=disable
PORT ?= 9000
BIN ?= bin/server

.PHONY: run build create-db migrate-up migrate-down sqlc tidy

run:
	PORT=$(PORT) go run ./cmd/server

build:
	go build -o $(BIN) ./cmd/server

create-db:
	-createdb -U postgres -h localhost gofus_dev

migrate-up:
	migrate -path db/migrations -database "$(DATABASE_URL)" up

migrate-down:
	migrate -path db/migrations -database "$(DATABASE_URL)" down 1

sqlc:
	sqlc generate

tidy:
	go mod tidy
