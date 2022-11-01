package sqlview

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	resultsModelStyle = lipgloss.NewStyle().
				Border(lipgloss.ThickBorder(), true, true, false, true).
				BorderForeground(lipgloss.Color("#ffffff"))

	infoStyle = lipgloss.NewStyle()
)

type ResultsView struct {
	Viewport viewport.Model
	err      error
}

func NewResultsView() *ResultsView {
	vp := viewport.New(5, 5)
	return &ResultsView{
		Viewport: vp,
	}
}

func (m *ResultsView) Init() tea.Cmd {
	return nil
}

func (m *ResultsView) Update(msg tea.Msg) (*ResultsView, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		resultsModelStyle.Width(msg.Width - 2)
		resultsModelStyle.Height(msg.Height - (SqlTextareaHeight + 3))
		m.Viewport.Width = msg.Width - 2
		m.Viewport.Height = msg.Height - (SqlTextareaHeight + 3)
	// We handle errors just like any other message
	case error:
		m.err = msg
		return m, nil
	}

	m.Viewport, cmd = m.Viewport.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m *ResultsView) View() string {
	renderedResults := resultsModelStyle.Render(m.Viewport.View())
	percent := fmt.Sprintf("%3.f%%", m.Viewport.ScrollPercent()*100)
	line := strings.Repeat("━", m.Viewport.Width-4)
	line = fmt.Sprintf("┗%s%s┛", line, infoStyle.Render(percent))
	out := lipgloss.JoinVertical(lipgloss.Center, renderedResults, line)
	return fmt.Sprintf(
		"%s\n",
		out,
	)
}
