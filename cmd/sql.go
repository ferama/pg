package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/ferama/pg/pkg/autocomplete"
	"github.com/ferama/pg/pkg/db"
	"github.com/ferama/pg/pkg/utils"
	"github.com/spf13/cobra"
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

var sqlCmd = &cobra.Command{
	Use:               "sql",
	Args:              cobra.MinimumNArgs(2),
	Short:             "Run sql query",
	ValidArgsFunction: autocomplete.Path(3),
	Run: func(cmd *cobra.Command, args []string) {
		path := utils.ParsePath(strings.Join(args, " "), false)

		if path.SchemaName != "" {
			sqlExecute(
				path.ConfigConnection,
				path.DatabaseName,
				path.SchemaName,
				args[1],
			)
		} else {
			fmt.Fprintf(os.Stderr, "schema '%s' not found", path.SchemaName)
			os.Exit(1)
		}
	},
}
