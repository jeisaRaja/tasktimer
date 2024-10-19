package ui

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/jeisaRaja/tasktimer/internal/models"
	"github.com/jeisaRaja/tasktimer/internal/task"
)

type Model struct {
	taskService *task.TaskService
	views       []tea.Model
	activeView  tea.Model
}

func newModel(ts *task.TaskService) Model {
	views := make([]tea.Model, 3)
	tasksToday, err := ts.GetTodayTasks()
	if err != nil {
		panic(fmt.Sprintf("something went wrong in newModel: %v", err))
	}
	todayTask := initialTodayTaskModel(tasksToday)
	createTask := initialTaskCreation()
	taskSelector := initialTaskSelector()

	views[viewTodayTask] = todayTask
	views[viewCreateTask] = createTask
	views[viewTaskSelector] = taskSelector

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
	case TaskSelectMsg:
		err := m.handleTaskSelect(msg)
		if err != nil {
			panic(err)
		}
	case TaskUpdateMsg:
		viewModel, cmd := m.activeView.Update(msg)
		m.activeView = viewModel
		return m, cmd
	case InsertTaskMsg:
		err := m.taskService.NewTask(msg.Task)
		if err != nil {
			panic(err)
		}
		m.activeView = m.views[viewTodayTask]
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+a":
			tasks, err := m.taskService.GetTasks()
			if err != nil {
				panic(err)
			}
			taskSelector := m.views[viewTaskSelector].(TaskSelector)
			taskSelector.SetTasks(tasks)
			m.views[viewTaskSelector] = taskSelector
			m.activeView = m.views[viewTaskSelector]
			return m, nil
		case "ctrl+c":
			return m, tea.Quit
		case "ctrl+n":
			m.activeView = m.views[viewCreateTask]
			return m, nil
		case "esc":
			m.activeView = m.views[viewTodayTask]
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

func (m *Model) handleTaskSelect(msg TaskSelectMsg) error {
	taskWithDaily := models.TaskWithDaily{
		Task: models.Task{
			ID:            msg.ID,
			Name:          msg.Name,
			Description:   msg.Description,
			RecurringDays: msg.RecurringDays,
			Tags:          msg.Tags,
			WeeklyTarget:  msg.WeeklyTarget,
		},
		DailyTask: models.DailyTask{
			TaskID:      msg.ID,
			Date:        time.Now(),
			DailyTarget: time.Hour,
			TimeSpent:   0,
		},
	}

	if err := m.taskService.InsertDailyTask(taskWithDaily.DailyTask); err != nil {
		return fmt.Errorf("error in model.go: %v", err)
	}

	todayTaskModel := m.views[viewTodayTask].(TodayTaskModel)
	todayTaskModel = todayTaskModel.AppendTask(taskWithDaily)
	m.views[viewTodayTask] = todayTaskModel

	m.activeView = m.views[viewTodayTask]

	return nil
}
