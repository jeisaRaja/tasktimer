package task

import (
	"database/sql"
	"time"
)

type TaskService struct {
	db *sql.DB
}

func NewTaskService(db *sql.DB) *TaskService {
	return &TaskService{
		db: db,
	}
}

// New creates a new Task instance with the provided name and optional configuration functions.
// The opts parameter is a variadic list of functions that modify the Task instance.
// Each function takes a pointer to the Task and applies specific changes or configurations.
func New(name string, opts ...func(*Task)) Task {
	task := Task{
		Name: name,
	}
	for _, opt := range opts {
		opt(&task)
	}
	return task
}

func WithDescription(description string) func(*Task) {
	return func(t *Task) {
		t.Description = description
	}
}

func WithRecurringDays(days []time.Weekday) func(*Task) {
	return func(t *Task) {
		t.RecurringDays = days
	}
}
