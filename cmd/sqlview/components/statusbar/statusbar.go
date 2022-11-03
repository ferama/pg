package statusbar

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ferama/pg/pkg/conf"
)

var (
	barStyle = lipgloss.NewStyle().
		Background(lipgloss.Color("#cccccc")).
		Foreground(lipgloss.Color("#000000")).
		Align(lipgloss.Right).
		Border(lipgloss.ThickBorder(), false, true).
		BorderForeground(lipgloss.Color(conf.ColorBlur))
)

type Model struct {
}

func New() *Model {
	m := &Model{}
	return m
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Update(msg tea.Msg) (*Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		barStyle = barStyle.Width(msg.Width - 2)
	}
	return m, tea.Batch(cmds...)
}

func (m *Model) View() string {
	out := "|ctrl+x| execute |ESC| exit"
	return barStyle.Render(out)
}
