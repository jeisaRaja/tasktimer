package storage

import (
	"reflect"
	"testing"
	"time"

	"github.com/jeisaRaja/tasktimer/internal/models"
)

func TestGetAllTasks(t *testing.T) {
	store := ConnectTestDB()
	if store == nil {
		t.Fatal("failed to conect to the database")
	}
	defer store.Close()

	tasks := []models.Task{
		{
			Name:         "task1",
			Description:  "task1 description",
			TimeSpent:    120,
			WeeklyTarget: 600,
			Tags:         []string{"work", "urgent"},
		},
		{
			Name:         "task2",
			Description:  "task2 description",
			TimeSpent:    60,
			WeeklyTarget: 300,
			Tags:         []string{"personal"},
		},
	}
	err := insertTasks(store, tasks)
	if err != nil {
		t.Fatalf("failed to insert task: %v", err)
	}

	retreivedTasks, err := store.GetAllTasks()
	if err != nil {
		t.Fatalf("failed to retreive tasks: %v", err)
	}

	for i, task := range tasks {
		expected := task
		actual := retreivedTasks[i]

		if expected.Name != actual.Name {
			t.Errorf("Task name mismatch: expected %s, got %s", expected.Name, actual.Name)
		}

		if expected.Description != actual.Description {
			t.Errorf("Task description mismatch: expected %s, got %s", expected.Description, actual.Description)
		}

		if expected.TimeSpent != actual.TimeSpent {
			t.Errorf("Task time spent mismatch: expected %v, got %v", expected.TimeSpent, actual.TimeSpent)
		}

		if expected.WeeklyTarget != actual.WeeklyTarget {
			t.Errorf("Task weekly target mismatch: expected %d, got %d", expected.WeeklyTarget, actual.WeeklyTarget)
		}

		if !reflect.DeepEqual(expected.Tags, actual.Tags) {
			t.Errorf("Task tags mismatch: expected %+v, got %+v", expected.Tags, actual.Tags)
		}

	}
}

func insertTasks(db *Storage, tasks []models.Task) error {
	for _, task := range tasks {
		err := db.InsertTask(task)
		if err != nil {
			return err
		}
	}
	return nil
}

func TestInsertDailyTask(t *testing.T) {
	store := ConnectTestDB()
	defer store.Close()
	if store == nil {
		t.Fatal("failed to connect to db")
	}

	task := models.DailyTask{
		TaskID:      1,
		Date:        time.Now(),
		DailyTarget: time.Hour * 5,
		TimeSpent:   time.Hour * 3,
	}

	err := store.InsertDailyTask(task)
	if err != nil {
		t.Fatalf("failed to insert dailytask: %v", err)
	}

	query := `SELECT task_id, date, daily_target, time_spent FROM daily_tasks WHERE task_id = ? LIMIT 1`
	var retrievedTaskID int
	var retrievedDate string
	var retrievedDailyTarget, retrievedTimeSpent int64
	err = store.DB.QueryRow(query, task.TaskID).Scan(&retrievedTaskID, &retrievedDate, &retrievedDailyTarget, &retrievedTimeSpent)
	if err != nil {
		t.Fatalf("failed to retrieve daily task: %v", err)
	}

	retrievedDateParsed, err := time.Parse(time.RFC3339, retrievedDate)
	if err != nil {
		t.Fatalf("failed to parse date: %v", err)
	}

	if retrievedTaskID != task.TaskID {
		t.Errorf("expected task ID %d, got %d", task.TaskID, retrievedTaskID)
	}

	if !retrievedDateParsed.Equal(task.Date.Truncate(24 * time.Hour)) {
		t.Errorf("expected date %v, got %v", task.Date.Truncate(24*time.Hour), retrievedDateParsed)
	}

	if retrievedDailyTarget != int64(task.DailyTarget.Nanoseconds()) {
		t.Errorf("expected daily target %v, got %v", int64(task.DailyTarget.Nanoseconds()), retrievedDailyTarget)
	}

	if retrievedTimeSpent != int64(task.TimeSpent.Nanoseconds()) {
		t.Errorf("expected time spent %v, got %v", int64(task.TimeSpent.Nanoseconds()), retrievedTimeSpent)
	}
}
