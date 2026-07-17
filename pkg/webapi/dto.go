package webapi

import "github.com/ferama/pg/pkg/db"

// QueryResult is the uniform envelope for tabular data returned by the API.
type QueryResult struct {
	Columns   []string   `json:"columns"`
	Rows      [][]string `json:"rows"`
	ElapsedMs int64      `json:"elapsedMs"`
	Message   string     `json:"message,omitempty"`
}

func fromQueryResults(r *db.QueryResults) QueryResult {
	rows := make([][]string, len(r.Rows))
	for i, row := range r.Rows {
		rows[i] = []string(row)
	}
	cols := []string(r.Columns)
	if cols == nil {
		cols = []string{}
	}
	if rows == nil {
		rows = [][]string{}
	}
	return QueryResult{
		Columns:   cols,
		Rows:      rows,
		ElapsedMs: r.Elapsed.Milliseconds(),
	}
}
