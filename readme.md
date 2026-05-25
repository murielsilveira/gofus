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

If you already have databases and changed migrations, wipe and re-apply:

```bash
make reset-db-all
```

## Run

```bash
make run
```

Server listens on `:9000` (override with `PORT`). API lives at `/api/v1`.

## Testing

Tests use a **separate** database (`gofus_test` by default). Create it with `make create-db`. Tests truncate tables between cases — never point them at dev. After migration changes, reset the test DB with `make reset-db-test`.

```bash
make test
```

Integration tests run in-process against the full stack (HTTP → handlers → services → PostgreSQL). If Postgres is unavailable, tests are skipped.

## Makefile

| Command              | Description                             |
| -------------------- | --------------------------------------- |
| `make run`           | Start the server                        |
| `make build`         | Build binary to `bin/server`            |
| `make create-db`     | Create dev and test databases           |
| `make reset-db`      | Drop and re-apply migrations on dev DB  |
| `make reset-db-test` | Drop and re-apply migrations on test DB |
| `make reset-db-all`  | Reset both dev and test databases       |
| `make migrate-up`    | Apply migrations                        |
| `make migrate-down`  | Roll back one migration                 |
| `make sqlc`          | Regenerate Go code from SQL             |
| `make tidy`          | Tidy Go modules                         |
| `make test`          | Run integration tests                   |

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

After changing `db/schema/` or `db/queries/`, run `make sqlc`. After changing migrations, run `make migrate-up` (or `make reset-db` to wipe and re-apply from scratch).
