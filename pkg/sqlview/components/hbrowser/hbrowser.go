package hbrowser

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ferama/pg/pkg/conf"
	"github.com/ferama/pg/pkg/history"
)

// https://github.com/charmbracelet/bubbletea/blob/master/examples/list-simple/main.go

type HBrowserStatesMsg struct {
}
type listItem struct {
	title, desc string
}

func (i listItem) Title() string       { return i.title }
func (i listItem) Description() string { return i.desc }
func (i listItem) FilterValue() string { return i.title }

var (
	borderStyle = lipgloss.ThickBorder()

	style = lipgloss.NewStyle().
		BorderTop(true).
		BorderRight(true).
		BorderLeft(true).
		BorderBottom(true).
		BorderForeground(lipgloss.Color(conf.ColorFocus)).
		BorderStyle(borderStyle)
)

type Model struct {
	err            error
	focused        bool
	terminalHeight int
	terminalWidth  int

	list list.Model
}

func New() *Model {

	m := &Model{
		err:     nil,
		focused: false,
	}
	m.setState()
	return m
}

func (m *Model) setState() tea.Msg {

	delegate := list.NewDefaultDelegate()

	delegate.Styles.SelectedTitle.
		BorderForeground(lipgloss.AdaptiveColor{Light: conf.ColorFocus, Dark: conf.ColorFocus}).
		Foreground(lipgloss.AdaptiveColor{Light: conf.ColorFocus, Dark: conf.ColorFocus})
	delegate.Styles.SelectedDesc.
		BorderForeground(lipgloss.AdaptiveColor{Light: conf.ColorFocus, Dark: conf.ColorFocus})

	delegate.ShowDescription = true

	listModel := list.New(make([]list.Item, 0), delegate, 0, 0)
	listModel.DisableQuitKeybindings()
	listModel.Styles.Title.
		UnsetBackground().
		Underline(true).
		Foreground(lipgloss.Color(conf.ColorTitle))
	listModel.Styles.FilterPrompt.Foreground(lipgloss.Color(conf.ColorFocus))

	h := history.GetInstance()
	hitems := h.GetList()
	items := make([]list.Item, 0)
	for _, i := range hitems {
		items = append(items, listItem{
			title: i, desc: "",
		})
	}

	delegate.ShowDescription = false
	listModel.SetDelegate(delegate)
	listModel.SetItems(items)
	listModel.Title = "Query History"

	m.list = listModel
	return HBrowserStatesMsg{}
}

func (m *Model) setDimensions() {
	style.Width(m.terminalWidth - 4)
	style.Height(m.terminalHeight - 4)

	m.list.SetSize(m.terminalWidth-4, m.terminalHeight-4)
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Focused() bool {
	return m.focused
}

func (m *Model) Focus() tea.Cmd {
	m.focused = true
	return m.setState
}

func (m *Model) Blur() {
	m.focused = false
}

func (m *Model) Update(msg tea.Msg) (*Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.terminalHeight = msg.Height
		m.terminalWidth = msg.Width
		m.setDimensions()

	case HBrowserStatesMsg:
		m.setDimensions()

	case tea.KeyMsg:
		if !m.focused {
			break
		}
		switch msg.Type {
		case tea.KeyEnter:
			// i := m.list.SelectedItem()
			// if i != nil {
			// 	cmd = m.nextState(i.(listItem))
			// 	cmds = append(cmds, cmd)
			// }
		}

	case error:
		m.err = msg
		return m, nil
	}

	m.list, cmd = m.list.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m *Model) View() string {
	return style.Render(m.list.View())
}
