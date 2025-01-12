-- name: CreateTask :one
INSERT INTO tasks (
  title,
  description,
  status,
  due_date
) VALUES (
  $1, $2, $3, $4
) RETURNING *;

-- name: GetTask :one
SELECT * FROM tasks
WHERE id = $1 LIMIT 1;

-- name: ListTasks :many
SELECT * FROM tasks
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: UpdateTasks :exec
UPDATE tasks
SET title = $2,
    description = $3,
    status = $4,
    due_date = $5
WHERE id = $1
RETURNING *;

-- name: DeleteTasks :exec
DELETE FROM tasks
WHERE id = $1;

