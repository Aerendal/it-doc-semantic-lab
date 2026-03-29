package ports

import (
	"context"

	"github.com/it-doc-semantic-lab/itdlab/internal/domain"
)

// RelationStore persists and queries document relations.
type RelationStore interface {
	Save(ctx context.Context, rel *domain.Relation) error
	FindByFromID(ctx context.Context, fromID string) ([]*domain.Relation, error)
	FindByToID(ctx context.Context, toID string) ([]*domain.Relation, error)
	FindByType(ctx context.Context, relType domain.RelationType) ([]*domain.Relation, error)
	ListAll(ctx context.Context) ([]*domain.Relation, error)
}

// RelationRuleStore stores and retrieves relation inference rules.
type RelationRuleStore interface {
	SaveRule(ctx context.Context, rule *domain.RelationRule) error
	ListRules(ctx context.Context) ([]*domain.RelationRule, error)
}
