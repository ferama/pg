package db

import (
	"context"
	"fmt"

	"github.com/ferama/pg/pkg/pool"
	"github.com/ferama/pg/pkg/utils"
	"github.com/jedib0t/go-pretty/v6/table"
)

func PrintQueryResults(connString, dbName, query string, fields []string) error {
	conn, err := pool.GetFromConf(connString, dbName)
	if err != nil {
		return fmt.Errorf("unable to connect to database: %v", err)
	}
	defer conn.Close()

	rows, err := conn.Query(context.Background(), query)
	if err != nil {
		return fmt.Errorf("queryRow failed: %v", err)
	}
	defer rows.Close()

	t := utils.GetTableWriter()
	var tr table.Row
	for _, f := range fields {
		tr = append(tr, f)
	}
	t.AppendHeader(tr)
	defer t.Render()

	for rows.Next() {
		pointers := make([]any, len(fields))
		container := make([]any, len(fields))

		for i := range pointers {
			pointers[i] = &container[i]
		}

		err = rows.Scan(pointers...)

		if err != nil {
			return fmt.Errorf("queryRow failed: %v", err)
		}
		var tr table.Row
		for i := range fields {
			tr = append(tr, container[i])
		}
		t.AppendRow(tr)
	}

	return nil
}
