package autocomplete

import (
	"fmt"
	"strings"

	"github.com/ferama/pg/pkg/conf"
	"github.com/ferama/pg/pkg/db"
	"github.com/ferama/pg/pkg/utils"
	"github.com/spf13/cobra"
)

func flatArray(items [][]string) []string {
	ret := make([]string, 0)
	for _, item := range items {
		ret = append(ret, fmt.Sprintf("%s/", item[0]))
	}
	return ret
}

// https://github.com/spf13/cobra/blob/main/shell_completions.md
func Path(level int) func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		shellDirective := cobra.ShellCompDirectiveNoFileComp

		if len(args) >= level {
			return nil, shellDirective
		}

		if len(args) == 0 {
			items := conf.GetAvailableConnections()
			for i := range items {
				items[i] = fmt.Sprintf("%s/", items[i])
			}
			return items, shellDirective
		}

		path := utils.ParsePath(strings.Join(args, " "))

		if path.SchemaName != "" && level > 3 {
			query := fmt.Sprintf(`
			SELECT table_name 
			FROM information_schema.tables
			WHERE table_schema = '%s' 
			ORDER BY table_name
			`, path.SchemaName)
			items, err := db.Query(path.ConfigConnection, path.DatabaseName, query)
			if err != nil {
				return nil, shellDirective
			}
			return flatArray(items), shellDirective
		}

		if path.DatabaseName != "" && level > 2 {
			query := `
				SELECT schema_name 
				FROM information_schema.schemata
				ORDER BY schema_name
			`
			items, err := db.Query(path.ConfigConnection, "", query)
			if err != nil {
				return nil, shellDirective
			}
			return flatArray(items), shellDirective
		}

		if path.ConfigConnection != "" && level > 1 {
			query := `
				SELECT d.datname
				FROM pg_database d
				WHERE d.datistemplate = false
				ORDER BY d.datname
				`
			items, err := db.Query(path.ConfigConnection, "", query)
			if err != nil {
				return nil, shellDirective
			}
			return flatArray(items), shellDirective
		}

		return nil, shellDirective
	}
}
