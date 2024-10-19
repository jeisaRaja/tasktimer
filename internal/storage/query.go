package storage

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/jeisaRaja/tasktimer/internal/models"
)

// UpdateGeneratedDate sets the 'date' in 'last_generated' to today's date, replacing any existing record with 'id' 1.
// Returns an error if the operation fails.
func (s *Storage) UpdateGeneratedDate() error {
	today := time.Now().Truncate(24 * time.Hour)
	_, err := s.DB.Exec("INSERT OR REPLACE INTO last_generated (id, date) VALUES (1, ?)", today)
	return err
}

// HasGeneratedToday checks if the 'last_generated' table has a record with today's date.
// It returns true if the process was generated today, false otherwise, along with any query error.
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
	intTimeSpent := int64(task.TimeSpent.Nanoseconds())
	intWeeklyTarget := int64(task.WeeklyTarget.Nanoseconds())

	_, err = s.DB.Exec(query, task.Name, task.Description, intTimeSpent, recurringDaysJSON, tagsJSON, intWeeklyTarget)
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
	intDailyTarget := int64(task.DailyTarget.Nanoseconds())
	intTimeSpent := int64(task.TimeSpent.Nanoseconds())

	_, err := s.DB.Exec(query, task.TaskID, taskDate, intDailyTarget, intTimeSpent)
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
			return nil, fmt.Errorf("Scan failed, %v, task: %v", err, task)
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

func (s *Storage) GetTodayTasks() ([]models.TaskWithDaily, error) {
	var tasks []models.TaskWithDaily

	today := time.Now().Truncate(time.Hour * 24)

	query := `
        SELECT 
            t.ID, t.Name, t.Description, t.TimeSpent, t.RecurringDays, t.Tags, t.WeeklyTarget,
            dt.TaskID, dt.Date, dt.DailyTarget, dt.TimeSpent
        FROM tasks t
        JOIN daily_tasks dt ON t.ID = dt.TaskID
        WHERE dt.Date = $1
    `
	rows, err := s.DB.Query(query, today)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var taskWithDaily models.TaskWithDaily
		var task models.Task
		var dailyTask models.DailyTask

		err := rows.Scan(
			&task.ID, &task.Name, &task.Description, &task.TimeSpent, &task.RecurringDays, &task.Tags, &task.WeeklyTarget,
			&dailyTask.TaskID, &dailyTask.Date, &dailyTask.DailyTarget, &dailyTask.TimeSpent,
		)
		if err != nil {
			return nil, err
		}

		taskWithDaily.Task = task
		taskWithDaily.DailyTask = dailyTask

		tasks = append(tasks, taskWithDaily)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}
