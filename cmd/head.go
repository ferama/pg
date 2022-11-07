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

func headGetWhereCondition(filters []string) (string, error) {
	validSplits := []string{
		"=",
		"!=",
		">",
		"<",
		"<=",
		">=",
	}

	where := ""
	if len(filters) > 0 {
		conditions := make([]string, 0)
		for _, f := range filters {
			t := ""
			for _, vs := range validSplits {
				f = strings.ToLower(f)
				parts := strings.Split(f, vs)
				if len(parts) != 2 {
					continue
				}
				t = fmt.Sprintf("%s%s'%s'", parts[0], vs, parts[1])
			}
			if t == "" {
				return "", fmt.Errorf("invalid filter: '%s'", f)
			}
			conditions = append(conditions, t)
		}
		where = strings.Join(conditions, " AND ")
		where = fmt.Sprintf("AND %s", where)
	}
	return where, nil
}

func headTable(
	connString, dbName, schema, tableName string,
	columns, filters []string,
) {

	cols := "*"
	if len(columns) > 0 {
		cols = strings.Join(columns, ",")
	}

	whereConditions, err := headGetWhereCondition(filters)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	query := fmt.Sprintf(`
		SELECT %s
		FROM %s
		WHERE true %s
		LIMIT 10
		`, cols, tableName, whereConditions)

	results, err := db.Query(connString, dbName, schema, query)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	db.PrintQueryResults(results)
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

		path := utils.ParsePath(args[0], false)

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
