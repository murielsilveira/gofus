-- name: CreateBoard :one
INSERT INTO boards (name)
VALUES ($1)
RETURNING *;

-- name: GetBoard :one
SELECT * FROM boards
WHERE id = $1;

-- name: ListBoards :many
SELECT * FROM boards
ORDER BY created_at;

-- name: UpdateBoard :one
UPDATE boards
SET name = $2, updated_at = now()
WHERE id = $1
RETURNING *;

-- name: DeleteBoard :execrows
DELETE FROM boards
WHERE id = $1;
