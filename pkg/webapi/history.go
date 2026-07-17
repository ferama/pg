package webapi

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/ferama/pg/pkg/history"
)

type historyItem struct {
	Idx   int    `json:"idx"`
	Query string `json:"query"`
}

func handleHistoryList(w http.ResponseWriter, r *http.Request) {
	h := history.GetInstance()
	items := h.GetList()

	out := make([]historyItem, 0, len(items))
	for idx := len(items) - 1; idx >= 0; idx-- {
		out = append(out, historyItem{Idx: idx, Query: items[idx]})
	}
	writeJSON(w, http.StatusOK, out)
}

func handleHistoryDelete(w http.ResponseWriter, r *http.Request) {
	idx, err := strconv.Atoi(r.PathValue("idx"))
	if err != nil {
		writeError(w, http.StatusBadRequest, fmt.Errorf("invalid idx: %v", err))
		return
	}
	history.GetInstance().DeleteAtIdx(idx)
	w.WriteHeader(http.StatusNoContent)
}
