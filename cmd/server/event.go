//go:generate go install github.com/dmarkham/enumer@latest
//go:generate enumer -type=EventType
package main

import (
	"encoding/json"
	"fmt"

	"tasktodo/db"

	"github.com/google/uuid"
)

type EventType int

const (
	TaskCreated EventType = iota
	TaskPending
	TaskStarted
	TaskCompleted
	TaskCanceled
	TaskUnknown
)

type EventPayload interface {
	ToJSON() ([]byte, error)
}

type Event struct {
	AggregateID uuid.UUID
	EventType   EventType
	Payload     EventPayload
}

func (e *Event) ToSaveEventParam() (db.SaveEventParams, error) {
	payload, err := e.Payload.ToJSON()
	if err != nil {
		return db.SaveEventParams{}, fmt.Errorf("failed to MarshalJSON Event.Payload: %w", err)
	}

	return db.SaveEventParams{
		AggregateID: e.AggregateID,
		EventType:   e.EventType.String(),
		Payload:     payload,
	}, nil
}

type TaskCreatedPayload struct {
	TaskName string `json:"TaskName"`
}

func (p TaskCreatedPayload) ToJSON() ([]byte, error) {
	return json.Marshal(p)
}

func TaskCreatedEvent(taskId uuid.UUID, taskName string) *Event {
	return &Event{
		AggregateID: taskId,
		EventType:   TaskCreated,
		Payload:     TaskCreatedPayload{TaskName: taskName},
	}
}

type TaskUpdatedPayload struct{}

func (p TaskUpdatedPayload) ToJSON() ([]byte, error) {
	return json.Marshal("")
}

func TaskUpdatedEvent(taskId uuid.UUID, status TaskStatus) *Event {
	eventType := TaskStatusToEventType(status)
	if eventType == TaskUnknown {
		panic("unknown even type")
	}
	return &Event{
		AggregateID: taskId,
		EventType:   eventType,
		Payload:     TaskUpdatedPayload{},
	}
}

func TaskStatusToEventType(status TaskStatus) EventType {
	switch status {
	case Pending:
		return TaskPending
	case Doing:
		return TaskStarted
	case Completed:
		return TaskCompleted
	case Cancelled:
		return TaskCanceled
	}
	return TaskUnknown
}

func EventTypeToTaskStatus(eventType EventType) TaskStatus {
	switch eventType {
	case TaskPending:
		return Pending
	case TaskStarted:
		return Doing
	case TaskCompleted:
		return Completed
	case TaskCanceled:
		return Cancelled
	}
	return Unknown
}
