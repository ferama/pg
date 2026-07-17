package webapi

import "net/http"

// NewHandler builds the JSON HTTP API exposed under /api/.
func NewHandler() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/connections", handleConnections)
	mux.HandleFunc("GET /api/ping", handlePing)

	mux.HandleFunc("GET /api/browse", handleBrowse)

	mux.HandleFunc("POST /api/mk", handleMk)
	mux.HandleFunc("POST /api/chown", handleChown)

	mux.HandleFunc("GET /api/search", handleSearch)

	mux.HandleFunc("POST /api/query", handleQuery)

	mux.HandleFunc("GET /api/history", handleHistoryList)
	mux.HandleFunc("DELETE /api/history/{idx}", handleHistoryDelete)

	mux.HandleFunc("GET /api/autocomplete", handleAutocomplete)

	mux.HandleFunc("GET /api/stats", handleStats)

	mux.HandleFunc("GET /api/users", handleUsersList)
	mux.HandleFunc("POST /api/users", handleUsersCreate)
	mux.HandleFunc("DELETE /api/users/{username}", handleUsersDelete)
	mux.HandleFunc("POST /api/users/{username}/attr", handleUsersAttr)

	return mux
}
