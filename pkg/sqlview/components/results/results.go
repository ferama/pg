package results

import (
	"fmt"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/charmbracelet/bubbles/viewport"
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

	detailsRowStyle = lipgloss.NewStyle().
			PaddingLeft(1)
)

type Model struct {
	table           table.Model
	detailsViewport viewport.Model

	height int
	width  int

	err          error
	focused      bool
	currentState int

	results *db.QueryResults
}

func New() *Model {
	tbl := table.New(nil, 0, 0)
	vp := viewport.New(0, 0)

	return &Model{
		focused:         false,
		results:         nil,
		currentState:    defaultState,
		table:           tbl,
		detailsViewport: vp,
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
			// do not truncate if the result contains just one colum
			// Think about explain queries
			if len(out) > conf.ItemMaxLen && len(results.Columns) > 1 {
				out = out[:conf.ItemMaxLen] + "..."
			}
			row = append(row, out)
		}
		rs = append(rs, row)
	}

	m.table = table.New(upperCols, 0, 0)

	m.table.SetRows(rs)
}

func (m *Model) SetSize(width, height int) {
	m.width = width
	m.height = height
	m.applySize()
}

func (m *Model) applySize() {
	style.Width(m.width - 2)
	style.Height(m.height - 2)

	detailsRowStyle.Width(m.width - 4)

	titleStyle.Width(m.width - 2)
	m.table.SetSize(m.width-2, m.height)

	m.detailsViewport.Width = m.width - 2
	m.detailsViewport.Height = m.height
}

func (m *Model) handleDetailsKeys(msg tea.Msg) {
	if m.currentState != detailsState {
		return
	}
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyDown:
			m.detailsViewport.LineDown(1)
		case tea.KeyUp:
			m.detailsViewport.LineUp(1)
		}
	}
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
		m.results = &db.QueryResults{
			Elapsed: msg.Elapsed,
		}
		m.applySize()

	case editor.QueryResultsMsg:
		m.setResults(msg.Results)
		m.applySize()

	case tea.KeyMsg:
		if !m.focused {
			break
		}
		m.handleDetailsKeys(msg)
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
	if m.currentState == detailsState && m.results != nil {

		m.detailsViewport.SetContent(m.renderDetails())
		return style.Render(
			lipgloss.JoinVertical(lipgloss.Left,
				titleStyle.Render("Item Details"),
				m.detailsViewport.View(),
			),
		)
	}
	title := "Query Results"
	if m.results != nil && m.results.Elapsed != 0 {
		title = fmt.Sprintf("%s [%s]", title, m.results.Elapsed.Round(1*time.Millisecond))
	}
	return style.Render(
		lipgloss.JoinVertical(lipgloss.Left,
			titleStyle.Render(title),
			m.table.View(),
		),
	)
}

func (m *Model) renderDetails() string {
	hStyle := lipgloss.NewStyle().
		Bold(true)

	// https://www.ditig.com/256-colors-cheat-sheet
	evenStyle := lipgloss.NewStyle().
		Padding(0, 1, 0, 1).
		Align(lipgloss.Left).
		Background(lipgloss.Color("235"))
	oddStyle := lipgloss.NewStyle().
		Padding(0, 1, 0, 1).
		Align(lipgloss.Left).
		Background(lipgloss.Color("239"))

	idx := m.table.Cursor()
	row := m.results.Rows[idx]

	var sb strings.Builder
	tw := &tabwriter.Writer{}
	tw.Init(&sb, 0, 4, 2, ' ', 0)

	for i, rawColumn := range m.results.Columns {
		cell := row[i]
		cellWidth := lipgloss.Width(cell)

		cellStyle := evenStyle
		if i%2 == 1 {
			cellStyle = oddStyle
		}
		parts := make([]string, 0)

		if cellWidth > m.width {
			sliceLen := int(m.width / 2)
			i := 0
			for {
				if (i+1)*sliceLen >= cellWidth {
					break
				}
				part := cell[i*sliceLen : (i+1)*sliceLen]
				parts = append(parts, part)
				i++
			}
			parts = append(parts, cell[i*sliceLen:])
		} else {
			parts = append(parts, cell)
		}

		for i, part := range parts {
			col := rawColumn
			if i != 0 {
				col = strings.Repeat(" ", len(rawColumn))
			}
			s := lipgloss.JoinHorizontal(lipgloss.Top,
				hStyle.Render(col),
				cellStyle.Render("\t"),
				cellStyle.Render(part),
				cellStyle.Render("\t"))

			detailsRowStyle.Background(cellStyle.GetBackground())
			s = detailsRowStyle.Render(s)

			fmt.Fprintln(tw, s)
		}
	}
	tw.Flush()

	return sb.String()
}
