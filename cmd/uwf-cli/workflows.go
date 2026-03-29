package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/spf13/cobra"
)

func newWorkflowsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "workflows",
		Short: "Manage workflows",
		Long:  `Create, list, get, and delete workflows.`,
	}

	cmd.AddCommand(newWorkflowsListCmd())
	cmd.AddCommand(newWorkflowsCreateCmd())
	cmd.AddCommand(newWorkflowsGetCmd())
	cmd.AddCommand(newWorkflowsDeleteCmd())

	return cmd
}

func newWorkflowsListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all workflows",
		RunE:  runWorkflowsListCmd,
	}

	cmd.Flags().IntP("limit", "l", 100, "Maximum number of workflows to return")
	cmd.Flags().Int("offset", 0, "Offset for pagination")

	return cmd
}

func runWorkflowsListCmd(cmd *cobra.Command, args []string) error {
	endpoint, _ := cmd.Flags().GetString("endpoint")
	limit, _ := cmd.Flags().GetInt("limit")
	offset, _ := cmd.Flags().GetInt("offset")
	output, _ := cmd.Flags().GetString("output")

	// Build URL with query parameters
	url := fmt.Sprintf("%s/api/v1/workflows", endpoint)
	if limit > 0 || offset > 0 {
		url += fmt.Sprintf("?limit=%d&offset=%d", limit, offset)
	}

	// Create HTTP client that skips TLS verification for self-signed certificates
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	// Make HTTP request
	resp, err := client.Get(url)
	if err != nil {
		return fmt.Errorf("failed to list workflows: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	var response map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return fmt.Errorf("failed to decode response: %v", err)
	}

	return printOutput(response, output)
}

func newWorkflowsCreateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new workflow",
		RunE:  runWorkflowsCreateCmd,
	}

	cmd.Flags().StringP("name", "n", "", "Workflow name (required)")
	cmd.Flags().StringP("description", "d", "", "Workflow description")
	cmd.Flags().StringP("input-file", "f", "", "JSON file with workflow definition")

	cmd.MarkFlagRequired("name")

	return cmd
}

func runWorkflowsCreateCmd(cmd *cobra.Command, args []string) error {
	name, _ := cmd.Flags().GetString("name")
	description, _ := cmd.Flags().GetString("description")
	output, _ := cmd.Flags().GetString("output")
	endpoint, _ := cmd.Flags().GetString("endpoint")

	// Create workflow payload
	workflowData := map[string]interface{}{
		"name":        name,
		"description": description,
		"steps":       []interface{}{},
	}

	data, err := json.Marshal(workflowData)
	if err != nil {
		return fmt.Errorf("failed to marshal workflow data: %v", err)
	}

	// Create HTTP client that skips TLS verification for self-signed certificates
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	// Make HTTP request
	url := fmt.Sprintf("%s/api/v1/workflows", endpoint)
	resp, err := client.Post(url, "application/json", bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("failed to create workflow: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	var response map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return fmt.Errorf("failed to decode response: %v", err)
	}

	return printOutput(response, output)
}

func newWorkflowsGetCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get [workflow-id]",
		Short: "Get workflow details",
		Args:  cobra.ExactArgs(1),
		RunE:  runWorkflowsGetCmd,
	}

	return cmd
}

func runWorkflowsGetCmd(cmd *cobra.Command, args []string) error {
	workflowID := args[0]
	output, _ := cmd.Flags().GetString("output")
	endpoint, _ := cmd.Flags().GetString("endpoint")

	// Create HTTP client that skips TLS verification for self-signed certificates
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	// Make HTTP request
	url := fmt.Sprintf("%s/api/v1/workflows/%s", endpoint, workflowID)
	resp, err := client.Get(url)
	if err != nil {
		return fmt.Errorf("failed to get workflow: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	var response map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return fmt.Errorf("failed to decode response: %v", err)
	}

	return printOutput(response, output)
}

func newWorkflowsDeleteCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete [workflow-id]",
		Short: "Delete a workflow",
		Args:  cobra.ExactArgs(1),
		RunE:  runWorkflowsDeleteCmd,
	}

	return cmd
}

func runWorkflowsDeleteCmd(cmd *cobra.Command, args []string) error {
	workflowID := args[0]
	output, _ := cmd.Flags().GetString("output")
	endpoint, _ := cmd.Flags().GetString("endpoint")

	// Create HTTP client that skips TLS verification for self-signed certificates
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	// Make HTTP request
	url := fmt.Sprintf("%s/api/v1/workflows/%s", endpoint, workflowID)
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create delete request: %v", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to delete workflow: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	var response map[string]interface{}
	if resp.StatusCode == http.StatusOK {
		if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
			return fmt.Errorf("failed to decode response: %v", err)
		}
	} else {
		// For 204 No Content, create a simple response
		response = map[string]interface{}{
			"message": "Workflow deleted successfully",
			"id":      workflowID,
		}
	}

	return printOutput(response, output)
}

// Helper function to print output in different formats
func printOutput(data interface{}, format string) error {
	switch format {
	case "json":
		encoder := json.NewEncoder(os.Stdout)
		encoder.SetIndent("", "  ")
		return encoder.Encode(data)
	case "yaml":
		// For now, fall back to JSON
		encoder := json.NewEncoder(os.Stdout)
		encoder.SetIndent("", "  ")
		return encoder.Encode(data)
	case "table":
		// Simple table output
		if m, ok := data.(map[string]interface{}); ok {
			for k, v := range m {
				fmt.Printf("%-20s: %v\n", k, v)
			}
		} else {
			fmt.Printf("%v\n", data)
		}
		return nil
	default:
		encoder := json.NewEncoder(os.Stdout)
		encoder.SetIndent("", "  ")
		return encoder.Encode(data)
	}
}
