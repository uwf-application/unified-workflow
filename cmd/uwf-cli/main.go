package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	version = "1.0.0"
	commit  = "dev"
	date    = "unknown"
)

func main() {
	rootCmd := &cobra.Command{
		Use:     "uwf-cli",
		Short:   "Unified Workflow CLI - Command line interface for workflow-api",
		Long:    `A comprehensive CLI tool for interacting with the Unified Workflow API.`,
		Version: fmt.Sprintf("%s (commit: %s, built: %s)", version, commit, date),
	}

	// Add global flags
	rootCmd.PersistentFlags().StringP("endpoint", "e", "https://af-test.qazpost.kz", "Workflow API endpoint")
	rootCmd.PersistentFlags().StringP("auth-token", "t", "", "Authentication token")
	rootCmd.PersistentFlags().StringP("output", "o", "json", "Output format (json, yaml, table)")
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "Verbose output")
	rootCmd.PersistentFlags().StringP("config", "c", "", "Config file path")

	// Add commands
	rootCmd.AddCommand(newHealthCmd())
	rootCmd.AddCommand(newConfigCmd())
	rootCmd.AddCommand(newWorkflowsCmd())
	rootCmd.AddCommand(newExecuteCmd())
	rootCmd.AddCommand(newExecutionsCmd())
	rootCmd.AddCommand(newTestCmd())
	rootCmd.AddCommand(newCompletionCmd())
	rootCmd.AddCommand(newDeployCmd())

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
