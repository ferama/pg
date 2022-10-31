package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/ferama/pg/pkg/autocomplete"
	"github.com/ferama/pg/pkg/db"
	"github.com/ferama/pg/pkg/utils"
	"github.com/spf13/cobra"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

const textareaHeight = 6

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
	path     *utils.PathParts
	viewport viewport.Model
	textarea textarea.Model
	err      error
}

func newModel(path *utils.PathParts) *model {
	ti := textarea.New()
	vp := viewport.New(20, 10)
	ti.Placeholder = "select ..."
	ti.SetWidth(10)
	ti.SetHeight(textareaHeight)
	ti.Focus()

	return &model{
		path:     path,
		viewport: vp,
		textarea: ti,
		err:      nil,
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
		m.viewport.Width = msg.Width
		m.viewport.Height = msg.Height - (textareaHeight + 2)
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEsc:
			m.textarea.Reset()
			return m, tea.Quit
		case tea.KeyCtrlX:
			m.viewport.SetContent("running query...")
			go func() {
				query := m.textarea.Value()
				res, err := sqlExecute(
					m.path.ConfigConnection,
					m.path.DatabaseName,
					m.path.SchemaName,
					query,
				)
				if err != nil {
					m.viewport.SetContent(err.Error())
				} else {
					m.viewport.SetContent(res)
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
	return fmt.Sprintf(
		"%s\n%s\n\n%s",
		m.viewport.View(),
		m.textarea.View(),
		"(ctrl+x to execute. ESC to cancel)",
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
