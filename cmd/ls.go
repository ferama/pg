package cmd

import (
	"fmt"

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
	defer t.Render()

	for _, item := range c {
		t.AppendRow(table.Row{
			item,
		})
	}
}

func listDatabases(connString string) {
	query := `
		SELECT d.datname, u.usename
		FROM pg_database d
		JOIN pg_user u ON (d.datdba = u.usesysid)
		WHERE datistemplate = false
		ORDER BY d.datname
	`
	db.PrintQueryResults(connString, "", query, []string{"Database", "Owner"})
}

func listSchemas(connString string, dbName string) {
	query := `
		SELECT schema_name 
		FROM information_schema.schemata
		ORDER BY schema_name
	`
	err := db.PrintQueryResults(connString, dbName, query, []string{"Schema"})
	if err != nil {
		fmt.Println(err)
	}
}

func listTables(connString string, dbName string, schema string) {
	query := fmt.Sprintf(`
		SELECT table_name 
		FROM information_schema.tables
		WHERE table_schema = '%s' 
		ORDER BY table_name
		`, schema)
	err := db.PrintQueryResults(connString, dbName, query, []string{"Table"})
	if err != nil {
		fmt.Println(err)
	}
}

func listColumns(connString, dbName, schema, tableName string) {
	query := fmt.Sprintf(`
		SELECT column_name, data_type, is_nullable
		FROM information_schema.columns 
		WHERE table_schema = '%s'
		AND table_name = '%s' 
		ORDER BY column_name
		`, schema, tableName)
	err := db.PrintQueryResults(connString, dbName, query, []string{"Column", "Data Type", "Nullable"})
	if err != nil {
		fmt.Println(err)
	}
}

var lsCmd = &cobra.Command{
	Use:  "ls",
	Args: cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			listConnections()
			return
		}

		path := utils.ParsePath(args[0])
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
