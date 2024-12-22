-- name: SaveEvent :exec
INSERT INTO events (aggregate_id, event_type, payload) VALUES ($1,$2,$3);

-- name: GetEvents :many
SELECT * FROM events;

-- name: GetEventsAfter :many
SELECT * FROM events
WHERE id > $1;

-- name: GetCheckpoint :one
SELECT last_checkpoint FROM checkpoints
WHERE projection_id = $1;

-- name: UpdateLastCheckpoiint :exec
UPDATE checkpoints SET last_checkpoint = $2
WHERE projection_id = $1;
