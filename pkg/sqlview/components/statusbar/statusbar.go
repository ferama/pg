package statusbar

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ferama/pg/pkg/conf"
	"github.com/ferama/pg/pkg/utils"
)

var (
	shortcutStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("#cccccc")).
			Foreground(lipgloss.Color("#000000")).
			Align(lipgloss.Right).
			Border(lipgloss.ThickBorder(), false, true, false, false).
			BorderForeground(lipgloss.Color(conf.ColorBlur))

	infoStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("#bbbbbb")).
			Foreground(lipgloss.Color("#000000")).
			Align(lipgloss.Left).
			Border(lipgloss.ThickBorder(), false, false, false, true).
			BorderForeground(lipgloss.Color(conf.ColorBlur))
)

type Model struct {
	path *utils.PathParts
}

func New(path *utils.PathParts) *Model {
	m := &Model{
		path: path,
	}
	return m
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Update(msg tea.Msg) (*Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		shortcutStyle = shortcutStyle.Width(msg.Width/2 - 2)
		infoStyle = infoStyle.Width(msg.Width - shortcutStyle.GetWidth() - 2)
	}
	return m, tea.Batch(cmds...)
}

func (m *Model) View() string {
	out := "<Ctrl+x> → execute・<ESC> → exit・"
	path := fmt.Sprintf("・%s/%s/%s",
		m.path.ConfigConnection,
		m.path.DatabaseName,
		m.path.SchemaName,
	)
	r := lipgloss.JoinHorizontal(lipgloss.Left,
		infoStyle.Render(path),
		shortcutStyle.Render(out))
	return r
}
