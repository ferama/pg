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

func Path(level int) func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		shellDirective := cobra.ShellCompDirectiveNoFileComp | cobra.ShellCompDirectiveNoSpace

		parts := strings.Split(toComplete, "/")

		basePath := ""
		suffix := ""
		if strings.HasSuffix(toComplete, "/") {
			basePath = strings.Join(parts, "/")
			basePath = strings.TrimSuffix(basePath, "/")
		} else {
			basePath = strings.Join(parts[:len(parts)-1], "/")
		}

		path := utils.ParsePath(toComplete)

		if len(parts)-1 >= level {
			return nil, shellDirective
		}

		if len(parts) == 1 {
			items := conf.GetAvailableConnections()
			for i := range items {
				items[i] = fmt.Sprintf("%s/", items[i])
			}
			return items, shellDirective
		}

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
			ret := flatArray(items)
			for i := range ret {
				ret[i] = fmt.Sprintf("%s/%s%s", basePath, ret[i], suffix)
			}
			return ret, shellDirective
		}

		if path.DatabaseName != "" && level > 2 {
			query := `
				SELECT schema_name 
				FROM information_schema.schemata
				ORDER BY schema_name
			`
			items, err := db.Query(path.ConfigConnection, path.DatabaseName, query)
			if err != nil {
				return nil, shellDirective
			}
			ret := flatArray(items)
			for i := range ret {
				ret[i] = fmt.Sprintf("%s/%s%s", basePath, ret[i], suffix)
			}
			return ret, shellDirective
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
			ret := flatArray(items)
			for i := range ret {
				ret[i] = fmt.Sprintf("%s/%s%s", basePath, ret[i], suffix)
			}
			return ret, shellDirective
		}

		return nil, shellDirective
	}
}
