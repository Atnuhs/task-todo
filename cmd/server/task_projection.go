package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"tasktodo/db"
)

type TaskProjection struct {
	queries       *db.Queries
	projectrionId string
}

func NewTaskProjection(queries *db.Queries) TaskProjection {
	return TaskProjection{
		queries:       queries,
		projectrionId: "task_projection",
	}
}

func (p TaskProjection) Run(ctx context.Context) error {
	slog.InfoContext(ctx, "Start projection", slog.String("projection id", p.projectrionId))
	checkpoint, err := p.queries.GetCheckpoint(ctx, p.projectrionId)
	if err != nil {
		return fmt.Errorf("failed get checkpoint task_projection: %w", err)
	}
	slog.InfoContext(ctx, "Last checkpoint retrieved",
		slog.String("projection_id", p.projectrionId),
		slog.Int("checkpoint", int(checkpoint)),
	)

	changes, err := p.queries.GetEventsAfter(ctx, checkpoint)
	if err != nil {
		return fmt.Errorf("failed get event after %d from events: %w", checkpoint, err)
	}
	slog.InfoContext(ctx, "Changes after last checkpoint retrieved",
		slog.String("projection_id", p.projectrionId),
		slog.Int("changes num", len(changes)),
	)

	if len(changes) == 0 {
		slog.Info("No changes detected")
		return nil
	}

	for _, c := range changes {
		eventType, err := EventTypeString(c.EventType)
		if err != nil {
			return fmt.Errorf("invalid event type: %s: %w", c.EventType, err)
		}

		switch eventType {
		case TaskCreated:
			var payload TaskCreatedPayload
			if err := json.Unmarshal(c.Payload, &payload); err != nil {
				return fmt.Errorf("failed to unmarshal event payload: %w", err)
			}
			if err := p.queries.CreateTask(ctx, db.CreateTaskParams{
				TaskID:   c.AggregateID,
				TaskName: payload.TaskName,
			}); err != nil {
				slog.Error("TaskCreated event projection failed", slog.Any("err", err))
				return fmt.Errorf("failed to projection TaskCreated event: %w", err)
			}
		case TaskStarted, TaskPending, TaskCompleted, TaskCanceled:
			status := EventTypeToTaskStatus(eventType)
			if status == Unknown {
				return fmt.Errorf("unknown status %s", eventType)
			}

			if err := p.queries.UpdateTaskState(ctx, db.UpdateTaskStateParams{
				TaskID: c.AggregateID,
				Status: db.TaskStatus(status.LowerString()),
			}); err != nil {
				slog.Error(fmt.Sprintf("%s event projection failed", eventType), slog.Any("err", err))
				return fmt.Errorf("failed to projection TaskStarted event: %w", err)
			}
		default:
			return fmt.Errorf("unknown event type %s", eventType)
		}
	}

	new_checkpoint := changes[len(changes)-1].ID
	p.queries.UpdateLastCheckpoiint(ctx, db.UpdateLastCheckpoiintParams{
		ProjectionID:   p.projectrionId,
		LastCheckpoint: new_checkpoint,
	})

	return nil
}
