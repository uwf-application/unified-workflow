package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/spf13/cobra"
)

func newExecutionsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "executions",
		Short: "Manage workflow executions",
		Long:  `List, monitor, and manage workflow executions.`,
	}

	cmd.AddCommand(newExecutionsListCmd())
	cmd.AddCommand(newExecutionsStatusCmd())
	cmd.AddCommand(newExecutionsResultCmd())
	cmd.AddCommand(newExecutionsDataCmd())
	cmd.AddCommand(newExecutionsMetricsCmd())
	cmd.AddCommand(newExecutionsCancelCmd())
	cmd.AddCommand(newExecutionsPauseCmd())
	cmd.AddCommand(newExecutionsResumeCmd())
	cmd.AddCommand(newExecutionsRetryCmd())
	cmd.AddCommand(newExecutionsWatchCmd())

	return cmd
}

func newExecutionsListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all executions",
		RunE:  runExecutionsListCmd,
	}

	cmd.Flags().IntP("limit", "l", 50, "Maximum number of executions to return")
	cmd.Flags().Int("offset", 0, "Offset for pagination")
	cmd.Flags().StringP("workflow-id", "w", "", "Filter by workflow ID")
	cmd.Flags().StringP("status", "s", "", "Filter by status")

	return cmd
}

func runExecutionsListCmd(cmd *cobra.Command, args []string) error {
	endpoint, _ := cmd.Flags().GetString("endpoint")
	limit, _ := cmd.Flags().GetInt("limit")
	offset, _ := cmd.Flags().GetInt("offset")
	workflowID, _ := cmd.Flags().GetString("workflow-id")
	status, _ := cmd.Flags().GetString("status")
	output, _ := cmd.Flags().GetString("output")

	url := fmt.Sprintf("%s/api/v1/executions?limit=%d&offset=%d", endpoint, limit, offset)
	if workflowID != "" {
		url += fmt.Sprintf("&workflow_id=%s", workflowID)
	}
	if status != "" {
		url += fmt.Sprintf("&status=%s", status)
	}

	tr := &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
	httpClient := &http.Client{Transport: tr}

	resp, err := httpClient.Get(url)
	if err != nil {
		return fmt.Errorf("failed to list executions: %v", err)
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

func newExecutionsStatusCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "status [run-id]",
		Short: "Get execution status",
		Args:  cobra.ExactArgs(1),
		RunE:  runExecutionsStatusCmd,
	}

	return cmd
}

func runExecutionsStatusCmd(cmd *cobra.Command, args []string) error {
	runID := args[0]
	endpoint, _ := cmd.Flags().GetString("endpoint")
	output, _ := cmd.Flags().GetString("output")

	url := fmt.Sprintf("%s/api/v1/executions/%s", endpoint, runID)

	tr := &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
	httpClient := &http.Client{Transport: tr}

	resp, err := httpClient.Get(url)
	if err != nil {
		return fmt.Errorf("failed to get execution status: %v", err)
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

func newExecutionsResultCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "result [run-id]",
		Short: "Get execution result",
		Args:  cobra.ExactArgs(1),
		RunE:  runExecutionsResultCmd,
	}

	cmd.Flags().IntP("wait-ms", "w", 0, "Wait time in milliseconds")
	cmd.Flags().BoolP("long-poll", "l", false, "Enable long polling")

	return cmd
}

func runExecutionsResultCmd(cmd *cobra.Command, args []string) error {
	runID := args[0]
	waitMs, _ := cmd.Flags().GetInt("wait-ms")
	longPoll, _ := cmd.Flags().GetBool("long-poll")
	endpoint, _ := cmd.Flags().GetString("endpoint")
	output, _ := cmd.Flags().GetString("output")

	url := fmt.Sprintf("%s/api/v1/executions/%s/result?wait_ms=%d&long_poll=%v", endpoint, runID, waitMs, longPoll)

	tr := &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
	httpClient := &http.Client{Transport: tr}

	resp, err := httpClient.Get(url)
	if err != nil {
		return fmt.Errorf("failed to get execution result: %v", err)
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

func newExecutionsDataCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "data [run-id]",
		Short: "Get execution data",
		Args:  cobra.ExactArgs(1),
		RunE:  runExecutionsDataCmd,
	}

	return cmd
}

func runExecutionsDataCmd(cmd *cobra.Command, args []string) error {
	runID := args[0]
	endpoint, _ := cmd.Flags().GetString("endpoint")
	output, _ := cmd.Flags().GetString("output")

	url := fmt.Sprintf("%s/api/v1/executions/%s/data", endpoint, runID)

	tr := &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
	httpClient := &http.Client{Transport: tr}

	resp, err := httpClient.Get(url)
	if err != nil {
		return fmt.Errorf("failed to get execution data: %v", err)
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

func newExecutionsMetricsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "metrics [run-id]",
		Short: "Get execution metrics",
		Args:  cobra.ExactArgs(1),
		RunE:  runExecutionsMetricsCmd,
	}

	return cmd
}

func runExecutionsMetricsCmd(cmd *cobra.Command, args []string) error {
	runID := args[0]
	endpoint, _ := cmd.Flags().GetString("endpoint")
	output, _ := cmd.Flags().GetString("output")

	url := fmt.Sprintf("%s/api/v1/executions/%s/metrics", endpoint, runID)

	tr := &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
	httpClient := &http.Client{Transport: tr}

	resp, err := httpClient.Get(url)
	if err != nil {
		return fmt.Errorf("failed to get execution metrics: %v", err)
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

func newExecutionsCancelCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cancel [run-id]",
		Short: "Cancel execution",
		Args:  cobra.ExactArgs(1),
		RunE:  runExecutionsCancelCmd,
	}

	return cmd
}

func runExecutionsCancelCmd(cmd *cobra.Command, args []string) error {
	runID := args[0]
	endpoint, _ := cmd.Flags().GetString("endpoint")
	output, _ := cmd.Flags().GetString("output")

	url := fmt.Sprintf("%s/api/v1/executions/%s/cancel", endpoint, runID)

	tr := &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
	httpClient := &http.Client{Transport: tr}

	resp, err := httpClient.Post(url, "application/json", bytes.NewBufferString("{}"))
	if err != nil {
		return fmt.Errorf("failed to cancel execution: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	var response map[string]interface{}
	if resp.StatusCode == http.StatusNoContent {
		response = map[string]interface{}{
			"message": "cancelled",
			"run_id":  runID,
		}
	} else {
		if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
			return fmt.Errorf("failed to decode response: %v", err)
		}
	}

	return printOutput(response, output)
}

func newExecutionsPauseCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pause [run-id]",
		Short: "Pause execution",
		Args:  cobra.ExactArgs(1),
		RunE:  runExecutionsPauseCmd,
	}

	return cmd
}

func runExecutionsPauseCmd(cmd *cobra.Command, args []string) error {
	runID := args[0]
	endpoint, _ := cmd.Flags().GetString("endpoint")
	output, _ := cmd.Flags().GetString("output")

	url := fmt.Sprintf("%s/api/v1/executions/%s/pause", endpoint, runID)

	tr := &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
	httpClient := &http.Client{Transport: tr}

	resp, err := httpClient.Post(url, "application/json", bytes.NewBufferString("{}"))
	if err != nil {
		return fmt.Errorf("failed to pause execution: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	var response map[string]interface{}
	if resp.StatusCode == http.StatusNoContent {
		response = map[string]interface{}{
			"message": "paused",
			"run_id":  runID,
		}
	} else {
		if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
			return fmt.Errorf("failed to decode response: %v", err)
		}
	}

	return printOutput(response, output)
}

func newExecutionsResumeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "resume [run-id]",
		Short: "Resume execution",
		Args:  cobra.ExactArgs(1),
		RunE:  runExecutionsResumeCmd,
	}

	return cmd
}

func runExecutionsResumeCmd(cmd *cobra.Command, args []string) error {
	runID := args[0]
	endpoint, _ := cmd.Flags().GetString("endpoint")
	output, _ := cmd.Flags().GetString("output")

	url := fmt.Sprintf("%s/api/v1/executions/%s/resume", endpoint, runID)

	tr := &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
	httpClient := &http.Client{Transport: tr}

	resp, err := httpClient.Post(url, "application/json", bytes.NewBufferString("{}"))
	if err != nil {
		return fmt.Errorf("failed to resume execution: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	var response map[string]interface{}
	if resp.StatusCode == http.StatusNoContent {
		response = map[string]interface{}{
			"message": "resumed",
			"run_id":  runID,
		}
	} else {
		if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
			return fmt.Errorf("failed to decode response: %v", err)
		}
	}

	return printOutput(response, output)
}

func newExecutionsRetryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "retry [run-id]",
		Short: "Retry execution",
		Args:  cobra.ExactArgs(1),
		RunE:  runExecutionsRetryCmd,
	}

	return cmd
}

func runExecutionsRetryCmd(cmd *cobra.Command, args []string) error {
	runID := args[0]
	endpoint, _ := cmd.Flags().GetString("endpoint")
	output, _ := cmd.Flags().GetString("output")

	url := fmt.Sprintf("%s/api/v1/executions/%s/retry", endpoint, runID)

	tr := &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
	httpClient := &http.Client{Transport: tr}

	resp, err := httpClient.Post(url, "application/json", bytes.NewBufferString("{}"))
	if err != nil {
		return fmt.Errorf("failed to retry execution: %v", err)
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

func newExecutionsWatchCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "watch [run-id]",
		Short: "Watch execution status in real-time",
		Args:  cobra.ExactArgs(1),
		RunE:  runExecutionsWatchCmd,
	}

	cmd.Flags().DurationP("interval", "i", 2*time.Second, "Polling interval")
	cmd.Flags().DurationP("timeout", "t", 5*time.Minute, "Maximum watch time")
	cmd.Flags().BoolP("exit-on-completion", "e", true, "Exit when execution completes")

	return cmd
}

func runExecutionsWatchCmd(cmd *cobra.Command, args []string) error {
	runID := args[0]
	endpoint, _ := cmd.Flags().GetString("endpoint")
	interval, _ := cmd.Flags().GetDuration("interval")
	timeout, _ := cmd.Flags().GetDuration("timeout")
	output, _ := cmd.Flags().GetString("output")

	url := fmt.Sprintf("%s/api/v1/executions/%s", endpoint, runID)

	tr := &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
	httpClient := &http.Client{Transport: tr}

	deadline := time.After(timeout)
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		resp, err := httpClient.Get(url)
		if err != nil {
			return fmt.Errorf("failed to get execution status: %v", err)
		}

		if resp.StatusCode != http.StatusOK {
			resp.Body.Close()
			return fmt.Errorf("API returned status %d", resp.StatusCode)
		}

		var response map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
			resp.Body.Close()
			return fmt.Errorf("failed to decode response: %v", err)
		}
		resp.Body.Close()

		if err := printOutput(response, output); err != nil {
			return err
		}

		if st, ok := response["status"].(string); ok {
			if st == "completed" || st == "failed" || st == "cancelled" {
				return nil
			}
		}

		select {
		case <-deadline:
			return fmt.Errorf("watch timed out after %v", timeout)
		case <-ticker.C:
		}
	}
}
