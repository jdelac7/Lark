package session

// Store defines the interface for session persistence.
type Store interface {
	Get(id string) (*Session, error)
	Save(s *Session) error
	Delete(id string) error
}
