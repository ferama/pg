package webapi

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ferama/pg/pkg/db"
	"github.com/ferama/pg/pkg/utils"
)

type chownRequest struct {
	Path  string `json:"path"`
	Owner string `json:"owner"`
}

type chownResponse struct {
	Messages []string `json:"messages"`
}

// handleChown mirrors cmd/chown.go: the database and schema branches are
// independent (not mutually exclusive), so a path with both a database and
// a schema segment runs both ALTER statements.
func handleChown(w http.ResponseWriter, r *http.Request) {
	var req chownRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	path := utils.ParsePath(req.Path, false)

	if path.DatabaseName == "" && path.SchemaName == "" {
		writeError(w, http.StatusBadRequest, fmt.Errorf("path must contain a database or schema"))
		return
	}

	resp := chownResponse{Messages: []string{}}

	if path.DatabaseName != "" {
		query := fmt.Sprintf(`
			ALTER DATABASE %s
			OWNER to %s
			`, path.DatabaseName, req.Owner)

		conn, err := db.GetDBFromConf(path.ConfigConnection, path.DatabaseName)
		if err != nil {
			writeError(w, http.StatusInternalServerError, fmt.Errorf("unable to connect to database: %v", err))
			return
		}
		defer conn.Close()
		if _, err = conn.Exec(query); err != nil {
			writeError(w, http.StatusInternalServerError, err)
			return
		}
		resp.Messages = append(resp.Messages, "database owner set")
	}

	if path.SchemaName != "" {
		query := fmt.Sprintf(`
			ALTER SCHEMA %s
			OWNER to %s
			`, path.SchemaName, req.Owner)

		conn, err := db.GetDBFromConf(path.ConfigConnection, path.DatabaseName)
		if err != nil {
			writeError(w, http.StatusInternalServerError, fmt.Errorf("unable to connect to database: %v", err))
			return
		}
		defer conn.Close()
		if _, err = conn.Exec(query); err != nil {
			writeError(w, http.StatusInternalServerError, err)
			return
		}
		resp.Messages = append(resp.Messages, "schema owner set")
	}

	writeJSON(w, http.StatusOK, resp)
}
