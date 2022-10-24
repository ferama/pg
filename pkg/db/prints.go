package db

import (
	"context"
	"fmt"

	"github.com/ferama/pg/pkg/pool"
	"github.com/ferama/pg/pkg/utils"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jedib0t/go-pretty/v6/table"
)

func castType(item any) any {
	value := item

	if c, ok := item.(pgtype.Numeric); ok {
		if tmp, err := c.Value(); err == nil {
			value = tmp
		}
	}

	// if c, ok := item.(pgtype.B); ok {
	// 	if tmp, err := c.Value(); err == nil {
	// 		value = tmp
	// 	}
	// }

	return value
}

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

	// https://github.com/GoAdminGroup/go-admin/blob/master/modules/db/performer.go
	for rows.Next() {
		res, _ := rows.Values()
		var tr table.Row

		for _, item := range res {
			tr = append(tr, castType(item))
		}
		t.AppendRow(tr)
	}

	return nil
}
