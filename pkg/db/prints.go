package db

import (
	"fmt"
	"strings"

	"github.com/ferama/pg/pkg/components/table"
	"github.com/ferama/pg/pkg/conf"
)

func PrintQueryResults(results *QueryResults, truncate bool) {
	res := RenderQueryResults(results, truncate)
	fmt.Printf("\n%s\n\n", res)
}

func RenderQueryResults(results *QueryResults, truncate bool) string {

	var upper []string
	for _, c := range results.Columns {
		upper = append(upper, strings.ToUpper(c))
	}
	t := table.NewStatic(upper)
	var rs []table.SimpleRow

	for _, row := range results.Rows {
		var tr table.SimpleRow
		for _, item := range row {
			out := item
			if len(out) > conf.ItemMaxLen && truncate {
				out = out[:conf.ItemMaxLen] + "..."
			}
			tr = append(tr, out)
		}
		rs = append(rs, tr)
	}

	t.SetRows(rs)

	return t.Render()
}
