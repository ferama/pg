package db

import (
	"context"
	"errors"
	"fmt"

	"github.com/ferama/pg/pkg/pool"
	"github.com/jackc/pgx/v5/pgtype"
)

type ResultsFields []string

type ResultsRows [][]string

func castType(item any) string {
	value := item

	if c, ok := item.(pgtype.Numeric); ok {
		if tmp, err := c.Value(); err == nil {
			value = tmp
		}
	}
	return fmt.Sprint(value)
}

func Query(connString, dbName, schema, query string) (ResultsFields, ResultsRows, error) {
	if query == "" {
		return nil, nil, errors.New("query is empty")
	}
	conn, err := pool.GetPoolFromConf(connString, dbName)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to connect to database: %v", err)
	}
	defer conn.Close()

	if schema != "" {
		_, err := conn.Exec(context.Background(), fmt.Sprintf("set search_path to %s", schema))
		if err != nil {
			return nil, nil, fmt.Errorf("failed to select schema: %v", err)
		}
	}

	rows, err := conn.Query(context.Background(), query)
	if err != nil {
		// return nil, nil, fmt.Errorf("queryRow2 failed: %s", err.Error())
		return nil, nil, err
	}
	defer rows.Close()

	var out ResultsRows
	out = make(ResultsRows, 0)

	fields := make(ResultsFields, 0)
	for _, f := range rows.FieldDescriptions() {
		fields = append(fields, f.Name)
	}

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
	return fields, out, nil
}
