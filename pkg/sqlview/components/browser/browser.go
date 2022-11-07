package browser

import (
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	err            error
	terminalHeight int
	terminalWidth  int
	focused        bool
}

func New() *Model {

	return &Model{
		err:     nil,
		focused: false,
	}
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Focused() bool {
	return m.focused
}

func (m *Model) Blur() {
	m.focused = false
}

func (m *Model) Update(msg tea.Msg) (*Model, tea.Cmd) {
	var cmds []tea.Cmd
	// var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.terminalHeight = msg.Height
		m.terminalWidth = msg.Width
	// We handle errors just like any other message
	case error:
		m.err = msg
		return m, nil
	}

	return m, tea.Batch(cmds...)
}

func (m *Model) View() string {
	return "Browser"
}
