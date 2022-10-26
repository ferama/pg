package cmd

import (
	"fmt"
	"os"

	"github.com/ferama/pg/pkg/db"
	"github.com/ferama/pg/pkg/utils"
	"github.com/spf13/cobra"
)

func headTable(connString, dbName, schema, tableName string) {

	fields, _ := db.GetTableColumns(connString, dbName, schema, tableName)

	query := fmt.Sprintf(`
		SELECT *
		FROM %s.%s
		LIMIT 10
		`, schema, tableName)

	err := db.PrintQueryResults(connString, dbName, query, fields)
	if err != nil {
		fmt.Println(err)
	}
}

func init() {
	rootCmd.AddCommand(headCmd)
}

var headCmd = &cobra.Command{
	Use:   "head",
	Args:  cobra.MinimumNArgs(1),
	Short: "Display first table records",
	Run: func(cmd *cobra.Command, args []string) {

		path := utils.ParsePath(args[0])
		if path.TableName != "" {
			headTable(
				path.ConfigConnection,
				path.DatabaseName,
				path.SchemaName,
				path.TableName,
			)
		} else {
			fmt.Fprintln(os.Stderr, "table name not found")
			os.Exit(1)
		}
	},
}
