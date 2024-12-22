-- name: GetTasks :many
SELECT * FROM tasks
ORDER BY created_at;

-- name: GetTask :one
SELECT * FROM tasks
WHERE task_id = $1;

-- name: CreateTask :exec
INSERT INTO tasks (task_id, task_name) 
VALUES ($1, $2);

-- name: UpdateTaskState :exec
UPDATE tasks SET status = $2
WHERE task_id = $1;