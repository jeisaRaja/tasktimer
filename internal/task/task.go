package task

import "time"

type Task struct {
	ID            int
	Name          string
	Description   string
	TimeSpent     time.Duration
	RecurringDays []time.Weekday
	Tags          []string
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
