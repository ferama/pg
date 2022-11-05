package results

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ferama/pg/pkg/conf"
	"github.com/ferama/pg/pkg/db"
	"github.com/ferama/pg/pkg/stripansi"
)

var (
	borderStyle = lipgloss.NormalBorder()

	headerStyle = lipgloss.NewStyle().
			Background(lipgloss.Color(lipgloss.Color(conf.ColorFocus))).
			BorderTop(true).
			BorderRight(true).
			BorderLeft(true).
			BorderForeground(lipgloss.Color(conf.ColorBlur)).
			BorderStyle(borderStyle)

	style = lipgloss.NewStyle().
		BorderRight(true).
		BorderLeft(true).
		BorderForeground(lipgloss.Color(conf.ColorBlur)).
		BorderStyle(borderStyle)

	textStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(conf.ColorBlur))
)

type Model struct {
	viewport viewport.Model
	header   viewport.Model

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
	ha := viewport.New(5, 1)

	vp.KeyMap.Down = key.Binding{}
	vp.KeyMap.Up = key.Binding{}

	return &Model{
		focused:  false,
		viewport: vp,
		header:   ha,
	}
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Focus() {
	borderStyle = lipgloss.ThickBorder()
	style = style.BorderStyle(borderStyle).
		BorderForeground(lipgloss.Color(conf.ColorFocus))

	headerStyle = headerStyle.BorderStyle(borderStyle).
		BorderForeground(lipgloss.Color(conf.ColorFocus))

	textStyle = textStyle.Foreground(lipgloss.Color(conf.ColorFocus))

	m.focused = true
}

func (m *Model) Focused() bool {
	return m.focused
}

func (m *Model) Blur() {
	borderStyle = lipgloss.NormalBorder()
	style = style.BorderStyle(borderStyle).
		BorderForeground(lipgloss.Color(conf.ColorBlur))

	headerStyle = headerStyle.BorderStyle(borderStyle).
		BorderForeground(lipgloss.Color(conf.ColorBlur))

	textStyle = textStyle.Foreground(lipgloss.Color(conf.ColorBlur))

	m.focused = false
}

func (m *Model) SetResults(fields db.ResultsFields, rows db.ResultsRows) {
	m.xPosition = 0

	m.rows = rows
	m.fields = fields

	r := db.RenderQueryResults(rows, fields)
	parts := strings.Split(r, "\n")
	m.header.SetContent(parts[0])
	renderedContent := strings.Join(parts[1:], "\n")
	m.viewport.SetContent(renderedContent)
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
	parts := strings.Split(rendered, "\n")
	header := parts[0]
	if len(stripansi.Strip(header)) > m.terminalWidth {
		m.header.SetContent(header)

		renderedContent := strings.Join(parts[1:], "\n")
		m.viewport.SetContent(renderedContent)
		m.xPosition = nextPos
	}
}

func (m *Model) setDimensions() {
	style.Width(m.terminalWidth - 2)
	style.Height(m.terminalHeight - (conf.SqlTextareaHeight + 4))

	m.viewport.Width = m.terminalWidth - 2
	m.viewport.Height = m.terminalHeight - (conf.SqlTextareaHeight + 4)

	headerStyle.Width(m.terminalWidth - 2)
	m.header.Width = m.terminalWidth - 2
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
		case tea.KeyDown:
			m.viewport.LineDown(1)
		case tea.KeyUp:
			m.viewport.LineUp(1)
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

	m.header, cmd = m.header.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m *Model) View() string {
	var renderedResults string

	percent := fmt.Sprintf("%3.f%%", m.viewport.ScrollPercent()*100)

	bottomLine := strings.Repeat(borderStyle.Bottom, m.viewport.Width-4)
	bottomLine = fmt.Sprintf("%s%s%s%s",
		borderStyle.BottomLeft,
		bottomLine,
		percent,
		borderStyle.BottomRight)

	renderedResults = style.Render(m.viewport.View())
	bottomLine = textStyle.Render(bottomLine)

	out := lipgloss.JoinVertical(lipgloss.Center,
		headerStyle.Render(m.header.View()),
		renderedResults,
		bottomLine)
	return out
}
