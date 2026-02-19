package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/joshburnsxyz/lark/api"
	"github.com/joshburnsxyz/lark/server/ai"
	"github.com/joshburnsxyz/lark/server/session"
)

// ScenarioHandler handles scenario-related endpoints.
type ScenarioHandler struct {
	AI       ai.Client
	Sessions session.Store
}

// List handles GET /scenarios.
func (h *ScenarioHandler) List(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(api.Scenarios)
}

// Languages handles GET /languages.
func (h *ScenarioHandler) Languages(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(api.Languages)
}

// Start handles POST /scenarios/start.
func (h *ScenarioHandler) Start(w http.ResponseWriter, r *http.Request) {
	var req api.StartRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	scenario := resolveScenario(&req)
	if scenario == nil {
		writeError(w, http.StatusBadRequest, "unknown scenario: "+req.ScenarioID)
		return
	}

	lang := api.LanguageByCode(req.Language)
	if lang == nil {
		writeError(w, http.StatusBadRequest, "unsupported language: "+req.Language)
		return
	}

	playerID := r.Header.Get("X-Player-ID")
	if playerID == "" {
		playerID = "anonymous"
	}

	msg, history, err := h.AI.StartScenario(r.Context(), scenario, lang)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to start scenario: "+err.Error())
		return
	}

	sess := &session.Session{
		ID:           generateID(),
		PlayerID:     playerID,
		ScenarioID:   req.ScenarioID,
		CustomPrompt: req.CustomPrompt,
		Language:     req.Language,
		TurnCount:    1,
		History:      history,
		VocabSeen:    msg.Vocabulary,
		LastMessage:  msg,
		CreatedAt:    time.Now(),
		LastActive:   time.Now(),
	}

	if err := h.Sessions.Save(sess); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to save session")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(api.StartResponse{
		SessionID: sess.ID,
		Message:   *msg,
	})
}

// StartStream handles POST /scenarios/start/stream with SSE.
func (h *ScenarioHandler) StartStream(w http.ResponseWriter, r *http.Request) {
	var req api.StartRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	scenario := resolveScenario(&req)
	if scenario == nil {
		writeError(w, http.StatusBadRequest, "unknown scenario: "+req.ScenarioID)
		return
	}

	lang := api.LanguageByCode(req.Language)
	if lang == nil {
		writeError(w, http.StatusBadRequest, "unsupported language: "+req.Language)
		return
	}

	playerID := r.Header.Get("X-Player-ID")
	if playerID == "" {
		playerID = "anonymous"
	}

	flusher, ok := w.(http.Flusher)
	if !ok {
		writeError(w, http.StatusInternalServerError, "streaming not supported")
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	callback := func(token string) {
		data, _ := json.Marshal(map[string]string{"token": token})
		fmt.Fprintf(w, "data: %s\n\n", data)
		flusher.Flush()
	}

	msg, history, err := h.AI.StartScenarioStream(r.Context(), scenario, lang, callback)
	if err != nil {
		errData, _ := json.Marshal(map[string]string{"error": err.Error()})
		fmt.Fprintf(w, "data: %s\n\n", errData)
		flusher.Flush()
		return
	}

	sess := &session.Session{
		ID:           generateID(),
		PlayerID:     playerID,
		ScenarioID:   req.ScenarioID,
		CustomPrompt: req.CustomPrompt,
		Language:     req.Language,
		TurnCount:    1,
		History:      history,
		VocabSeen:    msg.Vocabulary,
		LastMessage:  msg,
		CreatedAt:    time.Now(),
		LastActive:   time.Now(),
	}

	if err := h.Sessions.Save(sess); err != nil {
		errData, _ := json.Marshal(map[string]string{"error": "failed to save session"})
		fmt.Fprintf(w, "data: %s\n\n", errData)
		flusher.Flush()
		return
	}

	final := api.StartResponse{
		SessionID: sess.ID,
		Message:   *msg,
	}
	doneData, _ := json.Marshal(map[string]any{"done": true, "sessionId": sess.ID, "message": final.Message})
	fmt.Fprintf(w, "data: %s\n\n", doneData)
	flusher.Flush()
}
