# gofus

Kanban-style REST API (boards, columns, tasks) built with Go, Fiber, PostgreSQL, [sqlc](https://github.com/sqlc-dev/sqlc), and [golang-migrate](https://github.com/golang-migrate/migrate).

## Prerequisites

- Go 1.26+
- PostgreSQL
- CLI tools: `migrate`, `sqlc` (`go install github.com/golang-migrate/migrate/v4/cmd/migrate@latest` and `go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest`)

## Setup

```bash
make create-db
make migrate-up
make sqlc
```

## Run

```bash
make run
```

Server listens on `:9000` (override with `PORT`). API lives at `/api/v1`.

## Makefile

| Command             | Description                  |
| ------------------- | ---------------------------- |
| `make run`          | Start the server             |
| `make build`        | Build binary to `bin/server` |
| `make create-db`    | Create `gofus_dev` database  |
| `make migrate-up`   | Apply migrations             |
| `make migrate-down` | Roll back one migration      |
| `make sqlc`         | Regenerate Go code from SQL  |
| `make tidy`         | Tidy Go modules              |

## Layout

```
cmd/server/          entrypoint
internal/board/      board domain (handler, service)
internal/column/     column domain
internal/task/       task domain
internal/platform/   shared db, http, server wiring
internal/platform/server/routes.go   all API routes in one place
internal/db/sqlc/    generated query code (do not edit)
db/schema/           DDL for sqlc
db/queries/          SQL queries for sqlc
db/migrations/       schema migrations
```

After changing `db/schema/` or `db/queries/`, run `make sqlc`. After changing migrations, run `make migrate-up`.
