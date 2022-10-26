package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/ferama/pg/pkg/autocomplete"
	"github.com/ferama/pg/pkg/db"
	"github.com/ferama/pg/pkg/utils"
	"github.com/spf13/cobra"
)

func headTable(
	connString, dbName, schema, tableName string,
	columns, filters []string,
) {

	fields, _ := db.GetTableColumns(connString, dbName, schema, tableName)

	cols := "*"
	if len(columns) > 0 {
		cols = strings.Join(columns, ",")
		fields = columns
	}
	where := ""
	if len(filters) > 0 {
		conditions := make([]string, 0)
		for _, f := range filters {
			parts := strings.Split(f, "=")
			t := fmt.Sprintf("%s='%s'", parts[0], parts[1])
			conditions = append(conditions, t)
		}
		where = strings.Join(conditions, " AND ")
		where = fmt.Sprintf("AND %s", where)
	}
	query := fmt.Sprintf(`
		SELECT %s
		FROM %s.%s
		WHERE true %s
		LIMIT 10
		`, cols, schema, tableName, where)

	items, err := db.Query(connString, dbName, query)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	db.PrintQueryResults(items, fields)
}

func init() {
	rootCmd.AddCommand(headCmd)

	headCmd.Flags().StringSliceP("columns", "c", nil, "include columns. empty to include all")
	headCmd.Flags().StringSliceP("filters", "f", nil, "filter column by value")
}

var headCmd = &cobra.Command{
	Use:   "head",
	Args:  cobra.MinimumNArgs(1),
	Short: "Display first table records",
	Example: `
  # get only some columns
  $ pg head myconn/testdb/public/sales -c age,sex,city
  # or
  $ pg head myconn/testdb/public/sales -c age -c sex
  # add where conditions
  $ pg head myconn/testdb/public/sales -f sex=M
	`,
	ValidArgsFunction: autocomplete.Path(4),
	Run: func(cmd *cobra.Command, args []string) {

		columns, _ := cmd.Flags().GetStringSlice("columns")
		filters, _ := cmd.Flags().GetStringSlice("filters")

		path := utils.ParsePath(strings.Join(args, " "))

		if path.TableName != "" {
			headTable(
				path.ConfigConnection,
				path.DatabaseName,
				path.SchemaName,
				path.TableName,
				columns,
				filters,
			)
		} else {
			fmt.Fprintf(os.Stderr, "table '%s' not found", path.TableName)
			os.Exit(1)
		}
	},
}
