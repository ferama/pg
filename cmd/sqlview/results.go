package sqlview

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ferama/pg/pkg/db"
)

const (
	color = "#66ccff"
)

var (
	resultsModelStyle = lipgloss.NewStyle().
				Border(lipgloss.ThickBorder(), true, true, false, true).
				BorderForeground(lipgloss.Color(color))

	resultsTextStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(color))
)

type resultsView struct {
	viewport       viewport.Model
	err            error
	terminalHeight int
	terminalWidth  int
	xPosition      int

	rows   db.ResultsRows
	fields db.ResultsFields
}

func newResultsView() *resultsView {
	vp := viewport.New(5, 5)
	return &resultsView{
		viewport: vp,
	}
}

func (m *resultsView) Init() tea.Cmd {
	return nil
}

func (m *resultsView) SetResults(fields db.ResultsFields, rows db.ResultsRows) {
	m.xPosition = 0

	m.rows = rows
	m.fields = fields

	r := db.RenderQueryResults(rows, fields)
	m.viewport.SetContent(r)
}

func (m *resultsView) SetContent(value string) {
	m.xPosition = 0
	m.viewport.SetContent(value)
}

func (m *resultsView) scrollHorizontally(amount int) {
	nextPos := m.xPosition + amount
	if nextPos < 0 || nextPos >= len(m.fields) {
		return
	}
	m.xPosition = nextPos

	rrows := make(db.ResultsRows, 0)
	for _, r := range m.rows {
		rrows = append(rrows, r[m.xPosition:])
	}

	rendered := db.RenderQueryResults(rrows, m.fields[m.xPosition:])
	m.viewport.SetContent(rendered)
}

func (m *resultsView) setDimensions() {
	resultsModelStyle.Width(m.terminalWidth - 2)
	resultsModelStyle.Height(m.terminalHeight - (sqlTextareaHeight + 3))

	m.viewport.Width = m.terminalWidth - 2
	m.viewport.Height = m.terminalHeight - (sqlTextareaHeight + 3)
}

func (m *resultsView) Update(msg tea.Msg) (*resultsView, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.terminalHeight = msg.Height
		m.terminalWidth = msg.Width
		m.setDimensions()
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyLeft:
			m.scrollHorizontally(-1)
		case tea.KeyRight:
			m.scrollHorizontally(1)
		}
	// We handle errors just like any other message
	case error:
		m.err = msg
		return m, nil
	}

	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m *resultsView) View() string {
	var renderedResults string

	percent := fmt.Sprintf("%3.f%%", m.viewport.ScrollPercent()*100)
	line := strings.Repeat("━", m.viewport.Width-4)
	line = fmt.Sprintf("┗%s%s┛", line, percent)

	renderedResults = resultsModelStyle.Render(m.viewport.View())
	line = resultsTextStyle.Render(line)

	out := lipgloss.JoinVertical(lipgloss.Center, renderedResults, line)
	return fmt.Sprintf(
		"%s\n",
		out,
	)
}
