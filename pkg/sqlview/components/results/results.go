package results

import (
	"fmt"
	"strings"
	"text/tabwriter"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ferama/pg/pkg/components/table"
	"github.com/ferama/pg/pkg/conf"
	"github.com/ferama/pg/pkg/db"
	"github.com/ferama/pg/pkg/sqlview/components/editor"
)

const (
	defaultState int = iota
	detailsState
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

	titleStyle = lipgloss.NewStyle().
			Align(lipgloss.Center).
			Bold(true).
			Foreground(lipgloss.Color(conf.ColorTitle)).
			PaddingBottom(1).
			Underline(true)
)

type Model struct {
	table table.Model

	err            error
	terminalHeight int
	terminalWidth  int
	focused        bool
	currentState   int

	results *db.QueryResults
}

func New() *Model {
	tbl := table.New(nil, 0, 0)

	return &Model{
		focused:      false,
		results:      nil,
		currentState: defaultState,
		table:        tbl,
	}
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) HandleEsc() bool {
	if !m.focused {
		return false
	}
	if m.currentState == detailsState {
		m.currentState = defaultState
		return true
	}
	return false
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
	m.results = results

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

	titleStyle.Width(m.terminalWidth - 2)
	m.table.SetSize(m.terminalWidth-2, m.terminalHeight-(conf.SqlTextareaHeight+5))
}

func (m *Model) Update(msg tea.Msg) (*Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case editor.QueryStatusMsg:
		m.currentState = defaultState
		m.table = table.New([]string{"STATUS"}, 0, 0)
		m.table.SetRows([]table.Row{
			table.SimpleRow{msg.Content},
		})
		m.setDimensions()

	case editor.QueryResultsMsg:
		m.setResults(msg.Results)

	case tea.WindowSizeMsg:
		m.terminalHeight = msg.Height
		m.terminalWidth = msg.Width
		m.setDimensions()
	case tea.KeyMsg:
		if !m.focused {
			break
		}
		switch msg.Type {
		case tea.KeyEnter:
			m.currentState = detailsState
		}
	case error:
		m.err = msg
		return m, nil
	}

	if m.currentState == defaultState {
		m.table, cmd = m.table.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m *Model) View() string {
	hStyle := lipgloss.NewStyle().
		PaddingLeft(2).
		Bold(true)

	if m.currentState == detailsState && m.results != nil {
		idx := m.table.Cursor()
		row := m.results.Rows[idx]

		var sb strings.Builder
		tw := &tabwriter.Writer{}
		tw.Init(&sb, 0, 4, 2, ' ', 0)
		for i, col := range m.results.Columns {
			c := hStyle.Render(col)
			r := row[i]
			s := fmt.Sprintf("%s\t%s", c, r)
			fmt.Fprintln(tw, s)
		}
		tw.Flush()
		return style.Render(
			lipgloss.JoinVertical(lipgloss.Left,
				titleStyle.Render("Item Details"),
				sb.String(),
			),
		)
	}
	return style.Render(
		lipgloss.JoinVertical(lipgloss.Left,
			titleStyle.Render("Query Results"),
			m.table.View(),
		),
	)
}
