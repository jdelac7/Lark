package ai

import (
	"context"
	"fmt"

	"github.com/joshburnsxyz/lark/api"
)

// StreamCallback is called with each token as it arrives from the AI backend.
type StreamCallback func(token string)

// Client is the interface both Google and OpenRouter backends implement.
type Client interface {
	// StartScenario begins a new scenario and returns the first game message + opaque history.
	StartScenario(ctx context.Context, scenario *api.Scenario, lang *api.Language) (*api.GameMessage, any, error)
	// SendInput sends player input and returns the response + optional correction + updated history.
	SendInput(ctx context.Context, scenario *api.Scenario, lang *api.Language, history any, input string) (*api.GameMessage, *api.Correction, any, error)

	// StartScenarioStream is like StartScenario but streams tokens via callback.
	StartScenarioStream(ctx context.Context, scenario *api.Scenario, lang *api.Language, callback StreamCallback) (*api.GameMessage, any, error)
	// SendInputStream is like SendInput but streams tokens via callback.
	SendInputStream(ctx context.Context, scenario *api.Scenario, lang *api.Language, history any, input string, callback StreamCallback) (*api.GameMessage, *api.Correction, any, error)

	// BuildStartHistory reconstructs backend-specific conversation history from a
	// cached first-turn response, so subsequent SendInput calls work correctly.
	BuildStartHistory(scenario *api.Scenario, lang *api.Language, responseText string) any
}

// FormatChoiceInput formats a choice selection for the LLM.
func FormatChoiceInput(choiceIndex int, choiceText string) string {
	return fmt.Sprintf("[PLAYER CHOSE OPTION %d]: %q", choiceIndex+1, choiceText)
}

// FormatFreeTextInput formats free text input for the LLM.
func FormatFreeTextInput(text string) string {
	return fmt.Sprintf("[PLAYER FREE TEXT - evaluate grammar]: %q", text)
}
