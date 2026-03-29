package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/it-doc-semantic-lab/itdlab/internal/cli"
)

func main() {
	if err := newRootCmd().Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func newRootCmd() *cobra.Command {
	var dbPath string
	var logPath string

	root := &cobra.Command{
		Use:   "itdlab",
		Short: "IT Documentation Semantic Lab — analysis and audit CLI",
		Long: `itdlab is the CLI for the IT Documentation Semantic Lab.

It ingests, normalises, classifies, and relates IT documentation artifacts,
storing all state in SQLite and appending an append-only JSONL event log.`,
	}

	root.PersistentFlags().StringVar(&dbPath, "db", "db/semantic_index.sqlite", "path to SQLite database")
	root.PersistentFlags().StringVar(&logPath, "log", "runs/events.jsonl", "path to JSONL event log")

	root.AddCommand(
		cli.NewIngestCmd(),
		cli.NewNormalizeCmd(),
		cli.NewClassifyCmd(),
		cli.NewRelationsCmd(),
		cli.NewExportCmd(),
		cli.NewAuditCmd(),
	)

	return root
}
