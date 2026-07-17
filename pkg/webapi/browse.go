package webapi

import (
	"fmt"
	"net/http"

	"github.com/ferama/pg/pkg/db"
	"github.com/ferama/pg/pkg/utils"
)

func queryDatabases(connString string) (*db.QueryResults, error) {
	query := `
		SELECT
			d.datname as database,
			r.rolname as owner,
			numbackends as "active connections"
		FROM pg_database d
		LEFT JOIN pg_roles r
			ON (d.datdba = r.oid)
		LEFT JOIN pg_stat_database  s
			ON s.datname = d.datname
		WHERE d.datistemplate = false
		ORDER BY d.datname
	`
	return db.Query(connString, "", "", query)
}

func querySchemas(connString, dbName string) (*db.QueryResults, error) {
	query := `
		SELECT
			nspname AS schema,
			usename AS owner
  		FROM
			pg_namespace
  		JOIN
			pg_user ON nspowner = usesysid
   		WHERE
			nspname NOT LIKE 'pg_%' AND
			nspname NOT LIKE 'information_schema' AND
			nspname NOT LIKE 'tiger%'
		ORDER BY schema
	`
	return db.Query(connString, dbName, "", query)
}

func queryTables(connString, dbName, schema string, details bool) (*db.QueryResults, error) {
	filter := `
		AND c.relkind IN ('r', 'v', 'p', 'm')
	`
	if details {
		filter = ""
	}

	query := fmt.Sprintf(`
		SELECT
			c.relname AS name,
			CASE c.relkind
				WHEN 'r' THEN 'TABLE'
				WHEN 'i' THEN 'SECONDARY INDEX'
				WHEN 'S' THEN 'SEQUENCE'
				WHEN 'v' THEN 'VIEW'
				WHEN 'm' THEN 'MATERIALIZED VIEW'
				WHEN 'f' THEN 'FOREIGN TABLE'
				WHEN 'p' THEN 'PARTITIONED TABLE'
				WHEN 'I' THEN 'PARTITIONED INDEX'
				WHEN 't' THEN 'OUT OF LINE VALUE'
			END AS type,
			pg_catalog.pg_get_userbyid(c.relowner) as owner,
			pg_catalog.pg_size_pretty(pg_catalog.pg_table_size(c.oid)) as size
 		FROM pg_catalog.pg_class c
 		LEFT JOIN pg_catalog.pg_namespace n
   			ON n.oid = c.relnamespace
 		WHERE n.nspname = '%s' %s
		ORDER BY name, type
	`, schema, filter)

	return db.Query(connString, dbName, "", query)
}

func queryTableColumns(connString, dbName, schema, tableName string) (*db.QueryResults, error) {
	query := fmt.Sprintf(`
		SELECT
			c.column_name as column,
			c.data_type as "data type",
			c.is_nullable as nullable,
			c.numeric_precision as "precision",
			c.character_maximum_length as "max length"
		FROM information_schema.columns AS c
		WHERE c.table_schema = '%s'
			AND c.table_name = '%s'
		ORDER BY c.column_name
		`, schema, tableName)

	return db.Query(connString, dbName, "", query)
}

func queryTableIndexes(connString, dbName, tableName string) (*db.QueryResults, error) {
	query := fmt.Sprintf(`
		SELECT indexname as index, indexdef as def
		FROM pg_indexes
		WHERE tablename = '%s'
		ORDER BY indexname
		`, tableName)

	return db.Query(connString, dbName, "", query)
}

func queryTableConstraints(connString, dbName, schema, tableName string) (*db.QueryResults, error) {
	query := fmt.Sprintf(`
		SELECT
			conname as "constraint name",
			pg_get_constraintdef(c.oid, true) as definition
		FROM
			pg_constraint c
		JOIN
			pg_namespace n ON n.oid = c.connamespace
		JOIN
			pg_class cl ON cl.oid = c.conrelid
		WHERE
			n.nspname = '%s' AND
			relname = '%s'
		ORDER BY
			contype desc`, schema, tableName)

	return db.Query(connString, dbName, "", query)
}

type tableDetailResponse struct {
	Kind        string       `json:"kind"`
	Columns     QueryResult  `json:"columns"`
	Indexes     *QueryResult `json:"indexes,omitempty"`
	Constraints *QueryResult `json:"constraints,omitempty"`
}

type browseResponse struct {
	Kind   string      `json:"kind"`
	Result QueryResult `json:"result"`
}

func handleBrowse(w http.ResponseWriter, r *http.Request) {
	path := utils.ParsePath(r.URL.Query().Get("path"), false)
	more := r.URL.Query().Get("more") == "true"

	if path.TableName != "" {
		colRes, err := queryTableColumns(path.ConfigConnection, path.DatabaseName, path.SchemaName, path.TableName)
		if err != nil {
			writeError(w, http.StatusInternalServerError, err)
			return
		}
		resp := tableDetailResponse{
			Kind:    "table",
			Columns: fromQueryResults(colRes),
		}
		if more {
			idxRes, err := queryTableIndexes(path.ConfigConnection, path.DatabaseName, path.TableName)
			if err != nil {
				writeError(w, http.StatusInternalServerError, err)
				return
			}
			idx := fromQueryResults(idxRes)
			resp.Indexes = &idx

			consRes, err := queryTableConstraints(path.ConfigConnection, path.DatabaseName, path.SchemaName, path.TableName)
			if err != nil {
				writeError(w, http.StatusInternalServerError, err)
				return
			}
			cons := fromQueryResults(consRes)
			resp.Constraints = &cons
		}
		writeJSON(w, http.StatusOK, resp)
		return
	}

	if path.SchemaName != "" {
		res, err := queryTables(path.ConfigConnection, path.DatabaseName, path.SchemaName, more)
		if err != nil {
			writeError(w, http.StatusInternalServerError, err)
			return
		}
		writeJSON(w, http.StatusOK, browseResponse{Kind: "tables", Result: fromQueryResults(res)})
		return
	}

	if path.DatabaseName != "" {
		res, err := querySchemas(path.ConfigConnection, path.DatabaseName)
		if err != nil {
			writeError(w, http.StatusInternalServerError, err)
			return
		}
		writeJSON(w, http.StatusOK, browseResponse{Kind: "schemas", Result: fromQueryResults(res)})
		return
	}

	if path.ConfigConnection != "" {
		res, err := queryDatabases(path.ConfigConnection)
		if err != nil {
			writeError(w, http.StatusInternalServerError, err)
			return
		}
		writeJSON(w, http.StatusOK, browseResponse{Kind: "databases", Result: fromQueryResults(res)})
		return
	}

	writeError(w, http.StatusBadRequest, fmt.Errorf("path is required"))
}
