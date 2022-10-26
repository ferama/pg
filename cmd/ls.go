package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/ferama/pg/pkg/autocomplete"
	"github.com/ferama/pg/pkg/conf"
	"github.com/ferama/pg/pkg/db"
	"github.com/ferama/pg/pkg/utils"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(lsCmd)
}

func listConnections() {
	c := conf.GetAvailableConnections()

	t := utils.GetTableWriter()
	t.AppendHeader(table.Row{"Connection"})
	defer func() {
		fmt.Println()
		t.Render()
		fmt.Println()
	}()

	for _, item := range c {
		t.AppendRow(table.Row{
			item,
		})
	}
}

func listDatabases(connString string) {
	query := `
		SELECT d.datname, r.rolname, numbackends
		FROM pg_database d
		LEFT JOIN pg_roles r 
			ON (d.datdba = r.oid)
		LEFT JOIN pg_stat_database  s
			ON s.datname = d.datname
		WHERE d.datistemplate = false
		ORDER BY d.datname
	`

	items, err := db.Query(connString, "", query)
	if err != nil {
		fmt.Println(err)
	}
	db.PrintQueryResults(items, []string{"Database", "Owner", "Active Connections"})
}

func listSchemas(connString string, dbName string) {
	query := `
		SELECT schema_name 
		FROM information_schema.schemata
		ORDER BY schema_name
	`
	items, err := db.Query(connString, dbName, query)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	db.PrintQueryResults(items, []string{"Schema"})
}

func listTables(connString string, dbName string, schema string) {
	query := fmt.Sprintf(`
		SELECT table_name 
		FROM information_schema.tables
		WHERE table_schema = '%s' 
		ORDER BY table_name
		`, schema)

	items, err := db.Query(connString, dbName, query)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	db.PrintQueryResults(items, []string{"Table"})
}

func listColumns(connString, dbName, schema, tableName string) {
	query := fmt.Sprintf(`
		SELECT 
			c.column_name, c.data_type, c.is_nullable,
			c.numeric_precision, c.character_maximum_length, constraint_type
		FROM information_schema.table_constraints tc 
		JOIN information_schema.constraint_column_usage AS ccu 
			USING (constraint_schema, constraint_name) 
		RIGHT JOIN information_schema.columns AS c 
			ON c.table_schema = tc.constraint_schema
			AND tc.table_name = c.table_name AND ccu.column_name = c.column_name
		WHERE c.table_schema = '%s'
			AND c.table_name = '%s' 
		ORDER BY c.column_name
		`, schema, tableName)

	items, err := db.Query(connString, dbName, query)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	db.PrintQueryResults(items, []string{
		"Column", "Data Type", "Nullable", "Numeric Precision", "Max Length", "Key"})
}

var lsCmd = &cobra.Command{
	Use:     "ls",
	Aliases: []string{"list"},
	Short:   "List configuration, database, schemas, tables and fields",
	Args:    cobra.MinimumNArgs(0),
	// https://github.com/spf13/cobra/blob/main/shell_completions.md
	ValidArgsFunction: autocomplete.Path(4),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			listConnections()
			return
		}

		path := utils.ParsePath(strings.Join(args, " "))

		if path.TableName != "" {
			listColumns(
				path.ConfigConnection,
				path.DatabaseName,
				path.SchemaName,
				path.TableName,
			)
			return
		}
		if path.SchemaName != "" {
			listTables(
				path.ConfigConnection,
				path.DatabaseName,
				path.SchemaName)
			return
		}
		if path.DatabaseName != "" {
			listSchemas(path.ConfigConnection, path.DatabaseName)
			return
		}
		if path.ConfigConnection != "" {
			listDatabases(path.ConfigConnection)
			return
		}
	},
}
