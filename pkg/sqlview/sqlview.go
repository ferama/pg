package sqlview

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ferama/pg/pkg/conf"
	"github.com/ferama/pg/pkg/sqlview/components/editor"
	"github.com/ferama/pg/pkg/sqlview/components/hbrowser"
	"github.com/ferama/pg/pkg/sqlview/components/results"
	"github.com/ferama/pg/pkg/sqlview/components/statusbar"
	"github.com/ferama/pg/pkg/utils"
)

const (
	defaultState int = iota
	historyState
)

type MainView struct {
	path *utils.PathParts
	err  error

	terminalHeight int
	terminalWidth  int

	resultsView *results.Model
	queryView   *editor.Model
	statsuBar   *statusbar.Model

	currentState int

	historyBrowser *hbrowser.Model
}

func NewMainView(path *utils.PathParts) *MainView {
	queryView := editor.New(path)

	if path.TableName != "" {
		query := fmt.Sprintf("select *\nfrom %s\nlimit 10", path.TableName)

		queryView.SetValue(query)
	}

	return &MainView{
		resultsView: results.New(),
		queryView:   queryView,
		statsuBar:   statusbar.New(path),

		historyBrowser: hbrowser.New(),

		currentState: defaultState,

		path: path,
		err:  nil,
	}
}

func (m *MainView) setState() tea.Cmd {
	var cmd tea.Cmd

	switch m.currentState {
	case defaultState:
		m.queryView.Focus()
		m.resultsView.Blur()
		m.historyBrowser.Blur()
	case historyState:
		cmd = m.historyBrowser.Focus()
		m.queryView.Blur()
		m.resultsView.Blur()
	}

	return cmd
}

func (m *MainView) Init() tea.Cmd {
	return tea.Batch(
		m.queryView.Init(),
		m.queryView.Focus(),
		m.statsuBar.Init(),
		tea.EnterAltScreen,
	)
}

func (m *MainView) setDimensions() {
	m.resultsView.SetSize(m.terminalWidth, m.terminalHeight-(conf.SqlTextareaHeight+5))
	m.historyBrowser.SetSize(m.terminalWidth, m.terminalHeight)
}

func (m *MainView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case hbrowser.HBrowserSelectedMsg:
		m.currentState = defaultState
		cmd = m.setState()
		cmds = append(cmds, cmd)

	case tea.WindowSizeMsg:
		m.terminalHeight = msg.Height
		m.terminalWidth = msg.Width
		m.setDimensions()

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEsc, tea.KeyCtrlC:
			if m.currentState != defaultState {
				m.currentState = defaultState
			} else {
				if m.resultsView.HandleEsc() {
					break
				}
				return m, tea.Quit
			}
			cmd = m.setState()
			cmds = append(cmds, cmd)

		case tea.KeyCtrlO:
			m.currentState = historyState
			cmd = m.setState()
			cmds = append(cmds, cmd)

		case tea.KeyTab:
			if m.resultsView.Focused() {
				m.resultsView.Blur()
				m.statsuBar.Focus()
				cmd = m.queryView.Focus()
				cmds = append(cmds, cmd)
			} else {
				m.resultsView.Focus()
				m.statsuBar.Blur()
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

	m.queryView, cmd = m.queryView.Update(msg)
	cmds = append(cmds, cmd)

	m.statsuBar, cmd = m.statsuBar.Update(msg)
	cmds = append(cmds, cmd)

	m.historyBrowser, cmd = m.historyBrowser.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m *MainView) View() string {
	switch m.currentState {
	case historyState:
		return m.historyBrowser.View()
	default:
		return lipgloss.JoinVertical(lipgloss.Left,
			m.resultsView.View(),
			m.queryView.View(),
			m.statsuBar.View(),
		)
	}
}
