package cmd

import (
	"fmt"
	"os"

	"github.com/ferama/pg/pkg/autocomplete"
	"github.com/ferama/pg/pkg/db"
	"github.com/ferama/pg/pkg/utils"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(mkCmd)
}

var mkCmd = &cobra.Command{
	Use:               "mk",
	Args:              cobra.MinimumNArgs(1),
	Short:             "Create database and/or schema",
	ValidArgsFunction: autocomplete.Path(2),
	Example: `
  # Create a database and a schema
  $ pg mk myconn/mydb/myschema
  `,
	Run: func(cmd *cobra.Command, args []string) {
		path := utils.ParsePath(args[0], false)
		conn, err := db.GetDBFromConf(path.ConfigConnection, "")
		if err != nil {
			fmt.Printf("unable to connect to database: %v", err)
			os.Exit(1)
		}
		defer conn.Close()

		if path.DatabaseName != "" {
			query := fmt.Sprintf(`
			SELECT count(*) FROM pg_database WHERE datname = '%s'
			`, path.DatabaseName)

			rows1, err := conn.Query(query)
			if err != nil {
				fmt.Printf("queryRow failed: %v", err)
				os.Exit(1)
			}
			defer rows1.Close()

			rows1.Next()

			var exists int
			err = rows1.Scan(&exists)
			if err != nil {
				fmt.Printf("queryRow failed: %v\n", err)
				os.Exit(1)
			}
			if exists == 0 {
				query := fmt.Sprintf(`
				create database "%s"
				`, path.DatabaseName)
				_, err = conn.Exec(query)
				if err != nil {
					fmt.Printf("error: %v", err)
					os.Exit(1)
				}
				fmt.Printf("created database '%s'\n", path.DatabaseName)
			}
		}

		if path.SchemaName != "" {
			conn, err := db.GetDBFromConf(path.ConfigConnection, path.DatabaseName)
			if err != nil {
				fmt.Printf("unable to connect to database: %v", err)
				os.Exit(1)
			}
			query := fmt.Sprintf(`
			create schema if not exists "%s"
			`, path.SchemaName)
			_, err = conn.Exec(query)
			if err != nil {
				fmt.Printf("error: %v", err)
				os.Exit(1)
			}
		}
	},
}
