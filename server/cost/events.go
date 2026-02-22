package cost

import (
	"database/sql"
	"log"
)

// EventRecorder writes individual cost events for analytics.
type EventRecorder struct {
	db *sql.DB
}

// NewEventRecorder creates the cost_events table and returns a recorder.
func NewEventRecorder(db *sql.DB) (*EventRecorder, error) {
	const create = `
	CREATE TABLE IF NOT EXISTS cost_events (
		id          INTEGER PRIMARY KEY AUTOINCREMENT,
		player_id   TEXT NOT NULL,
		scenario_id TEXT NOT NULL,
		language    TEXT NOT NULL DEFAULT '',
		is_custom   INTEGER NOT NULL DEFAULT 0,
		amount      REAL NOT NULL,
		created_at  TEXT NOT NULL DEFAULT (datetime('now'))
	);
	CREATE INDEX IF NOT EXISTS idx_cost_events_player ON cost_events(player_id);
	CREATE INDEX IF NOT EXISTS idx_cost_events_created ON cost_events(created_at);
	`
	if _, err := db.Exec(create); err != nil {
		return nil, err
	}
	// Migration: add language column to existing tables
	db.Exec("ALTER TABLE cost_events ADD COLUMN language TEXT NOT NULL DEFAULT ''")
	return &EventRecorder{db: db}, nil
}

// Record writes a single cost event.
func (r *EventRecorder) Record(playerID, scenarioID, language string, isCustom bool, amount float64) {
	custom := 0
	if isCustom {
		custom = 1
	}
	_, err := r.db.Exec(
		`INSERT INTO cost_events (player_id, scenario_id, language, is_custom, amount) VALUES (?, ?, ?, ?, ?)`,
		playerID, scenarioID, language, custom, amount,
	)
	if err != nil {
		log.Printf("[cost] failed to record event for %s: %v", playerID, err)
	}
}
