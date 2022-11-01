package sqlview

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/ferama/pg/pkg/utils"
)

const (
	SqlTextareaHeight = 6
)

type QueryView struct {
	path     *utils.PathParts
	textarea textarea.Model
	err      error
}

func NewQueryView(path *utils.PathParts) *QueryView {
	ti := textarea.New()
	ti.Placeholder = "select ..."
	ti.Prompt = ""
	ti.SetWidth(10)
	ti.SetHeight(SqlTextareaHeight)
	ti.Focus()

	return &QueryView{
		path:     path,
		textarea: ti,
		err:      nil,
	}
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
	return fmt.Sprint(m.textarea.View())
}
