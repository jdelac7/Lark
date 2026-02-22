package cost

import (
	"database/sql"
	"log"

	_ "modernc.org/sqlite"
)

// SQLiteStore persists per-player AI cost in a SQLite database.
type SQLiteStore struct {
	db *sql.DB
}

// NewSQLiteStore opens (or creates) the SQLite database at path and
// initialises the costs table.
func NewSQLiteStore(path string) (*SQLiteStore, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}

	// WAL mode for better concurrent read performance.
	if _, err := db.Exec("PRAGMA journal_mode=WAL"); err != nil {
		db.Close()
		return nil, err
	}

	const create = `CREATE TABLE IF NOT EXISTS player_costs (
		player_id TEXT PRIMARY KEY,
		total_cost REAL NOT NULL DEFAULT 0
	)`
	if _, err := db.Exec(create); err != nil {
		db.Close()
		return nil, err
	}

	return &SQLiteStore{db: db}, nil
}

func (s *SQLiteStore) Add(playerID string, amount float64) {
	_, err := s.db.Exec(
		`INSERT INTO player_costs (player_id, total_cost) VALUES (?, ?)
		 ON CONFLICT(player_id) DO UPDATE SET total_cost = total_cost + excluded.total_cost`,
		playerID, amount,
	)
	if err != nil {
		log.Printf("[cost] failed to add cost for %s: %v", playerID, err)
	}
}

func (s *SQLiteStore) Get(playerID string) float64 {
	var total float64
	err := s.db.QueryRow("SELECT total_cost FROM player_costs WHERE player_id = ?", playerID).Scan(&total)
	if err != nil {
		return 0 // not found or error → treat as zero
	}
	return total
}

// DB returns the underlying database connection for use by EventRecorder.
func (s *SQLiteStore) DB() *sql.DB {
	return s.db
}

// Close closes the underlying database connection.
func (s *SQLiteStore) Close() error {
	return s.db.Close()
}
