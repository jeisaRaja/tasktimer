package task

import (
	"fmt"
	"time"

	"github.com/jeisaRaja/tasktimer/internal/models"
	"github.com/jeisaRaja/tasktimer/internal/storage"
)

type TaskService struct {
	db *storage.Storage
}

func NewTaskService(db *storage.Storage) *TaskService {
	return &TaskService{
		db: db,
	}
}

// New creates a new Task instance with the provided name and optional configuration functions.
// The opts parameter is a variadic list of functions that modify the Task instance.
// Each function takes a pointer to the Task and applies specific changes or configurations.
func (ts *TaskService) New(task models.Task) error {
	err := ts.db.InsertTask(task)
	if err != nil {
		return fmt.Errorf("there is an error while inserting task: %w", err)
	}
	return nil
}

func WithDescription(description string) func(*models.Task) {
	return func(t *models.Task) {
		t.Description = description
	}
}

func WithRecurringDays(days []time.Weekday) func(*models.Task) {
	return func(t *models.Task) {
		t.RecurringDays = days
	}
}
