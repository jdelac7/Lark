package session

import (
	"time"

	"github.com/joshburnsxyz/lark/api"
)

// Session represents an active game session.
type Session struct {
	ID           string
	PlayerID     string
	ScenarioID   string
	CustomPrompt string // non-empty for user-created scenarios
	Language     string
	TurnCount    int
	History      any // opaque history managed by the AI backend
	VocabSeen    []api.VocabItem
	LastMessage  *api.GameMessage
	CreatedAt    time.Time
	LastActive   time.Time
}
