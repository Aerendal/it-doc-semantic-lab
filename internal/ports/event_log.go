package ports

import (
	"context"
	"time"
)

// RunEvent is a single entry in the append-only run event log.
type RunEvent struct {
	Timestamp  time.Time         `json:"ts"`
	RunID      string            `json:"run_id"`
	Step       string            `json:"step"`
	Entity     string            `json:"entity"`
	EntityID   string            `json:"entity_id"`
	Action     string            `json:"action"`
	Before     string            `json:"before,omitempty"`
	After      string            `json:"after,omitempty"`
	Meta       map[string]string `json:"meta,omitempty"`
}

// EventLog is an append-only log of all run events for audit and reproducibility.
type EventLog interface {
	// Append writes a new event to the log. Must be safe for concurrent writes.
	Append(ctx context.Context, event RunEvent) error
	// ReadAll returns all events. Intended for audit and test use only.
	ReadAll(ctx context.Context) ([]RunEvent, error)
	// ReadRun returns events filtered to a specific run_id.
	ReadRun(ctx context.Context, runID string) ([]RunEvent, error)
}
