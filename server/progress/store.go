package progress

import "github.com/joshburnsxyz/lark/api"

// Store defines the interface for progress persistence.
type Store interface {
	Get(playerID string) (*api.ProgressResponse, error)
	AddCompleted(playerID string, scenario api.CompletedScenario) error
	AddVocab(playerID string, items []api.VocabItem) error
}
