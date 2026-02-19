package session

import (
	"fmt"
	"sync"
	"time"
)

const sessionTTL = 2 * time.Hour

// MemoryStore is an in-memory session store.
type MemoryStore struct {
	mu       sync.RWMutex
	sessions map[string]*Session
}

// NewMemoryStore creates a new in-memory session store and starts a reaper goroutine.
func NewMemoryStore() *MemoryStore {
	s := &MemoryStore{
		sessions: make(map[string]*Session),
	}
	go s.reaper()
	return s
}

func (m *MemoryStore) Get(id string) (*Session, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	s, ok := m.sessions[id]
	if !ok {
		return nil, fmt.Errorf("session not found: %s", id)
	}
	return s, nil
}

func (m *MemoryStore) Save(s *Session) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	s.LastActive = time.Now()
	m.sessions[s.ID] = s
	return nil
}

func (m *MemoryStore) Delete(id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.sessions, id)
	return nil
}

func (m *MemoryStore) reaper() {
	ticker := time.NewTicker(10 * time.Minute)
	defer ticker.Stop()
	for range ticker.C {
		m.mu.Lock()
		now := time.Now()
		for id, s := range m.sessions {
			if now.Sub(s.LastActive) > sessionTTL {
				delete(m.sessions, id)
			}
		}
		m.mu.Unlock()
	}
}
