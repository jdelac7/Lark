package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"net/http"

	"github.com/joshburnsxyz/lark/api"
	"github.com/joshburnsxyz/lark/server/session"
)

func writeError(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(api.ErrorResponse{Error: msg})
}

func generateID() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}

// resolveScenario returns a Scenario for the request. If CustomPrompt is set
// and the ID is not in the catalog, it builds an ad-hoc scenario from the prompt.
func resolveScenario(req *api.StartRequest) *api.Scenario {
	if s := api.ScenarioByID(req.ScenarioID); s != nil {
		return s
	}
	if req.CustomPrompt != "" {
		return &api.Scenario{
			ID:          req.ScenarioID,
			Name:        req.CustomPrompt,
			Description: req.CustomPrompt,
			Difficulty:  api.DifficultyBeginner,
		}
	}
	return nil
}

// resolveSessionScenario reconstructs a Scenario from a session. For custom
// scenarios it rebuilds from the stored prompt; for catalog scenarios it
// looks up by ID.
func resolveSessionScenario(sess *session.Session) *api.Scenario {
	if s := api.ScenarioByID(sess.ScenarioID); s != nil {
		return s
	}
	if sess.CustomPrompt != "" {
		return &api.Scenario{
			ID:          sess.ScenarioID,
			Name:        sess.CustomPrompt,
			Description: sess.CustomPrompt,
			Difficulty:  api.DifficultyBeginner,
		}
	}
	return nil
}
