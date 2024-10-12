package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/jeisaRaja/tasktimer/internal/task"
)

type Model struct {
	TaskService *task.TaskService
	EditTab     tea.Model
	createTask  TaskCreationModel
}

func newModel(ts *task.TaskService) Model {
	tc := initialTaskCreation()
	return Model{
		TaskService: ts,
		createTask:  tc,
	}
}

func (m Model) Init() tea.Cmd {
	var cmds []tea.Cmd

	cmds = append(cmds, tea.EnterAltScreen)
	cmds = append(cmds, m.createTask.Init())
	return tea.Batch(cmds...)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		}
	}
	updatedCreateTask, cmd := m.createTask.Update(msg)
	m.createTask = updatedCreateTask.(TaskCreationModel)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	return m.createTask.View()
}
