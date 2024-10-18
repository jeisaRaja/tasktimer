package storage

import (
	"reflect"
	"testing"

	"github.com/jeisaRaja/tasktimer/internal/models"
)

func TestGetAllTasks(t *testing.T) {
	store := ConnectTestDB()
	if store == nil {
		t.Fatal("failed to conect to the database")
	}
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
