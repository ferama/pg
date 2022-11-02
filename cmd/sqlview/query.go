package sqlview

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/ferama/pg/pkg/utils"
)

const (
	sqlTextareaHeight = 5
)

type queryView struct {
	path     *utils.PathParts
	textarea textarea.Model
	err      error
}

func newQueryView(path *utils.PathParts) *queryView {
	ti := textarea.New()
	ti.Placeholder = "select ..."
	ti.Prompt = ""
	ti.SetWidth(10)
	ti.SetHeight(sqlTextareaHeight)
	ti.Focus()

	return &queryView{
		path:     path,
		textarea: ti,
		err:      nil,
	}
}

func (m *queryView) Init() tea.Cmd {
	return tea.Batch(textarea.Blink)
}
func (m *queryView) Update(msg tea.Msg) (*queryView, tea.Cmd) {
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

func (m *queryView) View() string {
	return fmt.Sprint(m.textarea.View())
}
