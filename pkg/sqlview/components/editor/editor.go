package editor

import (
	"time"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ferama/bubble-texteditor/texteditor"
	"github.com/ferama/pg/pkg/conf"
	"github.com/ferama/pg/pkg/db"
	"github.com/ferama/pg/pkg/history"
	"github.com/ferama/pg/pkg/sqlview/components/hbrowser"
	"github.com/ferama/pg/pkg/utils"
)

var (
	style = lipgloss.NewStyle().
		Border(lipgloss.ThickBorder(), false, true, false, true).
		BorderForeground(lipgloss.Color(conf.ColorFocus))
)

type QueryStatusMsg struct {
	Content string
	Elapsed time.Duration
}

type QueryResultsMsg struct {
	Results *db.QueryResults
}

type Model struct {
	path       *utils.PathParts
	texteditor texteditor.Model

	history *history.History
	err     error
}

func New(path *utils.PathParts) *Model {
	te := texteditor.New()

	te.SetSyntax("sql")
	te.SetWidth(10)
	te.SetHeight(conf.SqlTextareaHeight)
	te.Focus()

	return &Model{
		path:       path,
		texteditor: te,
		history:    history.GetInstance(),
		err:        nil,
	}
}
func (m *Model) Focus() tea.Cmd {
	style.
		BorderStyle(lipgloss.ThickBorder()).
		BorderForeground(lipgloss.Color(conf.ColorFocus))
	return m.texteditor.Focus()
}

func (m *Model) Blur() {
	style.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color(conf.ColorBlur))
	m.texteditor.Blur()
}

func (m *Model) SetValue(value string) {
	m.texteditor.SetValue(value)
}

func (m *Model) value() string {
	return m.texteditor.Value()
}

func (m *Model) sqlExecute(connString, dbName, schema, query string) (*db.QueryResults, error) {
	results, err := db.Query(connString, dbName, schema, query)
	if err != nil {
		return nil, err
	}
	return results, nil
}

func (m *Model) doQuery() tea.Cmd {
	return func() tea.Msg {
		query := m.value()

		m.history.Append(query)

		results, err := m.sqlExecute(
			m.path.ConfigConnection,
			m.path.DatabaseName,
			m.path.SchemaName,
			query,
		)
		if err != nil {
			return QueryStatusMsg{
				Content: err.Error(),
			}
		} else {
			if len(results.Columns) == 0 {
				return QueryStatusMsg{
					Content: "done",
					Elapsed: results.Elapsed,
				}
			} else {
				return QueryResultsMsg{
					Results: results,
				}
			}
		}
	}
}

func (m *Model) Init() tea.Cmd {
	return tea.Batch(textarea.Blink)
}
func (m *Model) Update(msg tea.Msg) (*Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case hbrowser.HBrowserSelectedMsg:
		q, err := m.history.GetAtIdx(msg.Idx)
		if err == nil {
			m.texteditor.SetValue(q)
		}
	case tea.WindowSizeMsg:
		m.texteditor.SetWidth(msg.Width - 2)

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyShiftDown:
			q, err := m.history.GoNext()
			if err == nil {
				m.texteditor.SetValue(q)
			}

		case tea.KeyShiftUp:
			q, err := m.history.GoPrev()
			if err == nil {
				m.texteditor.SetValue(q)
			}

		case tea.KeyCtrlD:
			m.SetValue("")

		case tea.KeyCtrlX:
			cmd = func() tea.Msg {
				return QueryStatusMsg{
					Content: "running query...",
				}
			}
			cmds = append(cmds, cmd)

			cmd = m.doQuery()
			cmds = append(cmds, cmd)
		}
	case error:
		m.err = msg
		return m, nil
	}

	m.texteditor, cmd = m.texteditor.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m *Model) View() string {
	return style.Render(m.texteditor.View())
}
