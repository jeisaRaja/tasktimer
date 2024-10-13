package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/jeisaRaja/tasktimer/internal/task"
)

type Model struct {
	taskService *task.TaskService
	views       []tea.Model
	activeView  tea.Model
}

func newModel(ts *task.TaskService) Model {
	var views []tea.Model
	createTask := initialTaskCreation()
	todayTask := initialTodayTaskModel()

	views = append(views, createTask)
	views = append(views, todayTask)

	return Model{
		taskService: ts,
		views:       views,
		activeView:  todayTask,
	}
}

func (m Model) Init() tea.Cmd {
	var cmds []tea.Cmd

	cmds = append(cmds, tea.EnterAltScreen)
	for _, view := range m.views {
		cmds = append(cmds, view.Init())
	}
	return tea.Batch(cmds...)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		}
	}

	activeView, cmd := m.activeView.Update(msg)
	m.activeView = activeView
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	return m.activeView.View()
}
