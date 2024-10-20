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
	views := make([]tea.Model, 4)
	tasksToday, err := ts.GetTodayTasks()
	// panic(fmt.Sprintf("taskwithdaily: %v", tasksToday))
	if err != nil {
		panic(fmt.Sprintf("something went wrong in newModel: %v", err))
	}
	todayTask := initialTodayTaskModel(tasksToday)
	createTask := initialTaskCreation()
	taskSelector := initialTaskSelector()
	dailyTargetInput := initialDailyTargetInput()

	views[viewTodayTask] = todayTask
	views[viewCreateTask] = createTask
	views[viewTaskSelector] = taskSelector
	views[viewDailyTargetInput] = dailyTargetInput

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
	case DailyTargetInputMsg:
		m.activeView = m.views[viewDailyTargetInput]
		m.activeView, _ = m.activeView.Update(msg)
		return m, nil
	case TaskSelectMsg:
		m, err := m.handleTaskSelect(msg)
		if err != nil {
			panic(err)
		}
		return m, nil
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
			allTasks, err := m.taskService.GetTasks()
			todayTasks, err := m.taskService.GetTodayTasks()
			if err != nil {
				panic(err)
			}
			todayTaskIDs := make(map[int]bool)
			for _, todayTask := range todayTasks {
				todayTaskIDs[todayTask.ID] = true
			}
			var unassignedTasks []models.Task
			for _, task := range allTasks {
				if !todayTaskIDs[task.ID] {
					unassignedTasks = append(unassignedTasks, task)
				}
			}
			taskSelector := m.views[viewTaskSelector].(TaskSelector)
			taskSelector.SetTasks(unassignedTasks)
			m.views[viewTaskSelector] = taskSelector
			m.activeView = m.views[viewTaskSelector]
			return m, nil
		case "ctrl+c":
			return m, tea.Quit
		case "ctrl+n":
			m.activeView = m.views[viewCreateTask]
			if creationView, ok := m.activeView.(TaskCreationModel); ok {
				m.activeView = creationView.Clear()
			}
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

func (m *Model) handleTaskSelect(msg TaskSelectMsg) (tea.Model, error) {
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
			DailyTarget: msg.DailyTarget,
			TimeSpent:   0,
		},
	}

	if err := m.taskService.InsertDailyTask(taskWithDaily.DailyTask); err != nil {
		return m, fmt.Errorf("error in model.go: %v", err)
	}

	todayTaskModel := m.views[viewTodayTask].(TodayTaskModel)
	todayTaskModel = todayTaskModel.AppendTask(taskWithDaily)
	m.views[viewTodayTask] = todayTaskModel

	m.activeView = m.views[viewTodayTask]

	return m, nil
}
