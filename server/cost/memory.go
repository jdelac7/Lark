package cost

import "sync"

// MemoryStore is an in-memory cost store.
type MemoryStore struct {
	mu    sync.RWMutex
	costs map[string]float64
}

// NewMemoryStore creates a new in-memory cost store.
func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		costs: make(map[string]float64),
	}
}

func (m *MemoryStore) Add(playerID string, amount float64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.costs[playerID] += amount
}

func (m *MemoryStore) Get(playerID string) float64 {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.costs[playerID]
}
