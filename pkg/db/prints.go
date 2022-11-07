package db

import (
	"fmt"

	"github.com/ferama/pg/pkg/utils"
	"github.com/jedib0t/go-pretty/v6/table"
)

func PrintQueryResults(results *QueryResults) {
	res := RenderQueryResults(results)
	fmt.Printf("\n%s\n\n", res)
}

func RenderQueryResults(results *QueryResults) string {
	t := utils.GetTableWriter()

	var tr table.Row
	for _, f := range results.Columns {
		tr = append(tr, f)
	}
	t.AppendHeader(tr)

	for _, row := range results.Columns {
		var tr table.Row
		for _, item := range row {
			tr = append(tr, item)
		}
		t.AppendRow(tr)
	}

	return t.Render()
}
