package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/spf13/cobra"
)

func newExecuteCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "execute",
		Short: "Execute workflows",
		Long:  `Execute workflows synchronously or asynchronously.`,
	}

	cmd.AddCommand(newExecuteSyncCmd())
	cmd.AddCommand(newExecuteAsyncCmd())
	cmd.AddCommand(newExecuteBulkCmd())

	return cmd
}

func newExecuteSyncCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sync [workflow-id]",
		Short: "Execute workflow synchronously",
		Args:  cobra.ExactArgs(1),
		RunE:  runExecuteSyncCmd,
	}

	cmd.Flags().StringP("input", "i", "", "Input data as JSON string")
	cmd.Flags().StringP("input-file", "f", "", "Input data from JSON file")
	cmd.Flags().Int("timeout", 30, "Timeout in seconds")

	return cmd
}

func runExecuteSyncCmd(cmd *cobra.Command, args []string) error {
	workflowID := args[0]
	input, _ := cmd.Flags().GetString("input")
	inputFile, _ := cmd.Flags().GetString("input-file")
	timeout, _ := cmd.Flags().GetInt("timeout")
	output, _ := cmd.Flags().GetString("output")
	endpoint, _ := cmd.Flags().GetString("endpoint")

	var inputData map[string]interface{}
	if inputFile != "" {
		data, err := os.ReadFile(inputFile)
		if err != nil {
			return fmt.Errorf("failed to read input file: %v", err)
		}
		if err := json.Unmarshal(data, &inputData); err != nil {
			return fmt.Errorf("failed to parse input JSON: %v", err)
		}
	} else if input != "" {
		if err := json.Unmarshal([]byte(input), &inputData); err != nil {
			return fmt.Errorf("failed to parse input JSON: %v", err)
		}
	} else {
		inputData = map[string]interface{}{
			"timestamp": time.Now().Format(time.RFC3339),
			"source":    "uwf-cli",
		}
	}

	body, err := json.Marshal(map[string]interface{}{
		"input_data": inputData,
	})
	if err != nil {
		return fmt.Errorf("failed to marshal request body: %v", err)
	}

	tr := &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
	httpClient := &http.Client{Transport: tr}

	url := fmt.Sprintf("%s/api/v1/workflows/%s/execute", endpoint, workflowID)
	resp, err := httpClient.Post(url, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("failed to execute workflow: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusAccepted {
		return fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	var execResponse map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&execResponse); err != nil {
		return fmt.Errorf("failed to decode response: %v", err)
	}

	runID := ""
	if v, ok := execResponse["run_id"]; ok {
		runID, _ = v.(string)
	} else if v, ok := execResponse["runId"]; ok {
		runID, _ = v.(string)
	}

	if runID == "" {
		return printOutput(execResponse, output)
	}

	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	var lastResponse map[string]interface{}
	lastResponse = execResponse

	for time.Now().Before(deadline) {
		time.Sleep(1 * time.Second)

		pollURL := fmt.Sprintf("%s/api/v1/executions/%s", endpoint, runID)
		pollResp, err := httpClient.Get(pollURL)
		if err != nil {
			return fmt.Errorf("failed to poll execution status: %v", err)
		}

		var pollData map[string]interface{}
		decodeErr := json.NewDecoder(pollResp.Body).Decode(&pollData)
		pollResp.Body.Close()

		if decodeErr != nil {
			return fmt.Errorf("failed to decode poll response: %v", decodeErr)
		}

		lastResponse = pollData

		if status, ok := pollData["status"].(string); ok {
			if status == "completed" || status == "failed" || status == "cancelled" {
				return printOutput(lastResponse, output)
			}
		}
	}

	lastResponse["_warning"] = fmt.Sprintf("timeout of %d seconds elapsed; showing last known status", timeout)
	return printOutput(lastResponse, output)
}

func newExecuteAsyncCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "async [workflow-id]",
		Short: "Execute workflow asynchronously",
		Args:  cobra.ExactArgs(1),
		RunE:  runExecuteAsyncCmd,
	}

	cmd.Flags().StringP("input", "i", "", "Input data as JSON string")
	cmd.Flags().StringP("input-file", "f", "", "Input data from JSON file")
	cmd.Flags().Int("timeout", 30000, "Timeout in milliseconds")
	cmd.Flags().String("callback-url", "", "Callback URL for completion notification")
	cmd.Flags().BoolP("wait", "w", false, "Wait for completion")

	return cmd
}

func runExecuteAsyncCmd(cmd *cobra.Command, args []string) error {
	workflowID := args[0]
	input, _ := cmd.Flags().GetString("input")
	inputFile, _ := cmd.Flags().GetString("input-file")
	timeout, _ := cmd.Flags().GetInt("timeout")
	callbackURL, _ := cmd.Flags().GetString("callback-url")
	wait, _ := cmd.Flags().GetBool("wait")
	output, _ := cmd.Flags().GetString("output")
	endpoint, _ := cmd.Flags().GetString("endpoint")

	var inputData map[string]interface{}
	if inputFile != "" {
		data, err := os.ReadFile(inputFile)
		if err != nil {
			return fmt.Errorf("failed to read input file: %v", err)
		}
		if err := json.Unmarshal(data, &inputData); err != nil {
			return fmt.Errorf("failed to parse input JSON: %v", err)
		}
	} else if input != "" {
		if err := json.Unmarshal([]byte(input), &inputData); err != nil {
			return fmt.Errorf("failed to parse input JSON: %v", err)
		}
	} else {
		inputData = map[string]interface{}{
			"timestamp": time.Now().Format(time.RFC3339),
			"source":    "uwf-cli",
			"async":     true,
		}
	}

	payload := map[string]interface{}{
		"input_data":   inputData,
		"callback_url": callbackURL,
		"timeout_ms":   timeout,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal request body: %v", err)
	}

	tr := &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
	httpClient := &http.Client{Transport: tr}

	url := fmt.Sprintf("%s/api/v1/workflows/%s/execute", endpoint, workflowID)
	resp, err := httpClient.Post(url, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("failed to execute workflow: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusAccepted {
		return fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	var execResponse map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&execResponse); err != nil {
		return fmt.Errorf("failed to decode response: %v", err)
	}

	if !wait {
		return printOutput(execResponse, output)
	}

	runID := ""
	if v, ok := execResponse["run_id"]; ok {
		runID, _ = v.(string)
	} else if v, ok := execResponse["runId"]; ok {
		runID, _ = v.(string)
	}

	if runID == "" {
		return printOutput(execResponse, output)
	}

	deadline := time.Now().Add(time.Duration(timeout) * time.Millisecond)
	lastResponse := execResponse

	for time.Now().Before(deadline) {
		time.Sleep(2 * time.Second)

		pollURL := fmt.Sprintf("%s/api/v1/executions/%s", endpoint, runID)
		pollResp, err := httpClient.Get(pollURL)
		if err != nil {
			return fmt.Errorf("failed to poll execution status: %v", err)
		}

		var pollData map[string]interface{}
		decodeErr := json.NewDecoder(pollResp.Body).Decode(&pollData)
		pollResp.Body.Close()

		if decodeErr != nil {
			return fmt.Errorf("failed to decode poll response: %v", decodeErr)
		}

		lastResponse = pollData

		if status, ok := pollData["status"].(string); ok {
			if status == "completed" || status == "failed" || status == "cancelled" {
				return printOutput(lastResponse, output)
			}
		}
	}

	return printOutput(lastResponse, output)
}

func newExecuteBulkCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "bulk [file]",
		Short: "Execute multiple workflows from file",
		Args:  cobra.ExactArgs(1),
		RunE:  runExecuteBulkCmd,
	}

	cmd.Flags().BoolP("parallel", "p", false, "Execute in parallel")
	cmd.Flags().Int("concurrency", 5, "Maximum concurrent executions")
	cmd.Flags().BoolP("async", "a", true, "Execute asynchronously")

	return cmd
}

type bulkEntry struct {
	WorkflowID string                 `json:"workflow_id"`
	InputData  map[string]interface{} `json:"input_data"`
}

type bulkResult struct {
	WorkflowID string `json:"workflow_id"`
	RunID      string `json:"run_id,omitempty"`
	Status     string `json:"status"`
	Error      string `json:"error,omitempty"`
}

func runExecuteBulkCmd(cmd *cobra.Command, args []string) error {
	filename := args[0]
	parallel, _ := cmd.Flags().GetBool("parallel")
	concurrency, _ := cmd.Flags().GetInt("concurrency")
	output, _ := cmd.Flags().GetString("output")
	endpoint, _ := cmd.Flags().GetString("endpoint")

	data, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("failed to read bulk file: %v", err)
	}

	var entries []bulkEntry
	if err := json.Unmarshal(data, &entries); err != nil {
		return fmt.Errorf("failed to parse bulk file JSON: %v", err)
	}

	tr := &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
	httpClient := &http.Client{Transport: tr}

	submitOne := func(entry bulkEntry) bulkResult {
		body, err := json.Marshal(map[string]interface{}{
			"input_data": entry.InputData,
		})
		if err != nil {
			return bulkResult{WorkflowID: entry.WorkflowID, Status: "failed", Error: err.Error()}
		}

		url := fmt.Sprintf("%s/api/v1/workflows/%s/execute", endpoint, entry.WorkflowID)
		resp, err := httpClient.Post(url, "application/json", bytes.NewReader(body))
		if err != nil {
			return bulkResult{WorkflowID: entry.WorkflowID, Status: "failed", Error: err.Error()}
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusAccepted {
			return bulkResult{
				WorkflowID: entry.WorkflowID,
				Status:     "failed",
				Error:      fmt.Sprintf("API returned status %d", resp.StatusCode),
			}
		}

		var execResponse map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&execResponse); err != nil {
			return bulkResult{WorkflowID: entry.WorkflowID, Status: "failed", Error: err.Error()}
		}

		runID := ""
		if v, ok := execResponse["run_id"]; ok {
			runID, _ = v.(string)
		} else if v, ok := execResponse["runId"]; ok {
			runID, _ = v.(string)
		}

		status := "submitted"
		if v, ok := execResponse["status"].(string); ok {
			status = v
		}

		return bulkResult{WorkflowID: entry.WorkflowID, RunID: runID, Status: status}
	}

	results := make([]bulkResult, len(entries))

	if parallel {
		sem := make(chan struct{}, concurrency)
		var wg sync.WaitGroup
		var mu sync.Mutex

		for idx, entry := range entries {
			wg.Add(1)
			sem <- struct{}{}
			go func(i int, e bulkEntry) {
				defer wg.Done()
				defer func() { <-sem }()
				r := submitOne(e)
				mu.Lock()
				results[i] = r
				mu.Unlock()
			}(idx, entry)
		}
		wg.Wait()
	} else {
		for i, entry := range entries {
			results[i] = submitOne(entry)
		}
	}

	submitted := 0
	failed := 0
	for _, r := range results {
		if r.Error != "" {
			failed++
		} else {
			submitted++
		}
	}

	summary := map[string]interface{}{
		"total":     len(entries),
		"submitted": submitted,
		"failed":    failed,
		"results":   results,
	}

	return printOutput(summary, output)
}
