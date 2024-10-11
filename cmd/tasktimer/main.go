package main

import (
	"github.com/jeisaRaja/tasktimer/internal/storage"
	"github.com/jeisaRaja/tasktimer/internal/task"
	"github.com/jeisaRaja/tasktimer/internal/ui"
)

func main() {
	db := storage.ConnectDB()
	defer db.Close()
	taskService := task.NewTaskService(db)
	ui.Start(taskService)
}
