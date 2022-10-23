package cmd

import (
	"context"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/ferama/pg/pkg/conf"
	"github.com/ferama/pg/pkg/pool"
	"github.com/ferama/pg/pkg/utils"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(recordCmd)
}

func listConnections() {
	c := conf.GetAvailableConnections()
	w := tabwriter.NewWriter(os.Stdout, 5, 5, 5, ' ', 0)
	fmt.Fprintln(w, "Connection")
	header := "---------"
	fmt.Fprintf(w, "%s\n", header)
	for _, item := range c {
		fmt.Fprintln(w, item)
	}
	w.Flush()
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

	for rows.Next() {
		var datname string
		err = rows.Scan(&datname)
		if err != nil {
			fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
			os.Exit(1)
		}
		fmt.Println(datname)
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

	for rows.Next() {
		var schema string
		err = rows.Scan(&schema)
		if err != nil {
			fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
			os.Exit(1)
		}
		fmt.Println(schema)
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

	for rows.Next() {
		var tableName string
		err = rows.Scan(&tableName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
			os.Exit(1)
		}
		fmt.Println(tableName)
	}
}

func listColumns(connString, db, schema, table string) {
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
		AND table_name = '%s' `, schema, table)
	rows, err := conn.Query(context.Background(), query)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		os.Exit(1)
	}
	defer rows.Close()

	w := tabwriter.NewWriter(os.Stdout, 5, 5, 5, ' ', 0)
	fmt.Fprintln(w, "Name\tType\tNullable")
	header := "---------"
	fmt.Fprintf(w, "%s\t%s\t%s\n", header, header, header)
	for rows.Next() {
		var columnName, dataType, isNullable string
		err = rows.Scan(&columnName, &dataType, &isNullable)
		if err != nil {
			fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
			os.Exit(1)
		}
		fmt.Fprintf(w, "%s\t%s\t%s\t\n", columnName, dataType, isNullable)
	}

	w.Flush()
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
