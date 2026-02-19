package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

func newConfigCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Manage CLI configuration",
		Long:  `View and manage CLI configuration settings.`,
	}

	cmd.AddCommand(newConfigShowCmd())
	cmd.AddCommand(newConfigSetCmd())
	cmd.AddCommand(newConfigInitCmd())

	return cmd
}

func newConfigShowCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show",
		Short: "Show current configuration",
		RunE:  runConfigShowCmd,
	}

	return cmd
}

func runConfigShowCmd(cmd *cobra.Command, args []string) error {
	endpoint, _ := cmd.Flags().GetString("endpoint")
	authToken, _ := cmd.Flags().GetString("auth-token")
	output, _ := cmd.Flags().GetString("output")
	verbose, _ := cmd.Flags().GetBool("verbose")

	config := map[string]interface{}{
		"endpoint":      endpoint,
		"auth_token":    maskToken(authToken),
		"output":        output,
		"verbose":       verbose,
		"config_file":   getConfigFilePath(),
		"config_exists": configFileExists(),
	}

	return printOutput(config, output)
}

func newConfigSetCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set [key] [value]",
		Short: "Set configuration value",
		Args:  cobra.ExactArgs(2),
		RunE:  runConfigSetCmd,
	}

	return cmd
}

func runConfigSetCmd(cmd *cobra.Command, args []string) error {
	key := args[0]
	value := args[1]
	output, _ := cmd.Flags().GetString("output")

	response := map[string]interface{}{
		"message": "Configuration updated",
		"key":     key,
		"value":   value,
		"note":    "Configuration persistence not yet implemented",
	}

	return printOutput(response, output)
}

func newConfigInitCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize configuration file",
		RunE:  runConfigInitCmd,
	}

	return cmd
}

func runConfigInitCmd(cmd *cobra.Command, args []string) error {
	output, _ := cmd.Flags().GetString("output")

	configDir := filepath.Join(os.Getenv("HOME"), ".uwf-cli")
	configFile := filepath.Join(configDir, "config.yaml")

	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %v", err)
	}

	// Create default config
	defaultConfig := `# Unified Workflow CLI Configuration
endpoint: "https://af-test.qazpost.kz"
auth_token: ""
timeout: 30
max_retries: 3
output_format: "json"
default_workflow_name: "Test Workflow"
test_data_path: "./test-data"
`

	if err := os.WriteFile(configFile, []byte(defaultConfig), 0644); err != nil {
		return fmt.Errorf("failed to write config file: %v", err)
	}

	response := map[string]interface{}{
		"message":     "Configuration file initialized",
		"config_file": configFile,
		"config_dir":  configDir,
	}

	return printOutput(response, output)
}

func getConfigFilePath() string {
	return filepath.Join(os.Getenv("HOME"), ".uwf-cli", "config.yaml")
}

func configFileExists() bool {
	_, err := os.Stat(getConfigFilePath())
	return err == nil
}

func maskToken(token string) string {
	if token == "" {
		return ""
	}
	if len(token) <= 8 {
		return "***"
	}
	return token[:4] + "***" + token[len(token)-4:]
}
