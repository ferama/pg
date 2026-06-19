package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/ferama/pg/pkg/autocomplete"
	"github.com/ferama/pg/pkg/db"
	"github.com/ferama/pg/pkg/sqlview"
	"github.com/ferama/pg/pkg/utils"
	"github.com/spf13/cobra"

	tea "github.com/charmbracelet/bubbletea"
)

func init() {
	rootCmd.AddCommand(sqlCmd)

	sqlCmd.Flags().StringP("query", "q", "", "run sql query")
}

func runCommand(connString, dbName, schema, query string) {

	results, err := db.Query(connString, dbName, schema, query)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	db.PrintQueryResults(results, true)
}

var sqlCmd = &cobra.Command{
	Use:   "sql",
	Args:  cobra.MinimumNArgs(1),
	Short: "Run sql query",
	Example: `  # Open the interactive SQL editor for a given schema
  pg sql myconn/mydb/public

  # Open the editor starting from a specific table (pre-fills a SELECT)
  pg sql myconn/mydb/public/mytable

  # Run a one-shot query without opening the editor
  pg sql myconn/mydb/public -q "SELECT * FROM mytable LIMIT 10"`,
	ValidArgsFunction: autocomplete.Path(4),
	Run: func(cmd *cobra.Command, args []string) {
		path := utils.ParsePath(args[0], false)

		query, _ := cmd.Flags().GetString("query")

		if path.DatabaseName != "" {
			if query == "" {
				model := sqlview.NewMainView(path)
				p := tea.NewProgram(model)

				if _, err := p.Run(); err != nil {
					log.Fatal(err)
				}
			} else {
				runCommand(
					path.ConfigConnection,
					path.DatabaseName,
					path.SchemaName,
					query,
				)
			}

		} else {
			fmt.Fprintf(os.Stderr, "database not provided")
			os.Exit(1)
		}
	},
}
