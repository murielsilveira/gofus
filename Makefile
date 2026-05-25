DATABASE_URL ?= postgres://postgres:postgres@localhost/gofus_dev?sslmode=disable
TEST_DATABASE_URL ?= postgres://postgres:postgres@localhost/gofus_test?sslmode=disable
PORT ?= 9000
BIN ?= bin/server

.PHONY: run build create-db migrate-up migrate-down sqlc tidy test

run:
	PORT=$(PORT) go run ./cmd/server

build:
	go build -o $(BIN) ./cmd/server

create-db:
	-createdb -U postgres -h localhost gofus_dev
	-createdb -U postgres -h localhost gofus_test

migrate-up:
	migrate -path db/migrations -database "$(DATABASE_URL)" up

migrate-down:
	migrate -path db/migrations -database "$(DATABASE_URL)" down 1

sqlc:
	sqlc generate

tidy:
	go mod tidy

test:
	TEST_DATABASE_URL=$(TEST_DATABASE_URL) go test ./...
