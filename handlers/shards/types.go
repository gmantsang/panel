package shards

import "time"

// Shard represents a bot shard
type Shard struct {
	ID        int
	Bot       string
	Number    int
	Timestamp time.Time
	UserID    string
}
