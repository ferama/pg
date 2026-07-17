package webapi

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/ferama/pg/pkg/db"
	"github.com/ferama/pg/pkg/utils"
)

// searchGetWhereCondition mirrors cmd/search.go's searchGetWhereCondition.
func searchGetWhereCondition(filters []string) (string, error) {
	validSplits := []string{
		"=",
		"!=",
		">",
		"<",
		"<=",
		">=",
		"like",
		"ilike",
		"is null",
		"not null",
	}

	where := ""
	if len(filters) > 0 {
		conditions := make([]string, 0)
		for _, f := range filters {
			t := ""
			for _, vs := range validSplits {
				f = strings.ToLower(f)
				parts := strings.Split(f, vs)
				if len(parts) != 2 {
					continue
				}
				t = fmt.Sprintf("%s %s '%s'",
					strings.TrimSpace(parts[0]),
					vs,
					strings.TrimSpace(parts[1]))
			}
			if t == "" {
				return "", fmt.Errorf("invalid filter: '%s'", f)
			}
			conditions = append(conditions, t)
		}
		where = strings.Join(conditions, " AND ")
		where = fmt.Sprintf("AND %s", where)
	}
	return where, nil
}

// splitCommaValues expands repeated query params that may also contain
// comma-separated values, matching cobra's StringSliceP flag behavior.
func splitCommaValues(values []string) []string {
	out := make([]string, 0, len(values))
	for _, v := range values {
		for _, part := range strings.Split(v, ",") {
			part = strings.TrimSpace(part)
			if part != "" {
				out = append(out, part)
			}
		}
	}
	return out
}

func handleSearch(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	path := utils.ParsePath(q.Get("path"), false)
	if path.TableName == "" {
		writeError(w, http.StatusBadRequest, fmt.Errorf("table '%s' not found", path.TableName))
		return
	}

	columns := splitCommaValues(q["columns"])
	filters := splitCommaValues(q["filters"])

	limit := 10
	if l := q.Get("limit"); l != "" {
		parsed, err := strconv.Atoi(l)
		if err != nil {
			writeError(w, http.StatusBadRequest, fmt.Errorf("invalid limit: %v", err))
			return
		}
		limit = parsed
	}

	cols := "*"
	if len(columns) > 0 {
		cols = strings.Join(columns, ",")
	}

	whereConditions, err := searchGetWhereCondition(filters)
	if err != nil {
		writeError(w, http.StatusUnprocessableEntity, err)
		return
	}

	query := fmt.Sprintf(`
		SELECT %s
		FROM %s
		WHERE true %s
		LIMIT %d
		`, cols, path.TableName, whereConditions, limit)

	results, err := db.Query(path.ConfigConnection, path.DatabaseName, path.SchemaName, query)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusOK, fromQueryResults(results))
}
