package db

import (
	"fmt"

	"github.com/ferama/pg/pkg/utils"
	"github.com/jedib0t/go-pretty/v6/table"
)

func PrintQueryResults(items [][]string, fields []string) {
	res := RenderQueryResults(items, fields)
	fmt.Printf("\n%s\n\n", res)
}

func RenderQueryResults(items [][]string, fields []string) string {
	t := utils.GetTableWriter()

	var tr table.Row
	for _, f := range fields {
		tr = append(tr, f)
	}
	t.AppendHeader(tr)

	for _, row := range items {
		var tr table.Row
		for _, item := range row {
			tr = append(tr, item)
		}
		t.AppendRow(tr)
	}

	return t.Render()
}
