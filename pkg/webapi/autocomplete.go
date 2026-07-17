package webapi

import (
	"fmt"
	"net/http"

	"github.com/ferama/pg/pkg/db"
	"github.com/ferama/pg/pkg/utils"
)

// sqlKeywords mirrors pkg/sqlview/components/intellisense.sqlKeywords.
var sqlKeywords = []string{
	"ALL", "ALTER", "AND", "AS", "ASC", "AVG",
	"BETWEEN", "BY",
	"CASE", "CAST", "COALESCE", "COUNT", "CREATE", "CROSS",
	"DELETE", "DESC", "DISTINCT", "DROP",
	"ELSE", "END", "EXISTS",
	"FULL", "FROM",
	"GROUP",
	"HAVING",
	"ILIKE", "IN", "INNER", "INSERT", "INTO", "IS",
	"JOIN",
	"LEFT", "LIKE", "LIMIT",
	"MAX", "MIN",
	"NOT", "NULL", "NULLIF",
	"OFFSET", "ON", "OR", "ORDER", "OUTER",
	"RETURNING", "RIGHT",
	"SELECT", "SET", "SUM",
	"TABLE", "THEN",
	"UNION", "UPDATE",
	"VALUES", "VIEW",
	"WHEN", "WHERE", "WITH",
}

func buildSchemaFilter(schemaName string) string {
	if schemaName != "" {
		return fmt.Sprintf("AND table_schema = '%s'", schemaName)
	}
	return "AND table_schema NOT IN ('pg_catalog', 'information_schema', 'pg_toast')"
}

type autocompleteResponse struct {
	Keywords []string `json:"keywords"`
	Tables   []string `json:"tables"`
	Columns  []string `json:"columns"`
}

func handleAutocomplete(w http.ResponseWriter, r *http.Request) {
	path := utils.ParsePath(r.URL.Query().Get("path"), false)

	resp := autocompleteResponse{
		Keywords: sqlKeywords,
		Tables:   []string{},
		Columns:  []string{},
	}

	if path.ConfigConnection == "" && path.DatabaseName == "" {
		writeJSON(w, http.StatusOK, resp)
		return
	}

	schemaFilter := buildSchemaFilter(path.SchemaName)

	tablesQuery := fmt.Sprintf(
		`SELECT table_name FROM information_schema.tables
		 WHERE table_type = 'BASE TABLE' %s
		 ORDER BY table_name`,
		schemaFilter)
	if res, err := db.Query(path.ConfigConnection, path.DatabaseName, path.SchemaName, tablesQuery); err == nil {
		for _, row := range res.Rows {
			if len(row) > 0 {
				resp.Tables = append(resp.Tables, row[0])
			}
		}
	}

	colQuery := fmt.Sprintf(
		`SELECT DISTINCT column_name FROM information_schema.columns
		 WHERE TRUE %s
		 ORDER BY column_name`,
		schemaFilter)
	if res, err := db.Query(path.ConfigConnection, path.DatabaseName, path.SchemaName, colQuery); err == nil {
		for _, row := range res.Rows {
			if len(row) > 0 {
				resp.Columns = append(resp.Columns, row[0])
			}
		}
	}

	writeJSON(w, http.StatusOK, resp)
}
