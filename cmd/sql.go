package cmd

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/ferama/pg/pkg/autocomplete"
	"github.com/ferama/pg/pkg/db"
	"github.com/ferama/pg/pkg/utils"
	"github.com/spf13/cobra"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	sqlTextareaHeight = 6
)

var (
	resultsModelStyle = lipgloss.NewStyle().
				Border(lipgloss.ThickBorder(), true, true, false, true).
				BorderForeground(lipgloss.Color("#ffffff"))

	infoStyle = lipgloss.NewStyle()
)

func init() {
	rootCmd.AddCommand(sqlCmd)
}

func sqlExecute(connString, dbName, schema, query string) (string, error) {
	fields, items, err := db.Query(connString, dbName, schema, query)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	return db.RenderQueryResults(items, fields), nil
}

type model struct {
	path        *utils.PathParts
	resultsView viewport.Model
	textarea    textarea.Model
	err         error
}

func newModel(path *utils.PathParts) *model {
	ti := textarea.New()
	ti.Placeholder = "select ..."
	ti.SetWidth(10)
	ti.SetHeight(sqlTextareaHeight)
	ti.Focus()

	vp := viewport.New(20, 10)

	return &model{
		path:        path,
		resultsView: vp,
		textarea:    ti,
		err:         nil,
	}
}
func (m *model) Init() tea.Cmd {
	return tea.Batch(textarea.Blink, tea.EnterAltScreen)
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.textarea.SetWidth(msg.Width)
		resultsModelStyle.Width(msg.Width - 2)
		resultsModelStyle.Height(msg.Height - (sqlTextareaHeight + 3))
		m.resultsView.Width = msg.Width - 2
		m.resultsView.Height = msg.Height - (sqlTextareaHeight + 3)
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEsc:
			m.textarea.Reset()
			return m, tea.Quit
		case tea.KeyCtrlDown, tea.KeyCtrlD:
			m.resultsView.HalfViewDown()
		case tea.KeyCtrlUp, tea.KeyCtrlU:
			m.resultsView.HalfViewUp()
		case tea.KeyCtrlX:
			m.resultsView.SetContent("running query...")
			go func() {
				query := m.textarea.Value()
				res, err := sqlExecute(
					m.path.ConfigConnection,
					m.path.DatabaseName,
					m.path.SchemaName,
					query,
				)
				if err != nil {
					m.resultsView.SetContent(err.Error())
				} else {
					m.resultsView.SetContent(res)
				}
			}()
		default:
			if !m.textarea.Focused() {
				cmd = m.textarea.Focus()
				cmds = append(cmds, cmd)
			}
		}

	// We handle errors just like any other message
	case error:
		m.err = msg
		return m, nil
	}

	m.textarea, cmd = m.textarea.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m *model) View() string {
	renderedResults := resultsModelStyle.Render(m.resultsView.View())
	percent := fmt.Sprintf("%3.f%%", m.resultsView.ScrollPercent()*100)
	line := strings.Repeat("━", m.resultsView.Width-4)
	line = fmt.Sprintf("┗%s%s┛", line, infoStyle.Render(percent))
	out := lipgloss.JoinVertical(lipgloss.Center, renderedResults, line)
	return fmt.Sprintf(
		"%s\n%s\n%s",
		// renderedResults,
		out,
		m.textarea.View(),
		"     |ctrl+x| execute |ESC| cancel |ctrl+down| scroll down |ctrl+up| scroll up",
	)
}

var sqlCmd = &cobra.Command{
	Use:               "sql",
	Args:              cobra.MinimumNArgs(1),
	Short:             "Run sql query",
	ValidArgsFunction: autocomplete.Path(3),
	Run: func(cmd *cobra.Command, args []string) {
		path := utils.ParsePath(args[0], false)

		if path.SchemaName != "" || path.DatabaseName != "" {
			model := newModel(path)
			p := tea.NewProgram(model)

			if err := p.Start(); err != nil {
				log.Fatal(err)
			}

		} else {
			fmt.Fprintf(os.Stderr, "database and schema not provided")
			os.Exit(1)
		}
	},
}
