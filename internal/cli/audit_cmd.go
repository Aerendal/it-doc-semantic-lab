package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

// NewAuditCmd returns the 'audit' subcommand tree.
func NewAuditCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "audit",
		Short: "Audit run history, evidence packs, and quality gates",
	}
	cmd.AddCommand(newAuditRunsCmd(), newAuditEvidenceCmd())
	return cmd
}

func newAuditRunsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "runs",
		Short: "List all recorded runs",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Fprintln(cmd.OutOrStdout(), "audit runs [not yet implemented]")
			return nil
		},
	}
}

func newAuditEvidenceCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "evidence [run-id]",
		Short: "Show evidence pack for a run",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			runID := ""
			if len(args) > 0 {
				runID = args[0]
			}
			fmt.Fprintf(cmd.OutOrStdout(), "audit evidence: run=%q [not yet implemented]\n", runID)
			return nil
		},
	}
}
