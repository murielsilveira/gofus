-- name: CreateColumn :one
INSERT INTO columns (board_id, name, position)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetColumn :one
SELECT * FROM columns
WHERE id = $1;

-- name: ListColumnsByBoard :many
SELECT * FROM columns
WHERE board_id = $1
ORDER BY position, created_at;

-- name: UpdateColumn :one
UPDATE columns
SET name = $2, position = $3
WHERE id = $1
RETURNING *;

-- name: DeleteColumn :execrows
WITH deleted_tasks AS (
    DELETE FROM tasks WHERE column_id = $1
)
DELETE FROM columns
WHERE columns.id = $1;
