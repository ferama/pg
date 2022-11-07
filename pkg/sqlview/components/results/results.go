package results

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ferama/pg/pkg/components/table"
	"github.com/ferama/pg/pkg/conf"
	"github.com/ferama/pg/pkg/db"
	"github.com/ferama/pg/pkg/sqlview/components/query"
)

var (
	borderStyle = lipgloss.NormalBorder()

	style = lipgloss.NewStyle().
		BorderTop(true).
		BorderRight(true).
		BorderLeft(true).
		BorderBottom(true).
		BorderForeground(lipgloss.Color(conf.ColorBlur)).
		BorderStyle(borderStyle)
)

type Model struct {
	table table.Model

	err            error
	terminalHeight int
	terminalWidth  int
	focused        bool
}

func New() *Model {
	tbl := table.New(nil, 0, 0)

	return &Model{
		focused: false,
		table:   tbl,
	}
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Focus() {
	borderStyle = lipgloss.ThickBorder()
	style = style.BorderStyle(borderStyle).
		BorderForeground(lipgloss.Color(conf.ColorFocus))

	m.focused = true
	m.table.Focus()
}

func (m *Model) Focused() bool {
	return m.focused
}

func (m *Model) Blur() {
	borderStyle = lipgloss.NormalBorder()
	style = style.BorderStyle(borderStyle).
		BorderForeground(lipgloss.Color(conf.ColorBlur))

	m.focused = false
	m.table.Blur()
}

func (m *Model) setResults(results *db.QueryResults) {

	upperCols := make(db.Columns, 0)
	for _, c := range results.Columns {
		upperCols = append(upperCols, strings.ToUpper(c))
	}

	var rs []table.Row
	for _, r := range results.Rows {
		row := table.SimpleRow{}
		for _, v := range r {
			out := v
			if len(out) > conf.ItemMaxLen {
				out = out[:conf.ItemMaxLen] + "..."
			}
			row = append(row, out)
		}
		rs = append(rs, row)
	}

	m.table = table.New(upperCols, 0, 0)

	m.table.SetRows(rs)
	m.setDimensions()
}

func (m *Model) setDimensions() {
	style.Width(m.terminalWidth - 2)
	style.Height(m.terminalHeight - (conf.SqlTextareaHeight + 3))

	m.table.SetSize(m.terminalWidth-2, m.terminalHeight-(conf.SqlTextareaHeight+3))
}

func (m *Model) Update(msg tea.Msg) (*Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case query.QueryStatusMsg:
		m.table = table.New([]string{"status"}, 0, 0)
		m.table.SetRows([]table.Row{
			table.SimpleRow{msg.Content},
		})
		m.setDimensions()

	case query.QueryResultsMsg:
		m.setResults(msg.Results)

	case tea.WindowSizeMsg:
		m.terminalHeight = msg.Height
		m.terminalWidth = msg.Width
		m.setDimensions()
	// We handle errors just like any other message
	case error:
		m.err = msg
		return m, nil
	}

	m.table, cmd = m.table.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m *Model) View() string {
	return style.Render(m.table.View())
}
