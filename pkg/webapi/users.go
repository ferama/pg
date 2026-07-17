package webapi

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ferama/pg/pkg/db"
)

func handleUsersList(w http.ResponseWriter, r *http.Request) {
	connName := r.URL.Query().Get("conn")

	query := `
		SELECT
			USENAME as USERNAME,
			USECREATEDB as CREATEDB,
			USESUPER as ISSUPER,
			USEREPL as REPL,
			USEBYPASSRLS as BYPASSRLS,
			VALUNTIL as VALUNTIL,
			USECONFIG as CONFIG
		FROM pg_catalog.pg_user
		`

	results, err := db.Query(connName, "", "", query)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusOK, fromQueryResults(results))
}

type usersCreateRequest struct {
	Conn     string `json:"conn"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func handleUsersCreate(w http.ResponseWriter, r *http.Request) {
	var req usersCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	conn, err := db.GetDBFromConf(req.Conn, "")
	if err != nil {
		writeError(w, http.StatusInternalServerError, fmt.Errorf("unable to connect to database: %v", err))
		return
	}
	defer conn.Close()

	query := fmt.Sprintf(`
			create user "%s"
			`, req.Username)
	if req.Password != "" {
		query = fmt.Sprintf(`
			create user "%s" with encrypted password '%s'
			`, req.Username, req.Password)
	}
	if _, err = conn.Exec(query); err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"message": fmt.Sprintf("created user '%s'", req.Username)})
}

func handleUsersDelete(w http.ResponseWriter, r *http.Request) {
	username := r.PathValue("username")
	connName := r.URL.Query().Get("conn")

	if r.URL.Query().Get("confirm") != "true" {
		writeError(w, http.StatusBadRequest, fmt.Errorf("confirm=true is required to drop a user"))
		return
	}

	conn, err := db.GetDBFromConf(connName, "")
	if err != nil {
		writeError(w, http.StatusInternalServerError, fmt.Errorf("unable to connect to database: %v", err))
		return
	}
	defer conn.Close()

	query := fmt.Sprintf(`
			drop user "%s"
			`, username)
	if _, err = conn.Exec(query); err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"message": fmt.Sprintf("deleted user '%s'", username)})
}

type usersAttrRequest struct {
	Conn        string `json:"conn"`
	Password    string `json:"password"`
	Superuser   bool   `json:"superuser"`
	NoSuperuser bool   `json:"nosuperuser"`
	Createdb    bool   `json:"createdb"`
	NoCreatedb  bool   `json:"nocreatedb"`
}

type attrResult struct {
	Attribute string `json:"attribute"`
	Success   bool   `json:"success"`
	Message   string `json:"message"`
}

// handleUsersAttr mirrors cmd/user-attr.go: each attribute is applied
// independently and a failure on one does not prevent the others from running.
func handleUsersAttr(w http.ResponseWriter, r *http.Request) {
	username := r.PathValue("username")

	var req usersAttrRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	conn, err := db.GetDBFromConf(req.Conn, "")
	if err != nil {
		writeError(w, http.StatusInternalServerError, fmt.Errorf("unable to connect to database: %v", err))
		return
	}
	defer conn.Close()

	results := make([]attrResult, 0)

	if req.Password != "" {
		query := fmt.Sprintf(`
		alter user "%s" with encrypted password '%s'
		`, username, req.Password)
		_, err = conn.Exec(query)
		results = append(results, attrResultFor("password", username, "changed password for user '%s'", err))
	}

	if req.Superuser {
		query := fmt.Sprintf(`
		alter user "%s" with superuser
		`, username)
		_, err = conn.Exec(query)
		results = append(results, attrResultFor("superuser", username, "user '%s' is now a superuser", err))
	}

	if req.NoSuperuser {
		query := fmt.Sprintf(`
		alter user "%s" with nosuperuser
		`, username)
		_, err = conn.Exec(query)
		results = append(results, attrResultFor("nosuperuser", username, "user '%s' is now a standard user", err))
	}

	if req.Createdb {
		query := fmt.Sprintf(`
		alter user "%s" with createdb
		`, username)
		_, err = conn.Exec(query)
		results = append(results, attrResultFor("createdb", username, "user '%s' can now create databases", err))
	}

	if req.NoCreatedb {
		query := fmt.Sprintf(`
		alter user "%s" with nocreatedb
		`, username)
		_, err = conn.Exec(query)
		results = append(results, attrResultFor("nocreatedb", username, "user '%s' can not create databases", err))
	}

	writeJSON(w, http.StatusOK, results)
}

func attrResultFor(attribute, username, successFmt string, err error) attrResult {
	if err != nil {
		return attrResult{Attribute: attribute, Success: false, Message: err.Error()}
	}
	return attrResult{Attribute: attribute, Success: true, Message: fmt.Sprintf(successFmt, username)}
}
