package models

import "time"

type TaskWithDaily struct {
	Task
	DailyTask
}

type Task struct {
	ID            int
	Name          string
	Description   string
	TimeSpent     time.Duration
	RecurringDays []time.Weekday
	Tags          []string
	WeeklyTarget  time.Duration
}

type DailyTask struct {
	TaskID      int
	Date        time.Time
	DailyTarget time.Duration
	TimeSpent   time.Duration
}

type Schedule struct {
	Day        time.Weekday
	DailyTasks []DailyTask
}

type Summary struct {
	WeeklySpent  map[string]time.Duration
	MonthlySpent map[string]time.Duration
	YearlySpent  map[string]time.Duration
}
