package domain

// SectionRole identifies the functional role a section plays within a document.
type SectionRole string

const (
	RolePurpose       SectionRole = "purpose"
	RoleScope         SectionRole = "scope"
	RoleDefinitions   SectionRole = "definitions"
	RoleResponsibility SectionRole = "responsibility"
	RoleProcedure     SectionRole = "procedure"
	RoleReference     SectionRole = "reference"
	RoleAppendix      SectionRole = "appendix"
	RoleUnknown       SectionRole = "unknown"
)

// SectionArchetype defines the canonical pattern for a section role.
type SectionArchetype struct {
	Role        SectionRole `json:"role"`
	Keywords    []string    `json:"keywords"`
	Required    bool        `json:"required"`
	Description string      `json:"description"`
}

// Section represents a parsed section within a document.
type Section struct {
	ID         string      `json:"id"`
	DocumentID string      `json:"document_id"`
	Heading    string      `json:"heading"`
	Level      int         `json:"level"`
	Role       SectionRole `json:"role"`
	Confidence float64     `json:"confidence"`
	Position   int         `json:"position"`
	WordCount  int         `json:"word_count"`
}

// SectionAnomaly records a detected problem with a section.
type SectionAnomaly struct {
	SectionID   string `json:"section_id"`
	DocumentID  string `json:"document_id"`
	AnomalyType string `json:"anomaly_type"`
	Description string `json:"description"`
}
