package webapi

import (
	"fmt"
	"net/http"

	"github.com/ferama/pg/pkg/db"
)

type statsResponse struct {
	MaxConnections     string `json:"maxConnections"`
	CurrentConnections string `json:"currentConnections"`
}

func handleStats(w http.ResponseWriter, r *http.Request) {
	connName := r.URL.Query().Get("conn")

	conn, err := db.GetDBFromConf(connName, "")
	if err != nil {
		writeError(w, http.StatusInternalServerError, fmt.Errorf("unable to connect to database: %v", err))
		return
	}
	defer conn.Close()

	rows1, err := conn.Query("show max_connections")
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	rows1.Next()
	var maxConnections string
	if err := rows1.Scan(&maxConnections); err != nil {
		rows1.Close()
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	rows1.Close()

	rows2, err := conn.Query("SELECT sum(numbackends) as current_connections FROM pg_stat_database")
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	rows2.Next()
	var currentConnections string
	if err := rows2.Scan(&currentConnections); err != nil {
		rows2.Close()
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	rows2.Close()

	writeJSON(w, http.StatusOK, statsResponse{
		MaxConnections:     maxConnections,
		CurrentConnections: currentConnections,
	})
}
