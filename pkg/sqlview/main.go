package sqlview

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ferama/pg/pkg/sqlview/components/browser"
	"github.com/ferama/pg/pkg/sqlview/components/query"
	"github.com/ferama/pg/pkg/sqlview/components/results"
	"github.com/ferama/pg/pkg/sqlview/components/statusbar"
	"github.com/ferama/pg/pkg/utils"
)

const (
	queryLayoutStatus int = iota
	browserLayoutStatus
)

type MainView struct {
	path *utils.PathParts
	err  error

	resultsView *results.Model
	queryView   *query.Model
	statsuBar   *statusbar.Model
	browserView *browser.Model

	activeLayout int
}

func NewMainView(path *utils.PathParts) *MainView {
	resultsView := results.New()
	queryView := query.New(path)
	statusBar := statusbar.New(path)
	browserView := browser.New()

	return &MainView{
		resultsView: resultsView,
		queryView:   queryView,
		statsuBar:   statusBar,
		browserView: browserView,

		path: path,
		err:  nil,

		activeLayout: queryLayoutStatus,
	}
}

func (m *MainView) Init() tea.Cmd {
	// m.queryView.SetValue("select * from pg_replication_slots")
	// m.queryView.SetValue("select * from sales limit 100")
	return tea.Batch(m.queryView.Init(), m.queryView.Focus(), tea.EnterAltScreen)
}

func (m *MainView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEsc:
			if m.activeLayout == browserLayoutStatus {
				m.activeLayout = queryLayoutStatus
			} else {
				return m, tea.Quit
			}
		case tea.KeyCtrlO:
			m.activeLayout = browserLayoutStatus
		case tea.KeyTab:
			if m.resultsView.Focused() {
				m.resultsView.Blur()
				cmd = m.queryView.Focus()
				cmds = append(cmds, cmd)
			} else {
				m.resultsView.Focus()
				m.queryView.Blur()
			}
		}
	// We handle errors just like any other message
	case error:
		m.err = msg
		return m, nil
	}

	m.resultsView, cmd = m.resultsView.Update(msg)
	cmds = append(cmds, cmd)

	m.browserView, cmd = m.browserView.Update(msg)
	cmds = append(cmds, cmd)

	m.queryView, cmd = m.queryView.Update(msg)
	cmds = append(cmds, cmd)

	m.statsuBar, cmd = m.statsuBar.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m *MainView) View() string {
	switch m.activeLayout {
	case queryLayoutStatus:
		return lipgloss.JoinVertical(lipgloss.Left,
			m.resultsView.View(),
			m.queryView.View(),
			m.statsuBar.View(),
		)
	case browserLayoutStatus:
		return lipgloss.JoinVertical(lipgloss.Left,
			m.browserView.View(),
			m.statsuBar.View(),
		)
	}
	return ""
}
