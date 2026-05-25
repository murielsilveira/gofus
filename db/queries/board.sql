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
WITH deleted_tasks AS (
    DELETE FROM tasks
    WHERE column_id IN (SELECT id FROM columns WHERE board_id = $1)
),
deleted_columns AS (
    DELETE FROM columns WHERE board_id = $1
)
DELETE FROM boards
WHERE boards.id = $1;
