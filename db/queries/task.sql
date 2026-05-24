-- name: CreateTask :one
INSERT INTO tasks (column_id, title, description, position)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetTask :one
SELECT * FROM tasks
WHERE id = $1;

-- name: ListTasksByColumn :many
SELECT * FROM tasks
WHERE column_id = $1
ORDER BY position, created_at;

-- name: UpdateTask :one
UPDATE tasks
SET column_id = $2, title = $3, description = $4, position = $5, updated_at = now()
WHERE id = $1
RETURNING *;

-- name: DeleteTask :execrows
DELETE FROM tasks
WHERE id = $1;
