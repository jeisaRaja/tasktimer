package main

import (
	"database/sql"
	"time"

	"github.com/jeisaRaja/tasktimer/internal/storage"
	"github.com/jeisaRaja/tasktimer/internal/task"
	"github.com/jeisaRaja/tasktimer/internal/ui"
)

type App struct {
	day   time.Weekday
	db    *sql.DB
	tasks []task.DailyTask
}

func NewApp() *App {
	return &App{
		day: time.Now().Weekday(),
	}
}

func main() {
	app := NewApp()
	app.db = storage.ConnectDB()
	ui.Start()
}
