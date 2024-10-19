package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/jeisaRaja/tasktimer/internal/models"
)

type TaskSelector struct {
	list []models.Task
	curr int
}

func initialTaskSelector() TaskSelector {
	return TaskSelector{
		curr: 0,
		list: []models.Task{},
	}
}

func (ts *TaskSelector) SetTasks(tasks []models.Task) {
	ts.list = tasks
	ts.curr = 0
}

func (m TaskSelector) Init() tea.Cmd {
	return nil
}

func (m TaskSelector) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			return m, func() tea.Msg {
				return TaskSelectMsg(m.list[m.curr])
			}
		case "k":
			m.curr -= 1
			if m.curr < 0 {
				m.curr = 0
			}
		case "j":
			m.curr += 1
			if m.curr >= len(m.list) {
				m.curr = len(m.list) - 1
			}
		}
	}
	return m, nil
}

func (m TaskSelector) View() string {
	if len(m.list) == 0 {
		return "No tasks available"
	}
	s := "Select a task to add:\n\n"
	for i, task := range m.list {
		cursor := " "
		if m.curr == i {
			cursor = ">"
		}
		s += cursor + " " + task.Name + "\n"
	}
	return s
}

func (m TaskSelector) WithTasks(tasks []models.Task) TaskSelector {
	m.list = tasks
	return m
}
