package domain

import "time"

// AuthorityRef links a document or section to an external regulatory or standards authority.
type AuthorityRef struct {
	ID         string    `json:"id"`
	DocumentID string    `json:"document_id"`
	SectionID  string    `json:"section_id,omitempty"`
	Authority  string    `json:"authority"`
	Clause     string    `json:"clause"`
	URL        string    `json:"url,omitempty"`
	LinkedAt   time.Time `json:"linked_at"`
}

// AuthorityCoverage summarizes how well a document set covers a given authority.
type AuthorityCoverage struct {
	Authority     string  `json:"authority"`
	TotalClauses  int     `json:"total_clauses"`
	LinkedClauses int     `json:"linked_clauses"`
	CoveragePct   float64 `json:"coverage_pct"`
}
