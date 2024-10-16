package task

import (
	"fmt"

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
	if err := ts.validateTask(&task); err != nil {
		return err
	}

	if task.RecurringDays != nil && len(task.RecurringDays) > 0 {
		if err := ts.handleRecurringTask(&task); err != nil {
			return fmt.Errorf("error processing recurring task: %w", err)
		}
	}

	err := ts.db.InsertTask(task)
	if err != nil {
		return fmt.Errorf("there is an error while inserting task: %w", err)
	}
	return nil
}

func (ts *TaskService) NewDailyTask(task models.DailyTask) error {
	err := ts.db.InsertDailyTask(task)
	if err != nil {
		return err
	}

	return nil
}

func (ts *TaskService) validateTask(task *models.Task) error {
	if task.Name == "" {
		return fmt.Errorf("task name cannot be empty")
	}
	if task.TimeSpent < 0 {
		return fmt.Errorf("time spent cannot be negative")
	}
	return nil
}

func (ts *TaskService) handleRecurringTask(task *models.Task) error {
	fmt.Printf("Recurring task setup for days: %v\n", task.RecurringDays)
	return nil
}

func (ts *TaskService) GetTasks(args ...string) ([]models.Task, error) {
	if len(args) == 0 {
		res, err := ts.db.GetAllTasks()
		if err != nil {
			return nil, err
		}
		return res, nil
	}

	return nil, nil
}
