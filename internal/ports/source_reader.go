package ports

import (
	"context"
	"io"

	"github.com/it-doc-semantic-lab/itdlab/internal/domain"
)

// SourceReader abstracts reading raw IT documentation sources from any storage backend.
type SourceReader interface {
	// ListSources returns paths to all readable source files under a root.
	ListSources(ctx context.Context, root string) ([]string, error)
	// ReadSource returns the content of a single source file.
	ReadSource(ctx context.Context, path string) (io.ReadCloser, error)
}

// DocumentStore persists and queries documents.
type DocumentStore interface {
	Save(ctx context.Context, doc *domain.Document) error
	FindByID(ctx context.Context, id string) (*domain.Document, error)
	FindByCanonicalID(ctx context.Context, canonicalID string) ([]*domain.Document, error)
	ListAll(ctx context.Context) ([]*domain.Document, error)
}
