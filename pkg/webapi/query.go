package webapi

import (
	"encoding/json"
	"net/http"

	"github.com/ferama/pg/pkg/db"
	"github.com/ferama/pg/pkg/history"
	"github.com/ferama/pg/pkg/utils"
)

type queryRequest struct {
	Path  string `json:"path"`
	Query string `json:"query"`
}

func handleQuery(w http.ResponseWriter, r *http.Request) {
	var req queryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	path := utils.ParsePath(req.Path, false)

	// Append to history before execution, same order as the TUI editor.
	history.GetInstance().Append(req.Query)

	results, err := db.Query(path.ConfigConnection, path.DatabaseName, path.SchemaName, req.Query)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	out := fromQueryResults(results)
	if len(out.Columns) == 0 {
		out.Message = "done"
	}
	writeJSON(w, http.StatusOK, out)
}
