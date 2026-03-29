package domain

import "time"

// DocumentClass identifies the semantic category of a document.
type DocumentClass string

// DocumentStatus tracks lifecycle state within the lab.
type DocumentStatus string

const (
	StatusRaw        DocumentStatus = "raw"
	StatusIngested   DocumentStatus = "ingested"
	StatusNormalized DocumentStatus = "normalized"
	StatusClassified DocumentStatus = "classified"
	StatusExported   DocumentStatus = "exported"
)

// Document is the core entity representing a single IT documentation artifact.
type Document struct {
	ID          string         `json:"id"`
	CanonicalID string         `json:"canonical_id"`
	RawName     string         `json:"raw_name"`
	SourcePath  string         `json:"source_path"`
	Class       DocumentClass  `json:"class"`
	Industry    Industry       `json:"industry"`
	Phase       Phase          `json:"phase"`
	Status      DocumentStatus `json:"status"`
	Checksum    string         `json:"checksum"`
	IngestedAt  time.Time      `json:"ingested_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}

// DocumentFamily groups documents that represent the same concept across industries or phases.
type DocumentFamily struct {
	ID          string   `json:"id"`
	CanonicalID string   `json:"canonical_id"`
	MemberIDs   []string `json:"member_ids"`
}
