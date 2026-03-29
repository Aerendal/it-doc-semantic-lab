package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

// NewNormalizeCmd returns the 'normalize' subcommand tree.
func NewNormalizeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "normalize",
		Short: "Normalize document names and assign canonical IDs",
	}
	cmd.AddCommand(newNormalizePreviewCmd(), newNormalizeApplyCmd())
	return cmd
}

func newNormalizePreviewCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "preview",
		Short: "Preview normalization changes without applying them",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Fprintln(cmd.OutOrStdout(), "normalize preview [not yet implemented]")
			return nil
		},
	}
}

func newNormalizeApplyCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "apply",
		Short: "Apply normalization and persist results",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Fprintln(cmd.OutOrStdout(), "normalize apply [not yet implemented]")
			return nil
		},
	}
}
