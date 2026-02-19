package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/joshburnsxyz/lark/server/progress"
)

// ProgressHandler handles progress-related endpoints.
type ProgressHandler struct {
	Progress progress.Store
}

// Get handles GET /progress.
func (h *ProgressHandler) Get(w http.ResponseWriter, r *http.Request) {
	playerID := r.URL.Query().Get("playerId")
	if playerID == "" {
		playerID = r.Header.Get("X-Player-ID")
	}
	if playerID == "" {
		writeError(w, http.StatusBadRequest, "playerId query parameter or X-Player-ID header required")
		return
	}

	resp, err := h.Progress.Get(playerID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to get progress")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
