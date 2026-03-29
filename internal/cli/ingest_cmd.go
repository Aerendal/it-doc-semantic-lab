package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

// NewIngestCmd returns the 'ingest' subcommand tree.
func NewIngestCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ingest",
		Short: "Parse and store raw IT documentation sources",
	}
	cmd.AddCommand(newIngestRunCmd(), newIngestInspectCmd())
	return cmd
}

func newIngestRunCmd() *cobra.Command {
	var sourcePath string
	cmd := &cobra.Command{
		Use:   "run",
		Short: "Run full ingest pipeline on source directory",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Fprintf(cmd.OutOrStdout(), "ingest run: source=%s [not yet implemented]\n", sourcePath)
			return nil
		},
	}
	cmd.Flags().StringVarP(&sourcePath, "source", "s", "sources/", "path to source directory")
	return cmd
}

func newIngestInspectCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "inspect [source-path]",
		Short: "Inspect a source file without persisting",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			path := ""
			if len(args) > 0 {
				path = args[0]
			}
			fmt.Fprintf(cmd.OutOrStdout(), "ingest inspect: path=%q [not yet implemented]\n", path)
			return nil
		},
	}
}
