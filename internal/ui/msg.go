package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/jeisaRaja/tasktimer/internal/models"
)

func createInsertTaskMsg(task models.Task) tea.Cmd {
	return func() tea.Msg {
		return InsertTaskMsg{Task: task}
	}
}

type InsertTaskMsg struct {
	Task models.Task
}

type FetchTasksMsg struct {
	Tasks []models.Task
}

type TaskUpdateMsg []models.TaskWithDaily

type TaskSelectMsg models.Task
