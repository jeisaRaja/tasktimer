package storage

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/jeisaRaja/tasktimer/internal/models"
)

func (s *Storage) UpdateGeneratedDate() error {
	today := time.Now().Truncate(24 * time.Hour)
	_, err := s.DB.Exec("INSERT OR REPLACE INTO last_generated (id, date) VALUES (1, ?)", today)
	return err
}

func (s *Storage) HasGeneratedToday() (bool, error) {
	var lastDate time.Time

	query := `
    SELECT date FROM last_generated LIMIT 1
  `
	err := s.DB.QueryRow(query).Scan(&lastDate)
	if err != nil {
		return false, err
	}

	return lastDate.Equal(time.Now().Truncate(24 * time.Hour)), nil
}

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

func (s *Storage) InsertDailyTask(task models.DailyTask) error {
	taskDate := task.Date.Format("2006-01-02")
	query := `
    INSERT INTO daily_tasks (task_id, date, daily_target, time_spent)
    VALUES (?, ?, ?, ?)
  `

	_, err := s.DB.Exec(query, task.TaskID, taskDate, task.DailyTarget.Seconds(), task.TimeSpent.Seconds())
	if err != nil {
		return err
	}

	return nil
}

// GetAllTasks retrieves all tasks from the database.
// It executes a SQL query to fetch tasks including their ID, name, description, time spent, weekly target, and tags.
// The tags are stored in the database as a JSON string, so the function unmarshals them into a slice of strings ([]string) in Go.
func (s *Storage) GetAllTasks() ([]models.Task, error) {
	var tasks []models.Task
	query := `SELECT id, name, description, time_spent, weekly_target, tags FROM tasks;`
	result, err := s.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer result.Close()

	for result.Next() {
		var task models.Task
		var tagsJSON string
		err := result.Scan(&task.ID, &task.Name, &task.Description, &task.TimeSpent, &task.WeeklyTarget, &tagsJSON)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal([]byte(tagsJSON), &task.Tags)
		if err != nil {
			return nil, fmt.Errorf("error unmarshaling tags: %w", err)
		}
		tasks = append(tasks, task)
	}
	if err := result.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}
