package webapi

import (
	"context"
	"net/http"
	"time"

	"github.com/ferama/pg/pkg/conf"
	"github.com/ferama/pg/pkg/db"
)

func handleConnections(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, conf.GetAvailableConnections())
}

type pingResponse struct {
	Connected bool `json:"connected"`
}

func handlePing(w http.ResponseWriter, r *http.Request) {
	connName := r.URL.Query().Get("conn")

	conn, err := db.GetDBFromConf(connName, "")
	if err != nil {
		writeJSON(w, http.StatusOK, pingResponse{Connected: false})
		return
	}
	defer conn.Close()

	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	if err := conn.PingContext(ctx); err != nil {
		writeJSON(w, http.StatusOK, pingResponse{Connected: false})
		return
	}
	writeJSON(w, http.StatusOK, pingResponse{Connected: true})
}
