package results

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ferama/pg/pkg/conf"
	"github.com/ferama/pg/pkg/db"
	"github.com/ferama/pg/pkg/stripansi"
)

var (
	style = lipgloss.NewStyle().
		Border(lipgloss.NormalBorder(), true, true, false, true).
		BorderForeground(lipgloss.Color(conf.ColorBlur))

	focusedStyle = lipgloss.NewStyle().
			Border(lipgloss.ThickBorder(), true, true, false, true).
			BorderForeground(lipgloss.Color(conf.ColorFocus))

	textStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color(conf.ColorBlur))
	textFocusedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(conf.ColorFocus))
)

type Model struct {
	viewport       viewport.Model
	err            error
	terminalHeight int
	terminalWidth  int
	xPosition      int
	focused        bool

	rows   db.ResultsRows
	fields db.ResultsFields
}

func New() *Model {
	vp := viewport.New(5, 5)
	return &Model{
		focused:  false,
		viewport: vp,
	}
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Focus() {
	m.focused = true
}

func (m *Model) Focused() bool {
	return m.focused
}

func (m *Model) Blur() {
	m.focused = false
}

func (m *Model) SetResults(fields db.ResultsFields, rows db.ResultsRows) {
	m.xPosition = 0

	m.rows = rows
	m.fields = fields

	r := db.RenderQueryResults(rows, fields)
	m.viewport.SetContent(r)
}

func (m *Model) SetContent(value string) {
	m.xPosition = 0
	m.viewport.SetContent(value)
}

func (m *Model) scrollHorizontally(amount int) {
	nextPos := m.xPosition + amount
	if nextPos < 0 || nextPos >= len(m.fields) {
		return
	}

	rrows := make(db.ResultsRows, 0)
	for _, r := range m.rows {
		rrows = append(rrows, r[nextPos:])
	}

	rendered := db.RenderQueryResults(rrows, m.fields[nextPos:])
	line := strings.Split(rendered, "\n")[0]
	line = stripansi.Strip(line)
	if len(line) > m.terminalWidth {
		m.viewport.SetContent(rendered)
		m.xPosition = nextPos
	}
}

func (m *Model) setDimensions() {
	style.Width(m.terminalWidth - 2)
	style.Height(m.terminalHeight - (conf.SqlTextareaHeight + 3))

	focusedStyle.Width(m.terminalWidth - 2)
	focusedStyle.Height(m.terminalHeight - (conf.SqlTextareaHeight + 3))

	m.viewport.Width = m.terminalWidth - 2
	m.viewport.Height = m.terminalHeight - (conf.SqlTextareaHeight + 3)
}

func (m *Model) Update(msg tea.Msg) (*Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.terminalHeight = msg.Height
		m.terminalWidth = msg.Width
		m.setDimensions()
	case tea.KeyMsg:
		if !m.focused {
			break
		}
		switch msg.Type {
		case tea.KeyCtrlDown, tea.KeyCtrlD:
			m.viewport.HalfViewDown()
		case tea.KeyCtrlUp, tea.KeyCtrlU:
			m.viewport.HalfViewUp()
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

func (m *Model) View() string {
	var renderedResults string

	percent := fmt.Sprintf("%3.f%%", m.viewport.ScrollPercent()*100)

	borderStyle := lipgloss.NormalBorder()
	if m.focused {
		borderStyle = lipgloss.ThickBorder()
	}
	line := strings.Repeat(borderStyle.Bottom, m.viewport.Width-4)
	line = fmt.Sprintf("%s%s%s%s",
		borderStyle.BottomLeft,
		line,
		percent,
		borderStyle.BottomRight)

	if m.focused {
		renderedResults = focusedStyle.Render(m.viewport.View())
		line = textFocusedStyle.Render(line)
	} else {
		renderedResults = style.Render(m.viewport.View())
		line = textStyle.Render(line)
	}

	out := lipgloss.JoinVertical(lipgloss.Center, renderedResults, line)
	return out
}
