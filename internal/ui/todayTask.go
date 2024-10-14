package ui

import (
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/jeisaRaja/tasktimer/internal/models"
)

type TodayTaskModel struct {
	day       time.Weekday
	tasks     []models.Task
	currIndex int
}

func initialTodayTaskModel() TodayTaskModel {
	currDay := time.Now().Weekday()
	var tasks []models.Task

	tasks = append(tasks, models.Task{Name: "Code"})
	tasks = append(tasks, models.Task{Name: "Read"})
	tasks = append(tasks, models.Task{Name: "Run"})
	tasks = append(tasks, models.Task{Name: "Walk"})
	tasks = append(tasks, models.Task{Name: "Watch"})

	m := TodayTaskModel{
		day:   currDay,
		tasks: tasks,
	}

	return m
}

func (m TodayTaskModel) Init() tea.Cmd {
	return nil
}

func (m TodayTaskModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "j", "down", "tab":
			m.currIndex++
			if m.currIndex > len(m.tasks)-1 {
				m.currIndex = len(m.tasks) - 1
			}
		case "k", "up", "shift+tab":
			m.currIndex--
			if m.currIndex < 0 {
				m.currIndex = 0
			}
		}
	}

	return m, nil
}

func (m TodayTaskModel) View() string {
	var b strings.Builder
	for i := range m.tasks {
		if i == m.currIndex {
			b.WriteString(">")
			b.WriteString(m.tasks[i].Name)
		} else {
			b.WriteString(m.tasks[i].Name)
		}
		b.WriteRune('\n')
	}

	return b.String()
}
