package ui

import (
	"strconv"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/jeisaRaja/tasktimer/internal/models"
)

type DailyTargetInput struct {
	input textinput.Model
	task  models.Task
}

func initialDailyTargetInput() DailyTargetInput {
	ti := textinput.New()
	ti.Placeholder = "Input daily target in minutes (optional)"
	ti.Focus()
	return DailyTargetInput{
		input: ti,
	}
}

func (m DailyTargetInput) Init() tea.Cmd {
	return textinput.Blink
}

func (m DailyTargetInput) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case DailyTargetInputMsg:
		m.task = models.Task(msg)
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			task := models.TaskWithDaily{
				Task: m.task,
			}
			dailyTarget, err := strconv.Atoi(m.input.Value())
			if err != nil {
				panic(err)
			}
			task.DailyTarget = time.Duration(dailyTarget * int(time.Minute))
			return m, func() tea.Msg {
				return TaskSelectMsg(task)
			}
		default:
			if isNumeric(msg.String()) || msg.String() == "backspace" {
				m.input, cmd = m.input.Update(msg)
				return m, cmd
			}
		}
	}
	return m, nil
}

func (m DailyTargetInput) View() string {
	return m.input.View()
}
