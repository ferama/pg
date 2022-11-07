package db

import (
	"errors"
	"fmt"
	"strings"

	"github.com/ferama/pg/pkg/pool"
)

type Columns []string
type Rows [][]string

type QueryResults struct {
	Columns Columns
	Rows    Rows
}

// Returns a clean query without any comment statements
func cleanQuery(query string) string {
	lines := []string{}

	for _, line := range strings.Split(query, "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "--") {
			continue
		}
		lines = append(lines, line)
	}

	return strings.TrimSpace(strings.Join(lines, "\n"))
}

func Query(connString, dbName, schema, query string) (*QueryResults, error) {
	if query == "" {
		return nil, errors.New("query is empty")
	}
	query = cleanQuery(query)

	conn, err := pool.GetPoolFromConf(connString, dbName)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %v", err)
	}
	defer conn.Close()

	if schema != "" {
		_, err := conn.Exec(fmt.Sprintf("set search_path to %s", schema))
		if err != nil {
			return nil, fmt.Errorf("failed to select schema: %v", err)
		}
	}

	rows, err := conn.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out Rows
	out = make(Rows, 0)

	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var row []string
		row = make([]string, 0)

		values := make([]any, len(columns))
		for i := range values {
			values[i] = new(any)
		}
		err = rows.Scan(values...)
		if err != nil {
			return nil, err
		}

		for i := range columns {
			values[i] = *(values[i].(*any))
		}

		for _, item := range values {
			if item == nil {
				row = append(row, "-")
			} else {
				row = append(row, fmt.Sprint(item))
			}
		}
		out = append(out, row)
	}

	return &QueryResults{
		Columns: columns,
		Rows:    out,
	}, nil
}
