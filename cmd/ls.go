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
	rootCmd.AddCommand(recordCmd)
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

	rows, err := conn.Query(
		context.Background(),
		"select datname from pg_database where datistemplate = false")
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		os.Exit(1)
	}
	defer rows.Close()

	t := utils.GetTableWriter()
	t.AppendHeader(table.Row{"Database"})
	defer t.Render()

	for rows.Next() {
		var datname string
		err = rows.Scan(&datname)
		if err != nil {
			fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
			os.Exit(1)
		}
		t.AppendRow(table.Row{
			datname,
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

	rows, err := conn.Query(
		context.Background(),
		"select schema_name from information_schema.schemata")
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

var recordCmd = &cobra.Command{
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
