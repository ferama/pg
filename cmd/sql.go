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
	tea "github.com/charmbracelet/bubbletea"
)

func init() {
	rootCmd.AddCommand(sqlCmd)
}

func sqlExecute(connString, dbName, schema, query string) {
	fields, items, err := db.Query(connString, dbName, schema, query)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	db.PrintQueryResults(items, fields)
}

type model struct {
	textarea textarea.Model
	err      error
}

func newModel() *model {
	ti := textarea.New()
	ti.Placeholder = "select ..."
	ti.SetWidth(60)
	ti.SetHeight(12)
	ti.Focus()

	return &model{
		textarea: ti,
		err:      nil,
	}
}
func (m *model) Init() tea.Cmd {
	return textarea.Blink
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEsc:
			m.textarea.Reset()
			return m, tea.Quit
		case tea.KeyCtrlX:
			return m, tea.Quit
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
		"\n%s\n\n%s",
		m.textarea.View(),
		"(ctrl+x to execute. ESC to cancel)",
	) + "\n\n"
}

var sqlCmd = &cobra.Command{
	Use:               "sql",
	Args:              cobra.MinimumNArgs(1),
	Short:             "Run sql query",
	ValidArgsFunction: autocomplete.Path(3),
	Run: func(cmd *cobra.Command, args []string) {
		path := utils.ParsePath(args[0], false)

		if path.SchemaName != "" || path.DatabaseName != "" {
			model := newModel()
			p := tea.NewProgram(model)

			if err := p.Start(); err != nil {
				log.Fatal(err)
			}
			query := model.textarea.Value()
			if query != "" {
				sqlExecute(
					path.ConfigConnection,
					path.DatabaseName,
					path.SchemaName,
					query,
				)
			}
		} else {
			fmt.Fprintf(os.Stderr, "database and schema not provided")
			os.Exit(1)
		}
	},
}
