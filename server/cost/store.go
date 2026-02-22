package cost

// Store tracks cumulative AI cost per player.
type Store interface {
	Add(playerID string, amount float64)
	Get(playerID string) float64
}
