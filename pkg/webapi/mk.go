package webapi

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ferama/pg/pkg/db"
	"github.com/ferama/pg/pkg/utils"
)

type mkRequest struct {
	Path string `json:"path"`
}

type mkResponse struct {
	Messages []string `json:"messages"`
}

func handleMk(w http.ResponseWriter, r *http.Request) {
	var req mkRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	path := utils.ParsePath(req.Path, false)
	conn, err := db.GetDBFromConf(path.ConfigConnection, "")
	if err != nil {
		writeError(w, http.StatusInternalServerError, fmt.Errorf("unable to connect to database: %v", err))
		return
	}
	defer conn.Close()

	resp := mkResponse{Messages: []string{}}

	if path.DatabaseName != "" {
		query := fmt.Sprintf(`
			SELECT count(*) FROM pg_database WHERE datname = '%s'
			`, path.DatabaseName)

		rows1, err := conn.Query(query)
		if err != nil {
			writeError(w, http.StatusInternalServerError, err)
			return
		}

		rows1.Next()

		var exists int
		if err := rows1.Scan(&exists); err != nil {
			rows1.Close()
			writeError(w, http.StatusInternalServerError, err)
			return
		}
		rows1.Close()

		if exists == 0 {
			query := fmt.Sprintf(`
			create database "%s"
			`, path.DatabaseName)
			if _, err = conn.Exec(query); err != nil {
				writeError(w, http.StatusInternalServerError, err)
				return
			}
			resp.Messages = append(resp.Messages, fmt.Sprintf("created database '%s'", path.DatabaseName))
		}
	}

	if path.SchemaName != "" {
		dbConn, err := db.GetDBFromConf(path.ConfigConnection, path.DatabaseName)
		if err != nil {
			writeError(w, http.StatusInternalServerError, fmt.Errorf("unable to connect to database: %v", err))
			return
		}
		defer dbConn.Close()

		query := fmt.Sprintf(`
			create schema if not exists "%s"
			`, path.SchemaName)
		if _, err = dbConn.Exec(query); err != nil {
			writeError(w, http.StatusInternalServerError, err)
			return
		}
		resp.Messages = append(resp.Messages, fmt.Sprintf("created schema '%s'", path.SchemaName))
	}

	writeJSON(w, http.StatusOK, resp)
}
