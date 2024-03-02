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
	db.PrintQueryResults(results, false)
}

func listSchemas(connString string, dbName string) {
	query := `
		SELECT 
			nspname AS schema,
			usename AS owner
  		FROM 
			pg_namespace
  		JOIN 
			pg_user ON nspowner = usesysid
   		WHERE 
			nspname NOT LIKE 'pg_%' AND
			nspname NOT LIKE 'information_schema' AND
			nspname NOT LIKE 'tiger%'
		ORDER BY schema
	`
	results, err := db.Query(connString, dbName, "", query)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	db.PrintQueryResults(results, false)
}

func listTables(connString string, dbName string, schema string, details bool) {

	filter := `
		AND c.relkind IN ('r', 'v', 'p', 'm')
	`
	if details {
		filter = ""
	}

	query := fmt.Sprintf(`
		SELECT 
			c.relname AS name,
			CASE c.relkind
				WHEN 'r' THEN 'TABLE'
				WHEN 'i' THEN 'SECONDARY INDEX'
				WHEN 'S' THEN 'SEQUENCE'
				WHEN 'v' THEN 'VIEW'
				WHEN 'm' THEN 'MATERIALIZED VIEW'
				WHEN 'f' THEN 'FOREIGN TABLE'
				WHEN 'p' THEN 'PARTITIONED TABLE'
				WHEN 'I' THEN 'PARTITIONED INDEX'
				WHEN 't' THEN 'OUT OF LINE VALUE'
			END AS type,
			pg_catalog.pg_get_userbyid(c.relowner) as owner,
			pg_catalog.pg_size_pretty(pg_catalog.pg_table_size(c.oid)) as size
 		FROM pg_catalog.pg_class c
 		LEFT JOIN pg_catalog.pg_namespace n
   			ON n.oid = c.relnamespace
 		WHERE n.nspname = '%s' %s
		ORDER BY name, type
	`, schema, filter)

	results, err := db.Query(connString, dbName, "", query)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	db.PrintQueryResults(results, false)
}

func listTableDetails(connString, dbName, schema, tableName string, details bool) {
	// Columns
	query := fmt.Sprintf(`
		SELECT 
			c.column_name as column, 
			c.data_type as "data type", 
			c.is_nullable as nullable,
			c.numeric_precision as "precision", 
			c.character_maximum_length as "max length"
		FROM information_schema.columns AS c 
		WHERE c.table_schema = '%s'
			AND c.table_name = '%s' 
		ORDER BY c.column_name
		`, schema, tableName)

	results, err := db.Query(connString, dbName, "", query)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	db.PrintQueryResults(results, false)

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
		db.PrintQueryResults(results, false)

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
		db.PrintQueryResults(results, false)
	}

}

var lsCmd = &cobra.Command{
	Use:     "ls [conn]",
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

		showDetails, _ := cmd.Flags().GetBool("more")

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
