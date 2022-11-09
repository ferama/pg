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

func searchGetWhereCondition(filters []string) (string, error) {
	validSplits := []string{
		"=",
		"!=",
		">",
		"<",
		"<=",
		">=",
		"like",
		"ilike",
		"is null",
		"not null",
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
				t = fmt.Sprintf("%s %s '%s'",
					strings.TrimSpace(parts[0]),
					vs,
					strings.TrimSpace(parts[1]))
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

func searchTable(
	connString, dbName, schema, tableName string,
	columns, filters []string,
	limit int,
) {

	cols := "*"
	if len(columns) > 0 {
		cols = strings.Join(columns, ",")
	}

	whereConditions, err := searchGetWhereCondition(filters)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	query := fmt.Sprintf(`
		SELECT %s
		FROM %s
		WHERE true %s
		LIMIT %d
		`, cols, tableName, whereConditions, limit)

	results, err := db.Query(connString, dbName, schema, query)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	db.PrintQueryResults(results)
}

func init() {
	rootCmd.AddCommand(searchCmd)

	searchCmd.Flags().StringSliceP("columns", "c", nil, "include columns. empty to include all")
	searchCmd.Flags().StringSliceP("filters", "f", nil, "filter column by value")
	searchCmd.Flags().IntP("limit", "l", 10, "limit results")
}

var searchCmd = &cobra.Command{
	Use:   "search",
	Args:  cobra.MinimumNArgs(1),
	Short: "Search tables for records",
	Example: `
  # get only some columns
  $ pg search myconn/testdb/public/sales -c age,sex,city

  # or
  $ pg search myconn/testdb/public/sales -c age -c sex

  # add where conditions
  $ pg search myconn/testdb/public/sales -f sex=M

  # use like
  $ pg search myconn/testdb/public/sales -f "name like test%"
	`,
	ValidArgsFunction: autocomplete.Path(4),
	Run: func(cmd *cobra.Command, args []string) {

		columns, _ := cmd.Flags().GetStringSlice("columns")
		filters, _ := cmd.Flags().GetStringSlice("filters")
		limit, _ := cmd.Flags().GetInt("limit")

		path := utils.ParsePath(args[0], false)

		if path.TableName != "" {
			searchTable(
				path.ConfigConnection,
				path.DatabaseName,
				path.SchemaName,
				path.TableName,
				columns,
				filters,
				limit,
			)
		} else {
			fmt.Fprintf(os.Stderr, "table '%s' not found", path.TableName)
			os.Exit(1)
		}
	},
}
