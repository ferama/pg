package sqlview

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/ferama/pg/pkg/db"
	"github.com/ferama/pg/pkg/utils"
)

type MainView struct {
	path        *utils.PathParts
	err         error
	resultsView *ResultsView
	queryView   *QueryView
}

func NewMainView(path *utils.PathParts) *MainView {
	rv := NewResultsView()
	qv := NewQueryView(path)

	return &MainView{
		resultsView: rv,
		queryView:   qv,
		path:        path,
		err:         nil,
	}
}

func (m *MainView) Init() tea.Cmd {
	return nil
}

func (m *MainView) sqlExecute(connString, dbName, schema, query string) (string, error) {
	if query == "" {
		return "", nil
	}
	fields, items, err := db.Query(connString, dbName, schema, query)

	if err != nil {
		return "", err
	}
	return db.RenderQueryResults(items, fields), nil
}

func (m *MainView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyCtrlDown, tea.KeyCtrlD:
			m.resultsView.Viewport.HalfViewDown()
		case tea.KeyCtrlUp, tea.KeyCtrlU:
			m.resultsView.Viewport.HalfViewUp()
		case tea.KeyCtrlX:
			m.resultsView.Viewport.SetContent("running query...")
			go func() {
				query := m.queryView.Textarea.Value()
				res, err := m.sqlExecute(
					m.path.ConfigConnection,
					m.path.DatabaseName,
					m.path.SchemaName,
					query,
				)
				if err != nil {
					m.resultsView.Viewport.SetContent(err.Error())
				} else {
					m.resultsView.Viewport.SetContent(res)
				}
			}()
		default:
			if !m.queryView.Textarea.Focused() {
				cmd := m.queryView.Textarea.Focus()
				cmds = append(cmds, cmd)
			}
		}
	// We handle errors just like any other message
	case error:
		m.err = msg
		return m, nil
	}

	m.resultsView, cmd = m.resultsView.Update(msg)
	cmds = append(cmds, cmd)

	m.queryView, cmd = m.queryView.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m *MainView) View() string {
	return fmt.Sprintf(
		"%s%s",
		m.resultsView.View(),
		m.queryView.View(),
	)
}
