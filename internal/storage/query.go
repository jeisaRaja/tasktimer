package storage

import (
	"encoding/json"
	"fmt"

	"github.com/jeisaRaja/tasktimer/internal/models"
)

func (s *Storage) InsertTask(task models.Task) error {
	recurringDaysJSON, err := json.Marshal(task.RecurringDays)
	if err != nil {
		return fmt.Errorf("error marshalling RecurringDays: %w", err)
	}
	tagsJSON, err := json.Marshal(task.Tags)
	if err != nil {
		return fmt.Errorf("error marshalling Tags: %w", err)
	}
	query := `
		INSERT INTO tasks (name, description, time_spent, recurring_days, tags, weekly_target)
		VALUES (?, ?, ?, ?, ?, ?)
	`

	_, err = s.DB.Exec(query, task.Name, task.Description, task.TimeSpent.Seconds(), recurringDaysJSON, tagsJSON, task.WeeklyTarget.Seconds())
	if err != nil {
		return fmt.Errorf("errorr inserting task: %w", err)
	}

	return nil
}
