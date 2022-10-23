package cmd

import (
	"context"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/ferama/gopigi/pkg/conf"
	"github.com/ferama/gopigi/pkg/pool"
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
		var datname string
		err = rows.Scan(&datname)
		if err != nil {
			fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
			os.Exit(1)
		}
		fmt.Println(datname)
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
		parts := strings.Split(args[0], "/")
		if parts[len(parts)-1] == "" {
			parts = parts[:len(parts)-1]
		}
		switch len(parts) {
		case 1:
			conn := parts[0]
			listDatabases(conn)
		case 2:
			conn := parts[0]
			database := parts[1]
			listSchemas(conn, database)
		}
	},
}
