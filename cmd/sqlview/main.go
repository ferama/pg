package sqlview

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ferama/pg/cmd/sqlview/components/query"
	"github.com/ferama/pg/cmd/sqlview/components/results"
	"github.com/ferama/pg/cmd/sqlview/components/statusbar"
	"github.com/ferama/pg/pkg/db"
	"github.com/ferama/pg/pkg/utils"
)

type MainView struct {
	path        *utils.PathParts
	err         error
	resultsView *results.Model
	queryView   *query.Model
	statsuBar   *statusbar.Model
}

func NewMainView(path *utils.PathParts) *MainView {
	rv := results.New()
	qv := query.New(path)
	sb := statusbar.New()

	return &MainView{
		resultsView: rv,
		queryView:   qv,
		statsuBar:   sb,
		path:        path,
		err:         nil,
	}
}

func (m *MainView) Init() tea.Cmd {
	// m.queryView.SetValue("select * from pg_replication_slots")
	// m.queryView.SetValue("select * from sales limit 100")
	return tea.Batch(m.queryView.Init(), m.queryView.Focus(), tea.EnterAltScreen)
}

func (m *MainView) sqlExecute(connString, dbName, schema, query string) (db.ResultsFields, db.ResultsRows, error) {
	if query == "" {
		return nil, nil, nil
	}
	fields, items, err := db.Query(connString, dbName, schema, query)

	if err != nil {
		return nil, nil, err
	}
	// return db.RenderQueryResults(items, fields), nil
	return fields, items, nil
}

func (m *MainView) doQuery(query string) (db.ResultsFields, db.ResultsRows, error) {
	fields, items, err := m.sqlExecute(
		m.path.ConfigConnection,
		m.path.DatabaseName,
		m.path.SchemaName,
		query,
	)

	return fields, items, err
}

func (m *MainView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEsc, tea.KeyCtrlC:
			return m, tea.Quit
		case tea.KeyTab:
			if m.resultsView.Focused() {
				m.resultsView.Blur()
				cmd = m.queryView.Focus()
				cmds = append(cmds, cmd)
			} else {
				m.resultsView.Focus()
				m.queryView.Blur()
			}
		case tea.KeyCtrlX:
			m.resultsView.SetContent("running query...")

			query := m.queryView.Value()
			go func() {
				fields, items, err := m.doQuery(query)
				if err != nil {
					m.resultsView.SetContent(err.Error())
				} else {
					m.resultsView.SetResults(fields, items)
				}
			}()
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

	m.statsuBar, cmd = m.statsuBar.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m *MainView) View() string {
	return lipgloss.JoinVertical(lipgloss.Left,
		m.resultsView.View(),
		m.queryView.View(),
		m.statsuBar.View(),
	)
}
