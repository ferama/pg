package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/ferama/pg/pkg/pool"
	"github.com/ferama/pg/pkg/utils"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
)

func headTable(connString, db, schema, tableName string) {
	conn, err := pool.GetFromConf(connString, db)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close()
	query := fmt.Sprintf(`
		SELECT column_name
		FROM information_schema.columns 
		WHERE table_schema = '%s'
		AND table_name = '%s' 
		`, schema, tableName)
	rows, err := conn.Query(context.Background(), query)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		os.Exit(1)
	}
	defer rows.Close()

	fields := make([]string, 0)
	for rows.Next() {
		var columnName string
		err = rows.Scan(&columnName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
			os.Exit(1)
		}
		fields = append(fields, columnName)
	}

	query = fmt.Sprintf(`
		SELECT *
		FROM %s.%s
		LIMIT 10
		`, schema, tableName)
	rows, err = conn.Query(context.Background(), query)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		os.Exit(1)
	}
	defer rows.Close()

	t := utils.GetTableWriter()
	var tr table.Row
	for _, f := range fields {
		tr = append(tr, f)
	}
	t.AppendHeader(tr)
	defer t.Render()

	for rows.Next() {
		pointers := make([]any, len(fields))
		container := make([]any, len(fields))

		for i := range pointers {
			pointers[i] = &container[i]
		}

		err = rows.Scan(pointers...)

		if err != nil {
			fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
			os.Exit(1)
		}
		var tr table.Row
		for i := range fields {
			tr = append(tr, container[i])
		}
		t.AppendRow(tr)
	}

}

func init() {
	rootCmd.AddCommand(headCmd)
}

var headCmd = &cobra.Command{
	Use:  "head",
	Args: cobra.MinimumNArgs(1),
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
