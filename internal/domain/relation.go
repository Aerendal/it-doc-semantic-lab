package domain

import "time"

// RelationType classifies the semantic relationship between two documents.
type RelationType string

const (
	RelationDependsOn   RelationType = "depends_on"
	RelationImplements  RelationType = "implements"
	RelationReferences  RelationType = "references"
	RelationSupersedes  RelationType = "supersedes"
	RelationComplementsRelationType = "complements"
	RelationDerivedFrom RelationType = "derived_from"
)

// RelationSource describes how a relation was discovered.
type RelationSource string

const (
	SourceRuleEngine  RelationSource = "rule_engine"
	SourceExplicit    RelationSource = "explicit"
	SourceInferred    RelationSource = "inferred"
)

// Relation represents a directed semantic relationship between two documents.
type Relation struct {
	ID           string         `json:"id"`
	FromID       string         `json:"from_id"`
	ToID         string         `json:"to_id"`
	Type         RelationType   `json:"type"`
	Source       RelationSource `json:"source"`
	Confidence   float64        `json:"confidence"`
	Explanation  string         `json:"explanation"`
	RuleID       string         `json:"rule_id,omitempty"`
	DiscoveredAt time.Time      `json:"discovered_at"`
}

// RelationRule defines a rule used to infer relations between documents.
type RelationRule struct {
	ID          string       `json:"id"`
	Name        string       `json:"name"`
	FromClass   DocumentClass `json:"from_class"`
	ToClass     DocumentClass `json:"to_class"`
	RelType     RelationType  `json:"rel_type"`
	Description string        `json:"description"`
}
