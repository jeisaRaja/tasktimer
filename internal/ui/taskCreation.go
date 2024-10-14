package ui

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jeisaRaja/tasktimer/internal/models"
)

var (
	focusedStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	blurredStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	cursorStyle         = focusedStyle
	noStyle             = lipgloss.NewStyle()
	helpStyle           = blurredStyle
	cursorModeHelpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("244"))

	focusedButton = focusedStyle.Render("[ Submit ]")
	blurredButton = fmt.Sprintf("[ %s ]", blurredStyle.Render("Submit"))
)

type TaskCreationModel struct {
	focusIndex int
	inputs     []textinput.Model
	fieldLen   int
	cursorMode cursor.Mode
}

func initialTaskCreation() TaskCreationModel {
	m := TaskCreationModel{
		inputs: make([]textinput.Model, 2),
	}

	nameInput := textinput.New()
	nameInput.Placeholder = "New Task Name"
	nameInput.PromptStyle = focusedStyle
	nameInput.TextStyle = focusedStyle
	nameInput.Focus()

	descriptionInput := textinput.New()
	descriptionInput.Placeholder = "Description (optional)"

	durationInput := textinput.New()
	durationInput.Placeholder = "Weekly Target Duration (in hours)"

	tagsInput := textinput.New()
	tagsInput.Placeholder = "Tags (optional, comma-separated)"

	inputs := []textinput.Model{nameInput, descriptionInput, durationInput, tagsInput}

	m.inputs = inputs
	m.fieldLen = len(inputs)

	return m
}

func (m TaskCreationModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m TaskCreationModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+s":
			durationString := m.inputs[2].Value()
			durationInt, err := strconv.Atoi(durationString)
			if err != nil {
				fmt.Println("Error converting string to int:", err)
			}

			duration := time.Duration(durationInt) * time.Second

			fmt.Println("Duration:", duration)

			tags := strings.Split(m.inputs[3].Value(), ",")
			task := models.Task{
				Name:         m.inputs[0].Value(),
				Description:  m.inputs[1].Value(),
				WeeklyTarget: duration,
				Tags:         tags,
			}
			return m, createInsertTaskMsg(task)
		case "ctrl+c", "esc":
			return m, tea.Quit

		case "ctrl+r":
			m.cursorMode++
			if m.cursorMode > cursor.CursorHide {
				m.cursorMode = cursor.CursorBlink
			}
			cmds := make([]tea.Cmd, len(m.inputs))
			for i := range m.inputs {
				cmds[i] = m.inputs[i].Cursor.SetMode(m.cursorMode)
			}
			return m, tea.Batch(cmds...)

		case "tab", "shift+tab", "enter", "up", "down":
			s := msg.String()

			if s == "enter" && m.focusIndex == len(m.inputs) {
				return m, tea.Quit
			}

			if s == "up" || s == "shift+tab" {
				m.focusIndex--
			} else {
				m.focusIndex++
			}

			if m.focusIndex >= len(m.inputs) {
				m.focusIndex = 0
			} else if m.focusIndex < 0 {
				m.focusIndex = len(m.inputs) - 1
			}

			cmds := make([]tea.Cmd, len(m.inputs))
			for i := 0; i <= len(m.inputs)-1; i++ {
				if i == m.focusIndex {
					cmds[i] = m.inputs[i].Focus()
					m.inputs[i].PromptStyle = focusedStyle
					m.inputs[i].TextStyle = focusedStyle
					continue
				}
				m.inputs[i].Blur()
				m.inputs[i].PromptStyle = noStyle
				m.inputs[i].TextStyle = noStyle
			}

			return m, tea.Batch(cmds...)
		}
	}

	// Handle character input and blinking
	cmd := m.updateInputs(msg)

	return m, cmd
}

func (m *TaskCreationModel) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))

	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

func (m TaskCreationModel) View() string {
	var b strings.Builder

	for i := range m.inputs {
		b.WriteString(m.inputs[i].View())
		if i < len(m.inputs)-1 {
			b.WriteRune('\n')
		}
	}

	return b.String()
}
