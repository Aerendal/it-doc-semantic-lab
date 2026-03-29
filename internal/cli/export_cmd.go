package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

// NewExportCmd returns the 'export' subcommand.
func NewExportCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "export",
		Short: "Export stable metadata to target repository",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Fprintln(cmd.OutOrStdout(), "export [not yet implemented]")
			return nil
		},
	}
}
