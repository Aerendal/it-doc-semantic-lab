package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

// NewClassifyCmd returns the 'classify' subcommand.
func NewClassifyCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "classify",
		Short: "Detect document class for ingested documents",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Fprintln(cmd.OutOrStdout(), "classify [not yet implemented]")
			return nil
		},
	}
}
