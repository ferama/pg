package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/ferama/pg/pkg/autocomplete"
	"github.com/ferama/pg/pkg/sqlview"
	"github.com/ferama/pg/pkg/utils"
	"github.com/spf13/cobra"

	tea "github.com/charmbracelet/bubbletea"
)

func init() {
	rootCmd.AddCommand(sqlCmd)
}

var sqlCmd = &cobra.Command{
	Use:               "sql",
	Args:              cobra.MinimumNArgs(1),
	Short:             "Run sql query",
	ValidArgsFunction: autocomplete.Path(3),
	Run: func(cmd *cobra.Command, args []string) {
		path := utils.ParsePath(args[0], false)

		if path.SchemaName != "" || path.DatabaseName != "" {
			model := sqlview.NewMainView(path)
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
