package jsonl

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"github.com/it-doc-semantic-lab/itdlab/internal/ports"
)

// EventLog is an append-only JSONL-backed implementation of ports.EventLog.
type EventLog struct {
	path string
	mu   sync.Mutex
}

// New creates an EventLog that writes to path. The file is created if it does not exist.
func New(path string) (*EventLog, error) {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return nil, fmt.Errorf("jsonl event log: open %s: %w", path, err)
	}
	f.Close()
	return &EventLog{path: path}, nil
}

// Append writes a single event as a JSON line. Thread-safe.
func (l *EventLog) Append(_ context.Context, event ports.RunEvent) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	f, err := os.OpenFile(l.path, os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return fmt.Errorf("jsonl event log: open for append: %w", err)
	}
	defer f.Close()

	b, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("jsonl event log: marshal: %w", err)
	}

	_, err = fmt.Fprintf(f, "%s\n", b)
	return err
}

// ReadAll returns every event in the log.
func (l *EventLog) ReadAll(_ context.Context) ([]ports.RunEvent, error) {
	return l.scan(func(_ ports.RunEvent) bool { return true })
}

// ReadRun returns events filtered to a specific run_id.
func (l *EventLog) ReadRun(_ context.Context, runID string) ([]ports.RunEvent, error) {
	return l.scan(func(e ports.RunEvent) bool { return e.RunID == runID })
}

func (l *EventLog) scan(keep func(ports.RunEvent) bool) ([]ports.RunEvent, error) {
	l.mu.Lock()
	defer l.mu.Unlock()

	f, err := os.Open(l.path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("jsonl event log: read: %w", err)
	}
	defer f.Close()

	var events []ports.RunEvent
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}
		var e ports.RunEvent
		if err := json.Unmarshal(line, &e); err != nil {
			return nil, fmt.Errorf("jsonl event log: parse line: %w", err)
		}
		if keep(e) {
			events = append(events, e)
		}
	}
	return events, scanner.Err()
}
