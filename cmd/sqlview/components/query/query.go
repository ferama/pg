package query

import (
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ferama/pg/pkg/conf"
	"github.com/ferama/pg/pkg/utils"
)

var (
	borderStyle = lipgloss.NewStyle().
		Border(lipgloss.ThickBorder(), false, true, false, true).
		BorderForeground(lipgloss.Color(conf.ColorFocus))
)

type Model struct {
	path     *utils.PathParts
	textarea textarea.Model
	err      error
}

func New(path *utils.PathParts) *Model {
	ta := textarea.New()
	ta.Placeholder = "select ..."

	// ta.Prompt = lipgloss.ThickBorder().Left + " "
	ta.Prompt = ""

	ta.SetWidth(10)
	ta.SetHeight(conf.SqlTextareaHeight)
	ta.Focus()

	return &Model{
		path:     path,
		textarea: ta,
		err:      nil,
	}
}
func (m *Model) Focus() tea.Cmd {
	borderStyle.
		BorderStyle(lipgloss.ThickBorder()).
		BorderForeground(lipgloss.Color(conf.ColorFocus))
	return m.textarea.Focus()
}

func (m *Model) Blur() {
	borderStyle.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color(conf.ColorBlur))
	m.textarea.Blur()
}

func (m *Model) SetValue(value string) {
	m.textarea.SetValue(value)
}

func (m *Model) Value() string {
	return m.textarea.Value()
}

func (m *Model) Init() tea.Cmd {
	return tea.Batch(textarea.Blink)
}
func (m *Model) Update(msg tea.Msg) (*Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.textarea.SetWidth(msg.Width - 2)
	// We handle errors just like any other message
	case error:
		m.err = msg
		return m, nil
	}

	m.textarea, cmd = m.textarea.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m *Model) View() string {
	return borderStyle.Render(m.textarea.View())
}
