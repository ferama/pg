package query

import (
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ferama/pg/pkg/conf"
	"github.com/ferama/pg/pkg/utils"
)

var (
	style = lipgloss.NewStyle().
		Foreground(lipgloss.Color(conf.ColorBlur))

	focusedStyle = style.Copy().
			Foreground(lipgloss.Color(conf.ColorFocus))
)

type QueryView struct {
	path     *utils.PathParts
	textarea textarea.Model
	err      error
}

func NewQueryView(path *utils.PathParts) *QueryView {
	ta := textarea.New()
	ta.Placeholder = "select ..."

	ta.Prompt = lipgloss.ThickBorder().Left + " "

	ta.BlurredStyle = textarea.Style{
		Prompt: style,
	}
	ta.FocusedStyle = textarea.Style{
		Prompt: focusedStyle,
	}

	ta.SetWidth(10)
	ta.SetHeight(conf.SqlTextareaHeight)
	ta.Focus()

	return &QueryView{
		path:     path,
		textarea: ta,
		err:      nil,
	}
}
func (m *QueryView) Focus() tea.Cmd {
	return m.textarea.Focus()
}

func (m *QueryView) Blur() {
	m.textarea.Blur()
}

func (m *QueryView) SetValue(value string) {
	m.textarea.SetValue(value)
}

func (m *QueryView) Value() string {
	return m.textarea.Value()
}

func (m *QueryView) Init() tea.Cmd {
	return tea.Batch(textarea.Blink)
}
func (m *QueryView) Update(msg tea.Msg) (*QueryView, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.textarea.SetWidth(msg.Width)
	// We handle errors just like any other message
	case error:
		m.err = msg
		return m, nil
	}

	m.textarea, cmd = m.textarea.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m *QueryView) View() string {
	return m.textarea.View()
}
