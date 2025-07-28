package dto

import "time"

type Task struct {
	ID          int
	Title       string
	Description string
	IsDone      string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
