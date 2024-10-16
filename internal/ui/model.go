package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/jeisaRaja/tasktimer/internal/task"
)

type Model struct {
	taskService *task.TaskService
	views       []tea.Model
	activeView  tea.Model
	viewIndex   int
}

func newModel(ts *task.TaskService) Model {
	var views []tea.Model
	todayTask := initialTodayTaskModel(ts)
	createTask := initialTaskCreation(ts)

	views = append(views, todayTask)
	views = append(views, createTask)

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

// Model handles keys for quit and switching between views
// The Update function also update the activeView
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case TaskUpdateMsg:
		viewModel, cmd := m.activeView.Update(msg)
		m.activeView = viewModel
		return m, cmd
	case InsertTaskMsg:
		err := m.taskService.New(msg.Task)
		if err != nil {
			panic(err)
		}
		m.activeView = m.views[0]
		if view, ok := m.activeView.(TodayTaskModel); ok {
			refreshCmd := view.Refresh()
			return m, refreshCmd
		}
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "ctrl+n":
			m.activeView = m.views[1]
			return m, nil
		case "esc":
			m.activeView = m.views[0]
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

func (m Model) SwitchView() int {
	view := m.viewIndex
	view++
	if view > len(m.views)-1 {
		view = 0
	}

	return view
}
