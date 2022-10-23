package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/ferama/pg/pkg/conf"
	"github.com/ferama/pg/pkg/pool"
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
	conn, err := pool.GetFromConf(connString, "")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close()

	query := `
		SELECT d.datname, u.usename
		FROM pg_database d
		JOIN pg_user u ON (d.datdba = u.usesysid)
		WHERE datistemplate = false
	`
	rows, err := conn.Query(context.Background(), query)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		os.Exit(1)
	}
	defer rows.Close()

	t := utils.GetTableWriter()
	t.AppendHeader(table.Row{"Database", "Owner"})
	defer t.Render()

	for rows.Next() {
		var datname, owner string
		err = rows.Scan(&datname, &owner)
		if err != nil {
			fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
			os.Exit(1)
		}
		t.AppendRow(table.Row{
			datname, owner,
		})
	}
}

func listSchemas(connString string, db string) {
	conn, err := pool.GetFromConf(connString, db)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close()

	query := `
		SELECT schema_name 
		FROM information_schema.schemata
	`
	rows, err := conn.Query(context.Background(), query)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		os.Exit(1)
	}
	defer rows.Close()

	t := utils.GetTableWriter()
	t.AppendHeader(table.Row{"Schema"})
	defer t.Render()

	for rows.Next() {
		var schema string
		err = rows.Scan(&schema)
		if err != nil {
			fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
			os.Exit(1)
		}
		t.AppendRow(table.Row{
			schema,
		})
	}
}

func listTables(connString string, db string, schema string) {
	conn, err := pool.GetFromConf(connString, db)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close()

	query := fmt.Sprintf(`
		SELECT table_name 
		FROM information_schema.tables
		WHERE table_schema = '%s' `, schema)
	rows, err := conn.Query(context.Background(), query)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		os.Exit(1)
	}
	defer rows.Close()

	t := utils.GetTableWriter()
	t.AppendHeader(table.Row{"Schema"})
	defer t.Render()

	for rows.Next() {
		var tableName string
		err = rows.Scan(&tableName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
			os.Exit(1)
		}
		t.AppendRow(table.Row{
			tableName,
		})
	}
}

func listColumns(connString, db, schema, tableName string) {
	conn, err := pool.GetFromConf(connString, db)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close()

	query := fmt.Sprintf(`
		SELECT column_name, data_type, is_nullable
		FROM information_schema.columns 
		WHERE table_schema = '%s'
		AND table_name = '%s' `, schema, tableName)
	rows, err := conn.Query(context.Background(), query)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		os.Exit(1)
	}
	defer rows.Close()

	t := utils.GetTableWriter()
	t.AppendHeader(table.Row{"Name", "Type", "Nullable"})
	defer t.Render()

	for rows.Next() {
		var columnName, dataType, isNullable string
		err = rows.Scan(&columnName, &dataType, &isNullable)
		if err != nil {
			fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
			os.Exit(1)
		}
		t.AppendRow(table.Row{
			columnName, dataType, isNullable,
		})
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
