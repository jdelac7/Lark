package ai

import (
	"context"
	"fmt"

	"github.com/joshburnsxyz/lark/api"
)

// StreamCallback is called with each token as it arrives from the AI backend.
type StreamCallback func(token string)

// Client is the interface both Google and OpenRouter backends implement.
// Methods that call the AI return a float64 cost (USD) alongside the response.
type Client interface {
	// StartScenario begins a new scenario and returns the first game message + opaque history + cost.
	StartScenario(ctx context.Context, scenario *api.Scenario, lang *api.Language, explanationLang string) (*api.GameMessage, any, float64, error)
	// SendInput sends player input and returns the response + optional correction + updated history + cost.
	SendInput(ctx context.Context, scenario *api.Scenario, lang *api.Language, explanationLang string, history any, input string) (*api.GameMessage, *api.Correction, any, float64, error)

	// StartScenarioStream is like StartScenario but streams tokens via callback.
	StartScenarioStream(ctx context.Context, scenario *api.Scenario, lang *api.Language, explanationLang string, callback StreamCallback) (*api.GameMessage, any, float64, error)
	// SendInputStream is like SendInput but streams tokens via callback.
	SendInputStream(ctx context.Context, scenario *api.Scenario, lang *api.Language, explanationLang string, history any, input string, callback StreamCallback) (*api.GameMessage, *api.Correction, any, float64, error)

	// BuildStartHistory reconstructs backend-specific conversation history from a
	// cached first-turn response, so subsequent SendInput calls work correctly.
	BuildStartHistory(scenario *api.Scenario, lang *api.Language, explanationLang string, responseText string) any
}

// FormatChoiceInput formats a choice selection for the LLM.
func FormatChoiceInput(choiceIndex int, choiceText string) string {
	return fmt.Sprintf("[PLAYER CHOSE OPTION %d]: %q", choiceIndex+1, choiceText)
}

// FormatFreeTextInput formats free text input for the LLM.
func FormatFreeTextInput(text string) string {
	return fmt.Sprintf("[PLAYER FREE TEXT]: %q\nContinue the story with the full JSON response (narrative, translation, npcDialog, choices, vocabulary, correction, finished). Include a \"correction\" object evaluating the player's grammar.", text)
}
