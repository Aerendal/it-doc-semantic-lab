package ports

import "context"

// ReportWriter writes audit and analysis reports to a backend (filesystem, stdout, etc.).
type ReportWriter interface {
	// WriteReport persists a named report for a specific run.
	WriteReport(ctx context.Context, runID, reportName, content string) error
	// ReadReport retrieves a previously written report.
	ReadReport(ctx context.Context, runID, reportName string) (string, error)
}
