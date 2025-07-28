package usecase

import (
	"context"
	"time"
)

type Task struct {
	ID          int64
	Title       string
	Description string
	IsDone      string
	CreatedAt   time.Time
}

type TaskService interface {
	Create(ctx context.Context, t Task) (int64, error)
	List(ctx context.Context) ([]Task, error)
	Delete(ctx context.Context, id int64) error
	MarkDone(ctx context.Context, id int64) error
}
