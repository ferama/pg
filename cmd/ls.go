package cmd

import (
	"fmt"
	"os"

	"github.com/ferama/pg/pkg/autocomplete"
	"github.com/ferama/pg/pkg/components/table"
	"github.com/ferama/pg/pkg/conf"
	"github.com/ferama/pg/pkg/db"
	"github.com/ferama/pg/pkg/utils"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(lsCmd)

	lsCmd.Flags().BoolP("more", "m", false, "include more details")
}

func listConnections() {
	c := conf.GetAvailableConnections()

	t := table.NewStatic([]string{"CONNECTION"})
	var rs []table.SimpleRow
	for _, item := range c {
		rs = append(rs, table.SimpleRow{item})
	}
	t.SetRows(rs)
	fmt.Printf("\n%s\n\n", t.Render())
}

func listDatabases(connString string) {
	query := `
		SELECT 
			d.datname as database, 
			r.rolname as owner, 
			numbackends as "active connections"
		FROM pg_database d
		LEFT JOIN pg_roles r 
			ON (d.datdba = r.oid)
		LEFT JOIN pg_stat_database  s
			ON s.datname = d.datname
		WHERE d.datistemplate = false
		ORDER BY d.datname
	`

	results, err := db.Query(connString, "", "", query)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	db.PrintQueryResults(results)
}

func listSchemas(connString string, dbName string) {
	query := `
		SELECT schema_name as schema
		FROM information_schema.schemata
		ORDER BY schema_name
	`
	results, err := db.Query(connString, dbName, "", query)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	db.PrintQueryResults(results)
}

func listTables(connString string, dbName string, schema string, details bool) {
	query := fmt.Sprintf(`
		SELECT table_name as table, table_type as type
		FROM information_schema.tables
		WHERE table_schema = '%s' 
		ORDER BY table_name
		`, schema)

	results, err := db.Query(connString, dbName, "", query)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	db.PrintQueryResults(results)

	if details {
		query = fmt.Sprintf(`
			SELECT sequence_name as sequence
			FROM information_schema.sequences
			WHERE sequence_schema = '%s' 
			ORDER BY sequence_name
			`, schema)

		results, err = db.Query(connString, dbName, "", query)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		db.PrintQueryResults(results)
	}
}

func listTableDetails(connString, dbName, schema, tableName string, details bool) {
	// Columns
	query := fmt.Sprintf(`
		SELECT 
			c.column_name as column, 
			c.data_type as "data type", 
			c.is_nullable as nullable,
			c.numeric_precision as "precision", 
			c.character_maximum_length as "max length", 
			constraint_type as key,
			ccu.table_name as "key table"
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

	results, err := db.Query(connString, dbName, "", query)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	db.PrintQueryResults(results)

	if details {
		// Indexes
		query = fmt.Sprintf(`
			SELECT indexname as index, indexdef as def
			FROM pg_indexes
			WHERE tablename = '%s'
			ORDER BY indexname
			`, tableName)

		results, err = db.Query(connString, dbName, "", query)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		db.PrintQueryResults(results)

		// constraints
		query = fmt.Sprintf(`
			SELECT
				conname as "constraint name",
				pg_get_constraintdef(c.oid, true) as definition
			FROM
				pg_constraint c
			JOIN
				pg_namespace n ON n.oid = c.connamespace
			JOIN
				pg_class cl ON cl.oid = c.conrelid
			WHERE
				n.nspname = '%s' AND
				relname = '%s'
			ORDER BY
				contype desc`, schema, tableName)

		results, err = db.Query(connString, dbName, "", query)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		db.PrintQueryResults(results)
	}

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

		showDetails, _ := cmd.Flags().GetBool("details")

		path := utils.ParsePath(args[0], false)

		if path.TableName != "" {
			listTableDetails(
				path.ConfigConnection,
				path.DatabaseName,
				path.SchemaName,
				path.TableName,
				showDetails,
			)
			return
		}
		if path.SchemaName != "" {
			listTables(
				path.ConfigConnection,
				path.DatabaseName,
				path.SchemaName,
				showDetails)
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
