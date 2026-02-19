package ai

import (
	"encoding/json"
	"fmt"

	"github.com/joshburnsxyz/lark/api"
)

// TurnResponse is the raw JSON structure from the LLM.
type TurnResponse struct {
	Narrative            string          `json:"narrative"`
	Translation          string          `json:"translation"`
	NPCDialog            string          `json:"npcDialog"`
	NPCDialogTranslation string          `json:"npcDialogTranslation"`
	Choices              []api.Choice    `json:"choices"`
	Vocabulary           []api.VocabItem `json:"vocabulary"`
	Correction           *api.Correction `json:"correction"`
	Finished             bool            `json:"finished"`
}

// ParseTurnJSON parses raw JSON text into a GameMessage and optional Correction.
func ParseTurnJSON(text string) (*api.GameMessage, *api.Correction, error) {
	if text == "" {
		return nil, nil, fmt.Errorf("empty response from LLM")
	}

	var tr TurnResponse
	if err := json.Unmarshal([]byte(text), &tr); err != nil {
		return nil, nil, fmt.Errorf("unmarshaling response: %w (raw: %.500s)", err, text)
	}

	msg := &api.GameMessage{
		Narrative:            tr.Narrative,
		Translation:          tr.Translation,
		NPCDialog:            tr.NPCDialog,
		NPCDialogTranslation: tr.NPCDialogTranslation,
		Choices:              tr.Choices,
		Vocabulary:           tr.Vocabulary,
		Finished:             tr.Finished,
	}

	return msg, tr.Correction, nil
}
