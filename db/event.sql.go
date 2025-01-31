// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: event.sql

package db

import (
	"context"

	"github.com/google/uuid"
)

const getCheckpoint = `-- name: GetCheckpoint :one
SELECT last_checkpoint FROM checkpoints
WHERE projection_id = $1
`

func (q *Queries) GetCheckpoint(ctx context.Context, projectionID string) (int32, error) {
	row := q.db.QueryRow(ctx, getCheckpoint, projectionID)
	var last_checkpoint int32
	err := row.Scan(&last_checkpoint)
	return last_checkpoint, err
}

const getEvents = `-- name: GetEvents :many
SELECT id, aggregate_id, event_type, payload, created_at FROM events
`

func (q *Queries) GetEvents(ctx context.Context) ([]Event, error) {
	rows, err := q.db.Query(ctx, getEvents)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Event
	for rows.Next() {
		var i Event
		if err := rows.Scan(
			&i.ID,
			&i.AggregateID,
			&i.EventType,
			&i.Payload,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getEventsAfter = `-- name: GetEventsAfter :many
SELECT id, aggregate_id, event_type, payload, created_at FROM events
WHERE id > $1
`

func (q *Queries) GetEventsAfter(ctx context.Context, id int32) ([]Event, error) {
	rows, err := q.db.Query(ctx, getEventsAfter, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Event
	for rows.Next() {
		var i Event
		if err := rows.Scan(
			&i.ID,
			&i.AggregateID,
			&i.EventType,
			&i.Payload,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const saveEvent = `-- name: SaveEvent :exec
INSERT INTO events (aggregate_id, event_type, payload) VALUES ($1,$2,$3)
`

type SaveEventParams struct {
	AggregateID uuid.UUID
	EventType   string
	Payload     []byte
}

func (q *Queries) SaveEvent(ctx context.Context, arg SaveEventParams) error {
	_, err := q.db.Exec(ctx, saveEvent, arg.AggregateID, arg.EventType, arg.Payload)
	return err
}

const updateLastCheckpoiint = `-- name: UpdateLastCheckpoiint :exec
UPDATE checkpoints SET last_checkpoint = $2
WHERE projection_id = $1
`

type UpdateLastCheckpoiintParams struct {
	ProjectionID   string
	LastCheckpoint int32
}

func (q *Queries) UpdateLastCheckpoiint(ctx context.Context, arg UpdateLastCheckpoiintParams) error {
	_, err := q.db.Exec(ctx, updateLastCheckpoiint, arg.ProjectionID, arg.LastCheckpoint)
	return err
}
