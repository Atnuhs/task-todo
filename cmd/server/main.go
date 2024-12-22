package main

import (
	"context"
	"encoding/json"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"tasktodo/db"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humago"
	_ "github.com/danielgtaylor/huma/v2/formats/cbor"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type Task struct {
	Name   string    `json:"name" example:"空き缶を片付ける" doc:"task name"`
	Id     uuid.UUID `json:"task_id" example:"ea087856-822f-4db3-a9c7-23ce3dd337b3" doc:"uuid"`
	Status string    `json:"status" enum:"pending,doing,completed,cancelled" example:"pending" doc:"task status"`
}

type Status struct {
	Status string `json:"status" enum:"pending,doing,completed,cancelled" example:"pending" doc:"task status"`
}

type ResponseTasks struct {
	Body []Task
}

func GetTasksHandler(queries *db.Queries) func(ctx context.Context, _ *struct{}) (*ResponseTasks, error) {
	return func(ctx context.Context, i *struct{}) (*ResponseTasks, error) {
		tasks, err := queries.GetTasks(ctx)
		if err != nil {
			slog.Error("Error querying GetTasks", slog.Any("err", err))
			return nil, err
		}

		respBody := make([]Task, len(tasks))
		for i, task := range tasks {
			respBody[i] = Task{
				Name:   task.TaskName,
				Id:     task.TaskID,
				Status: string(task.Status),
			}
		}

		slog.Info("Tasks retrieved", slog.Any("tasks", respBody))
		return &ResponseTasks{
			Body: respBody,
		}, nil
	}
}

type (
	RequestCreateTask struct {
		Body struct {
			Name string `json:"name" example:"空き缶を片付ける" doc:"タスクの名前"`
		}
	}
	ResponseCreateTask struct {
		Body Task
	}
)

func CreateTask(queries *db.Queries) func(ctx context.Context, req *RequestCreateTask) (*ResponseCreateTask, error) {
	return func(ctx context.Context, req *RequestCreateTask) (*ResponseCreateTask, error) {
		taskId := uuid.New()
		saveEvent := TaskCreatedEvent(taskId, req.Body.Name)
		param, err := saveEvent.ToSaveEventParam()
		if err != nil {
			slog.Error("Error saveEvent to saveParam", slog.Any("err", err))
			return nil, huma.Error500InternalServerError("internal server error")
		}

		if err := queries.SaveEvent(ctx, param); err != nil {
			slog.Error("Error querying SaveEvent", slog.Any("err", err))
			return nil, huma.Error500InternalServerError("internal server error")
		}
		slog.Info("Event saved", slog.Any("param", saveEvent.Payload))

		p := NewTaskProjection(queries)
		if err := p.Run(ctx); err != nil {
			slog.Error("Error TaskProjection of CreateTask", slog.Any("err", err))
			return nil, huma.Error500InternalServerError("internal server error")
		}
		slog.Info("Task Projection excuted")

		return &ResponseCreateTask{
			Body: Task{
				Name:   req.Body.Name,
				Id:     taskId,
				Status: "pending",
			},
		}, nil
	}
}

func GetEventsHandler(queries *db.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.Background()
		events, err := queries.GetEvents(ctx)
		if err != nil {
			slog.Error("Error querying GetEvents", slog.Any("err", err))
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		slog.Info("Events length", slog.Int("length", len(events)))
		slog.Info("Events retrieved", slog.Any("events", events))
		w.Header().Set("Content-type", "application/json")
		json.NewEncoder(w).Encode(events)
	}
}

type RequestUpdateTask struct {
	TaskId uuid.UUID `path:"taskId" doc:"taskId" format:"uuid"`
	Body   Status
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		slog.Info(
			"Request started",
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
			slog.String("remote_addr", r.RemoteAddr),
		)
		next.ServeHTTP(w, r)
		slog.Info(
			"Request completed",
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
			slog.String("remote_addr", r.RemoteAddr),
			slog.Duration("duration", time.Since(start)),
		)
	})
}

func main() {
	// logging
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	// db
	ctx := context.Background()
	connString := "postgres://root:password@localhost:5432/appdb?sslmode=disable"
	conn, err := pgx.Connect(ctx, connString)
	if err != nil {
		slog.Error("postgres connection error", slog.Any("err", err))
	}
	defer conn.Close(ctx)

	queries := db.New(conn)

	// router
	mux := http.NewServeMux()
	api := humago.New(mux, huma.DefaultConfig("TaskTODO", "1.0.0"))
	huma.Get(api, "/tasks", GetTasksHandler(queries))
	huma.Post(api, "/tasks", CreateTask(queries))
	huma.Patch(api, "/tasks/{taskId}", func(ctx context.Context, req *RequestUpdateTask) (*struct{}, error) {
		slog.Info("UpdateTask started", slog.String("taskId", req.TaskId.String()), slog.String("status", req.Body.Status))
		status, err := TaskStatusString(req.Body.Status)
		if err != nil {
			slog.Error("Error validating task status", slog.Any("err", err))
			return nil, huma.Error500InternalServerError("internal server error")
		}

		event := TaskUpdatedEvent(req.TaskId, status)
		param, err := event.ToSaveEventParam()
		if err != nil {
			slog.Error("Error startEvent to param", slog.Any("err", err))
			return nil, huma.Error500InternalServerError("internal server error")
		}
		if err := queries.SaveEvent(ctx, param); err != nil {
			slog.Error("Error querying SaveEvent", slog.Any("err", err))
			return nil, huma.Error500InternalServerError("internal server error")
		}
		slog.Info("Event saved", slog.Any("param", event.Payload))

		p := NewTaskProjection(queries)
		if err := p.Run(ctx); err != nil {
			slog.Error("Error TaskProjection of UpdateTask", slog.Any("err", err))
			return nil, huma.Error500InternalServerError("internal server error")
		}
		slog.Info("Task Projection excuted")
		return nil, nil
	})
	mux.HandleFunc("GET /events", GetEventsHandler(queries))
	mux.Handle("GET /", http.FileServer(http.Dir("./static")))

	port := ":3000"
	log.Printf("Server running at http://localhost%s\n", port)

	err = http.ListenAndServe(port, LoggingMiddleware(mux))
	if err != nil {
		log.Fatal(err)
	}
}
