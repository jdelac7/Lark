package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/joshburnsxyz/lark/api"
	"github.com/joshburnsxyz/lark/server/ai"
	"github.com/joshburnsxyz/lark/server/progress"
	"github.com/joshburnsxyz/lark/server/session"
)

// GameHandler handles game interaction endpoints.
type GameHandler struct {
	AI       ai.Client
	Sessions session.Store
	Progress progress.Store
}

// Input handles POST /game/input.
func (h *GameHandler) Input(w http.ResponseWriter, r *http.Request) {
	var req api.PlayerInputRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	sess, err := h.Sessions.Get(req.SessionID)
	if err != nil {
		writeError(w, http.StatusNotFound, "session not found")
		return
	}

	scenario := resolveSessionScenario(sess)
	lang := api.LanguageByCode(sess.Language)
	if scenario == nil || lang == nil {
		writeError(w, http.StatusInternalServerError, "invalid session state")
		return
	}

	// Format input based on mode
	var inputText string
	switch req.Mode {
	case api.InputModeChoice:
		if req.ChoiceIndex < 0 || sess.LastMessage == nil || req.ChoiceIndex >= len(sess.LastMessage.Choices) {
			writeError(w, http.StatusBadRequest, "invalid choice index")
			return
		}
		choiceText := sess.LastMessage.Choices[req.ChoiceIndex].Text
		inputText = ai.FormatChoiceInput(req.ChoiceIndex, choiceText)
	case api.InputModeFreeText:
		if req.Text == "" {
			writeError(w, http.StatusBadRequest, "text is required for free_text mode")
			return
		}
		inputText = ai.FormatFreeTextInput(req.Text)
	default:
		writeError(w, http.StatusBadRequest, "mode must be 'choice' or 'free_text'")
		return
	}

	msg, correction, newHistory, err := h.AI.SendInput(r.Context(), scenario, lang, sess.History, inputText)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to process input: "+err.Error())
		return
	}

	// Update session
	sess.TurnCount++
	sess.History = newHistory
	sess.LastMessage = msg
	sess.VocabSeen = append(sess.VocabSeen, msg.Vocabulary...)

	if err := h.Sessions.Save(sess); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to save session")
		return
	}

	// If scenario finished, record progress
	if msg.Finished {
		playerID := r.Header.Get("X-Player-ID")
		if playerID == "" {
			playerID = sess.PlayerID
		}
		h.Progress.AddCompleted(playerID, api.CompletedScenario{
			ScenarioID: sess.ScenarioID,
			Language:   sess.Language,
			TurnCount:  sess.TurnCount,
		})
		h.Progress.AddVocab(playerID, sess.VocabSeen)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(api.PlayerInputResponse{
		Message:    *msg,
		Correction: correction,
	})
}

// InputStream handles POST /game/input/stream with SSE.
func (h *GameHandler) InputStream(w http.ResponseWriter, r *http.Request) {
	var req api.PlayerInputRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	sess, err := h.Sessions.Get(req.SessionID)
	if err != nil {
		writeError(w, http.StatusNotFound, "session not found")
		return
	}

	scenario := resolveSessionScenario(sess)
	lang := api.LanguageByCode(sess.Language)
	if scenario == nil || lang == nil {
		writeError(w, http.StatusInternalServerError, "invalid session state")
		return
	}

	var inputText string
	switch req.Mode {
	case api.InputModeChoice:
		if req.ChoiceIndex < 0 || sess.LastMessage == nil || req.ChoiceIndex >= len(sess.LastMessage.Choices) {
			writeError(w, http.StatusBadRequest, "invalid choice index")
			return
		}
		choiceText := sess.LastMessage.Choices[req.ChoiceIndex].Text
		inputText = ai.FormatChoiceInput(req.ChoiceIndex, choiceText)
	case api.InputModeFreeText:
		if req.Text == "" {
			writeError(w, http.StatusBadRequest, "text is required for free_text mode")
			return
		}
		inputText = ai.FormatFreeTextInput(req.Text)
	default:
		writeError(w, http.StatusBadRequest, "mode must be 'choice' or 'free_text'")
		return
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

	msg, correction, newHistory, err := h.AI.SendInputStream(r.Context(), scenario, lang, sess.History, inputText, callback)
	if err != nil {
		errData, _ := json.Marshal(map[string]string{"error": err.Error()})
		fmt.Fprintf(w, "data: %s\n\n", errData)
		flusher.Flush()
		return
	}

	sess.TurnCount++
	sess.History = newHistory
	sess.LastMessage = msg
	sess.VocabSeen = append(sess.VocabSeen, msg.Vocabulary...)

	if err := h.Sessions.Save(sess); err != nil {
		errData, _ := json.Marshal(map[string]string{"error": "failed to save session"})
		fmt.Fprintf(w, "data: %s\n\n", errData)
		flusher.Flush()
		return
	}

	if msg.Finished {
		playerID := r.Header.Get("X-Player-ID")
		if playerID == "" {
			playerID = sess.PlayerID
		}
		h.Progress.AddCompleted(playerID, api.CompletedScenario{
			ScenarioID: sess.ScenarioID,
			Language:   sess.Language,
			TurnCount:  sess.TurnCount,
		})
		h.Progress.AddVocab(playerID, sess.VocabSeen)
	}

	donePayload := map[string]any{
		"done":    true,
		"message": *msg,
	}
	if correction != nil {
		donePayload["correction"] = correction
	}
	doneData, _ := json.Marshal(donePayload)
	fmt.Fprintf(w, "data: %s\n\n", doneData)
	flusher.Flush()
}

// State handles GET /game/state.
func (h *GameHandler) State(w http.ResponseWriter, r *http.Request) {
	sessionID := r.URL.Query().Get("sessionId")
	if sessionID == "" {
		writeError(w, http.StatusBadRequest, "sessionId query parameter required")
		return
	}

	sess, err := h.Sessions.Get(sessionID)
	if err != nil {
		writeError(w, http.StatusNotFound, "session not found")
		return
	}

	resp := api.GameStateResponse{
		SessionID:  sess.ID,
		ScenarioID: sess.ScenarioID,
		Language:   sess.Language,
		TurnCount:  sess.TurnCount,
	}
	if sess.LastMessage != nil {
		resp.Message = *sess.LastMessage
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
