package progress

import (
	"sync"

	"github.com/joshburnsxyz/lark/api"
)

type playerData struct {
	Completed []api.CompletedScenario
	Vocab     []api.VocabItem
}

// MemoryStore is an in-memory progress store.
type MemoryStore struct {
	mu      sync.RWMutex
	players map[string]*playerData
}

// NewMemoryStore creates a new in-memory progress store.
func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		players: make(map[string]*playerData),
	}
}

func (m *MemoryStore) Get(playerID string) (*api.ProgressResponse, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	p := m.getOrCreate(playerID)
	return &api.ProgressResponse{
		PlayerID:           playerID,
		CompletedScenarios: p.Completed,
		VocabBank:          p.Vocab,
	}, nil
}

func (m *MemoryStore) AddCompleted(playerID string, scenario api.CompletedScenario) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	p := m.getOrCreate(playerID)
	p.Completed = append(p.Completed, scenario)
	return nil
}

func (m *MemoryStore) AddVocab(playerID string, items []api.VocabItem) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	p := m.getOrCreate(playerID)
	// Deduplicate by word
	existing := make(map[string]bool)
	for _, v := range p.Vocab {
		existing[v.Word] = true
	}
	for _, v := range items {
		if !existing[v.Word] {
			p.Vocab = append(p.Vocab, v)
			existing[v.Word] = true
		}
	}
	return nil
}

func (m *MemoryStore) getOrCreate(playerID string) *playerData {
	p, ok := m.players[playerID]
	if !ok {
		p = &playerData{
			Completed: []api.CompletedScenario{},
			Vocab:     []api.VocabItem{},
		}
		m.players[playerID] = p
	}
	return p
}
