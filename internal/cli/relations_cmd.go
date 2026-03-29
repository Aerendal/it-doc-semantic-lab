package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

// NewRelationsCmd returns the 'relations' subcommand tree.
func NewRelationsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "relations",
		Short: "Infer and inspect cross-document semantic relations",
	}
	cmd.AddCommand(newRelationsShowCmd(), newRelationsExplainCmd())
	return cmd
}

func newRelationsShowCmd() *cobra.Command {
	var docID string
	cmd := &cobra.Command{
		Use:   "show",
		Short: "Show relations for a document",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Fprintf(cmd.OutOrStdout(), "relations show: doc=%q [not yet implemented]\n", docID)
			return nil
		},
	}
	cmd.Flags().StringVar(&docID, "doc", "", "document ID to show relations for")
	return cmd
}

func newRelationsExplainCmd() *cobra.Command {
	var relID string
	cmd := &cobra.Command{
		Use:   "explain",
		Short: "Explain why a relation was inferred",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Fprintf(cmd.OutOrStdout(), "relations explain: rel=%q [not yet implemented]\n", relID)
			return nil
		},
	}
	cmd.Flags().StringVar(&relID, "rel", "", "relation ID to explain")
	return cmd
}
