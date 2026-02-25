package ai

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/joshburnsxyz/lark/api"
)

// TurnResponse is the raw JSON structure from the LLM.
type TurnResponse struct {
	Narrative            string          `json:"narrative"`
	Translation          string          `json:"translation"`
	NPCDialog            string          `json:"npcDialog"`
	NPCDialogTranslation string          `json:"npcDialogTranslation"`
	Vocabulary           []api.VocabItem `json:"vocabulary"`
	Choices              []api.Choice    `json:"choices"`
	Correction           *api.Correction `json:"correction"`
	Finished             bool            `json:"finished"`
}

// flexTurnResponse uses json.RawMessage for fields where the LLM sometimes
// returns strings instead of objects (e.g. vocabulary as ["word"] instead of
// [{"word":"word","translation":"..."}]).
type flexTurnResponse struct {
	Narrative            string            `json:"narrative"`
	Translation          string            `json:"translation"`
	NPCDialog            string            `json:"npcDialog"`
	NPCDialogTranslation string            `json:"npcDialogTranslation"`
	Vocabulary           json.RawMessage   `json:"vocabulary"`
	Choices              json.RawMessage   `json:"choices"`
	Correction           *api.Correction   `json:"correction"`
	Finished             bool              `json:"finished"`
}

// TrimHistoryJSON strips translations, vocabulary, and corrections from an
// assistant response JSON so that subsequent turns send fewer tokens.
// Keeps the same JSON structure so the LLM stays oriented.
func TrimHistoryJSON(text string) string {
	var tr TurnResponse
	if err := json.Unmarshal([]byte(text), &tr); err != nil {
		// Try flex parse in case of type mismatches
		msg, _, flexErr := flexParseTurn(text)
		if flexErr != nil {
			return text
		}
		tr = TurnResponse{
			Narrative:            msg.Narrative,
			Translation:          msg.Translation,
			NPCDialog:            msg.NPCDialog,
			NPCDialogTranslation: msg.NPCDialogTranslation,
			Choices:              msg.Choices,
			Finished:             msg.Finished,
		}
	}
	tr.Vocabulary = nil
	tr.Correction = nil
	out, err := json.Marshal(tr)
	if err != nil {
		return text
	}
	return string(out)
}

// ParseTurnJSON parses raw JSON text into a GameMessage and optional Correction.
func ParseTurnJSON(text string) (*api.GameMessage, *api.Correction, error) {
	if text == "" {
		return nil, nil, fmt.Errorf("empty response from LLM")
	}

	// Try strict parse first
	var tr TurnResponse
	if err := json.Unmarshal([]byte(text), &tr); err != nil {
		// Fall back to flexible parsing for type mismatches
		return flexParseTurn(text)
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

// flexParseTurn handles LLM responses where choices or vocabulary are returned
// as strings instead of structured objects.
func flexParseTurn(text string) (*api.GameMessage, *api.Correction, error) {
	var flex flexTurnResponse
	if err := json.Unmarshal([]byte(text), &flex); err != nil {
		return nil, nil, fmt.Errorf("unmarshaling response: %w (raw: %.500s)", err, text)
	}

	choices := parseFlexChoices(flex.Choices)
	vocab := parseFlexVocabulary(flex.Vocabulary)

	msg := &api.GameMessage{
		Narrative:            flex.Narrative,
		Translation:          flex.Translation,
		NPCDialog:            flex.NPCDialog,
		NPCDialogTranslation: flex.NPCDialogTranslation,
		Choices:              choices,
		Vocabulary:           vocab,
		Finished:             flex.Finished,
	}

	return msg, flex.Correction, nil
}

// parseFlexChoices handles choices as either []Choice or []string.
func parseFlexChoices(raw json.RawMessage) []api.Choice {
	if len(raw) == 0 {
		return nil
	}

	// Try structured first
	var choices []api.Choice
	if err := json.Unmarshal(raw, &choices); err == nil {
		return choices
	}

	// Fall back to string array
	var strs []string
	if err := json.Unmarshal(raw, &strs); err == nil {
		choices = make([]api.Choice, len(strs))
		for i, s := range strs {
			choices[i] = api.Choice{Text: s}
		}
		return choices
	}

	return nil
}

// parseFlexVocabulary handles vocabulary as either []VocabItem or []string.
func parseFlexVocabulary(raw json.RawMessage) []api.VocabItem {
	if len(raw) == 0 {
		return nil
	}

	// Try structured first
	var vocab []api.VocabItem
	if err := json.Unmarshal(raw, &vocab); err == nil {
		return vocab
	}

	// Fall back to string array
	var strs []string
	if err := json.Unmarshal(raw, &strs); err == nil {
		vocab = make([]api.VocabItem, len(strs))
		for i, s := range strs {
			vocab[i] = api.VocabItem{Word: s}
		}
		return vocab
	}

	// Fall back to mixed array (some objects, some strings)
	var mixed []json.RawMessage
	if err := json.Unmarshal(raw, &mixed); err == nil {
		for _, item := range mixed {
			trimmed := strings.TrimSpace(string(item))
			if trimmed[0] == '"' {
				var s string
				if json.Unmarshal(item, &s) == nil {
					vocab = append(vocab, api.VocabItem{Word: s})
				}
			} else {
				var v api.VocabItem
				if json.Unmarshal(item, &v) == nil {
					vocab = append(vocab, v)
				}
			}
		}
		return vocab
	}

	return nil
}
