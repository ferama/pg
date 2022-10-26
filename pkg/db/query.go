package db

import (
	"context"
	"fmt"

	"github.com/ferama/pg/pkg/pool"
	"github.com/jackc/pgx/v5/pgtype"
)

func castType(item any) string {
	value := item

	if c, ok := item.(pgtype.Numeric); ok {
		if tmp, err := c.Value(); err == nil {
			value = tmp
		}
	}
	return fmt.Sprint(value)
}

func Query(connString, dbName, query string) ([][]string, error) {
	conn, err := pool.GetPoolFromConf(connString, dbName)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %v", err)
	}
	defer conn.Close()

	rows, err := conn.Query(context.Background(), query)
	if err != nil {
		return nil, fmt.Errorf("queryRow failed: %v", err)
	}
	defer rows.Close()

	var out [][]string
	out = make([][]string, 0)

	for rows.Next() {
		res, _ := rows.Values()
		var row []string
		row = make([]string, 0)

		for _, item := range res {
			if item == nil {
				row = append(row, "-")
			} else {
				row = append(row, castType(item))
			}
		}
		out = append(out, row)
	}
	return out, nil
}
