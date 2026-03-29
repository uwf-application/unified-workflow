package io.unifiedworkflow.sdk;

import com.fasterxml.jackson.annotation.JsonIgnoreProperties;
import com.fasterxml.jackson.annotation.JsonProperty;
import com.fasterxml.jackson.core.type.TypeReference;
import io.unifiedworkflow.sdk.errors.*;
import io.unifiedworkflow.sdk.internal.HttpClientWrapper;
import io.unifiedworkflow.sdk.models.*;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.time.Duration;
import java.time.Instant;
import java.util.Collections;
import java.util.List;
import java.util.Map;
import java.util.Objects;
import java.util.concurrent.CompletableFuture;
import java.util.concurrent.Executor;
import java.util.concurrent.Executors;

/**
 * Main entry point for the Unified Workflow Java SDK.
 *
 * <p>Create an instance with a static factory method, then call the synchronous or
 * asynchronous API methods. This class is thread-safe and intended to be shared
 * across the application as a long-lived singleton.</p>
 *
 * <pre>{@code
 * // Minimal setup
 * try (UnifiedWorkflowSDK sdk = UnifiedWorkflowSDK.create("http://workflow-api:8080")) {
 *     SDKExecuteWorkflowResponse resp = sdk.executeWorkflow(
 *         "antifraud",
 *         Map.of("transactionId", "txn-001", "amount", 1500.0)
 *     );
 *     ExecutionStatus status = sdk.waitForCompletion(resp.runId(), Duration.ofSeconds(30));
 * }
 * }</pre>
 *
 * <p>Implements {@link AutoCloseable} — use in try-with-resources or call {@link #close()}
 * explicitly to release internal resources.</p>
 */
public final class UnifiedWorkflowSDK implements AutoCloseable {

    private static final Logger log = LoggerFactory.getLogger(UnifiedWorkflowSDK.class);

    private static final String SDK_VERSION = "1.2.0";

    private final SDKConfig config;
    private final HttpClientWrapper httpClient;
    private final Executor asyncExecutor;

    // -------------------------------------------------------------------------
    // Static factories
    // -------------------------------------------------------------------------

    /**
     * Creates a client with a fully customised {@link SDKConfig}.
     *
     * @param config SDK configuration; must not be null
     * @return a new {@code UnifiedWorkflowSDK} instance
     */
    public static UnifiedWorkflowSDK create(SDKConfig config) {
        Objects.requireNonNull(config, "config must not be null");
        return new UnifiedWorkflowSDK(config);
    }

    /**
     * Creates a client using the supplied API endpoint and all other defaults.
     *
     * @param endpoint base URL of the workflow API (e.g. {@code "http://localhost:8080"}); must not be null
     * @return a new {@code UnifiedWorkflowSDK} instance
     */
    public static UnifiedWorkflowSDK create(String endpoint) {
        Objects.requireNonNull(endpoint, "endpoint must not be null");
        return new UnifiedWorkflowSDK(SDKConfig.builder()
                .workflowApiEndpoint(endpoint)
                .build());
    }

    // -------------------------------------------------------------------------
    // Constructor
    // -------------------------------------------------------------------------

    private UnifiedWorkflowSDK(SDKConfig config) {
        this.config = config;
        this.httpClient = new HttpClientWrapper(config);
        this.asyncExecutor = Executors.newVirtualThreadPerTaskExecutor();
    }

    // -------------------------------------------------------------------------
    // Workflow execution — synchronous
    // -------------------------------------------------------------------------

    /**
     * Executes a workflow with the given input data.
     *
     * <p>Constructs a default {@link SDKExecuteWorkflowRequest} with validation
     * and sanitization settings taken from the SDK configuration.</p>
     *
     * @param workflowId the registered workflow identifier; must not be null
     * @param inputData  key-value input data for the workflow; must not be null
     * @return the server's execution response including the run identifier
     * @throws WorkflowNotFoundException if the workflow is not registered
     * @throws WorkflowSDKException      on authentication, network, or server errors
     */
    public SDKExecuteWorkflowResponse executeWorkflow(String workflowId,
                                                      Map<String, Object> inputData) {
        Objects.requireNonNull(workflowId, "workflowId must not be null");
        Objects.requireNonNull(inputData, "inputData must not be null");

        SDKExecuteWorkflowRequest request = SDKExecuteWorkflowRequest.builder()
                .inputData(inputData)
                .enableValidation(config.enableValidation())
                .enableSanitization(config.enableSanitization())
                .priority(config.defaultPriority())
                .build();

        return executeWorkflow(workflowId, request);
    }

    /**
     * Executes a workflow using the supplied request object.
     *
     * @param workflowId the registered workflow identifier; must not be null
     * @param request    fully constructed execution request; must not be null
     * @return the server's execution response
     * @throws WorkflowNotFoundException if the workflow is not registered
     * @throws ValidationException       if the request payload fails validation
     * @throws WorkflowSDKException      on authentication, network, or server errors
     */
    public SDKExecuteWorkflowResponse executeWorkflow(String workflowId,
                                                      SDKExecuteWorkflowRequest request) {
        Objects.requireNonNull(workflowId, "workflowId must not be null");
        Objects.requireNonNull(request, "request must not be null");

        String url = config.workflowApiEndpoint() + "/api/v1/workflows/"
                + workflowId + "/execute";
        log.debug("Executing workflow '{}' via {}", workflowId, url);

        return httpClient.post(url, request,
                new TypeReference<SDKExecuteWorkflowResponse>() {});
    }

    // -------------------------------------------------------------------------
    // Execution status
    // -------------------------------------------------------------------------

    /**
     * Returns the current status of a workflow execution.
     *
     * @param runId the execution run identifier; must not be null
     * @return current status snapshot
     * @throws ExecutionNotFoundException if no run with the given ID exists
     * @throws WorkflowSDKException       on other errors
     */
    public ExecutionStatus getExecutionStatus(String runId) {
        Objects.requireNonNull(runId, "runId must not be null");
        String url = config.workflowApiEndpoint() + "/api/v1/executions/" + runId;
        return httpClient.get(url, new TypeReference<ExecutionStatus>() {});
    }

    /**
     * Returns the full execution details for a completed or running workflow run.
     *
     * @param runId the execution run identifier; must not be null
     * @return execution details including per-step results
     * @throws ExecutionNotFoundException if no run with the given ID exists
     * @throws WorkflowSDKException       on other errors
     */
    public ExecutionDetails getExecutionDetails(String runId) {
        Objects.requireNonNull(runId, "runId must not be null");
        String url = config.workflowApiEndpoint() + "/api/v1/executions/" + runId;
        return httpClient.get(url, new TypeReference<ExecutionDetails>() {});
    }

    // -------------------------------------------------------------------------
    // Execution control
    // -------------------------------------------------------------------------

    /**
     * Requests cancellation of a running execution.
     *
     * @param runId the execution run identifier; must not be null
     * @return {@code true} if the cancellation request was accepted
     * @throws ExecutionNotFoundException if no run with the given ID exists
     * @throws WorkflowSDKException       on other errors
     */
    public boolean cancelExecution(String runId) {
        Objects.requireNonNull(runId, "runId must not be null");
        String url = config.workflowApiEndpoint() + "/api/v1/executions/" + runId + "/cancel";
        httpClient.post(url, null, new TypeReference<Void>() {});
        return true;
    }

    /**
     * Pauses a running execution.
     *
     * @param runId the execution run identifier; must not be null
     * @return {@code true} if the pause request was accepted
     * @throws ExecutionNotFoundException if no run with the given ID exists
     * @throws WorkflowSDKException       on other errors
     */
    public boolean pauseExecution(String runId) {
        Objects.requireNonNull(runId, "runId must not be null");
        String url = config.workflowApiEndpoint() + "/api/v1/executions/" + runId + "/pause";
        httpClient.post(url, null, new TypeReference<Void>() {});
        return true;
    }

    /**
     * Resumes a paused execution.
     *
     * @param runId the execution run identifier; must not be null
     * @return {@code true} if the resume request was accepted
     * @throws ExecutionNotFoundException if no run with the given ID exists
     * @throws WorkflowSDKException       on other errors
     */
    public boolean resumeExecution(String runId) {
        Objects.requireNonNull(runId, "runId must not be null");
        String url = config.workflowApiEndpoint() + "/api/v1/executions/" + runId + "/resume";
        httpClient.post(url, null, new TypeReference<Void>() {});
        return true;
    }

    // -------------------------------------------------------------------------
    // Batch execution
    // -------------------------------------------------------------------------

    /**
     * Executes multiple workflows in a single API call.
     *
     * @param request batch request; must not be null
     * @return batch response containing per-item results
     * @throws WorkflowSDKException on errors
     */
    public BatchExecuteWorkflowsResponse batchExecuteWorkflows(
            BatchExecuteWorkflowsRequest request) {
        Objects.requireNonNull(request, "request must not be null");
        String url = config.workflowApiEndpoint() + "/sdk/v1/workflows/batch/execute";
        return httpClient.post(url, request,
                new TypeReference<BatchExecuteWorkflowsResponse>() {});
    }

    // -------------------------------------------------------------------------
    // Workflow registry
    // -------------------------------------------------------------------------

    /**
     * Lists all registered workflow definitions.
     *
     * @return list of workflow definitions; never null
     * @throws WorkflowSDKException on errors
     */
    public List<WorkflowDefinition> listWorkflows() {
        String url = config.workflowApiEndpoint() + "/api/v1/workflows";
        // API returns {"workflows":[...], "count":N} — use a wrapper to unwrap
        WorkflowListEnvelope envelope = httpClient.get(url,
                new TypeReference<WorkflowListEnvelope>() {});
        if (envelope == null || envelope.workflows() == null) {
            return Collections.emptyList();
        }
        return envelope.workflows();
    }

    /** Private envelope type matching the API's list-workflows response shape. */
    @JsonIgnoreProperties(ignoreUnknown = true)
    private record WorkflowListEnvelope(
            @JsonProperty("workflows") List<WorkflowDefinition> workflows
    ) {}

    /**
     * Returns the definition for a single registered workflow.
     *
     * @param workflowId the workflow identifier; must not be null
     * @return the workflow definition
     * @throws WorkflowNotFoundException if the workflow is not registered
     * @throws WorkflowSDKException      on other errors
     */
    public WorkflowDefinition getWorkflow(String workflowId) {
        Objects.requireNonNull(workflowId, "workflowId must not be null");
        String url = config.workflowApiEndpoint() + "/api/v1/workflows/" + workflowId;
        return httpClient.get(url, new TypeReference<WorkflowDefinition>() {});
    }

    // -------------------------------------------------------------------------
    // Validation
    // -------------------------------------------------------------------------

    /**
     * Validates input data against a set of rules without executing a workflow.
     *
     * @param data  key-value data to validate; must not be null
     * @param rules list of validation rules to apply; must not be null
     * @return the validation result
     * @throws WorkflowSDKException on errors
     */
    public ValidationResult validate(Map<String, Object> data, List<ValidationRule> rules) {
        Objects.requireNonNull(data, "data must not be null");
        Objects.requireNonNull(rules, "rules must not be null");

        String url = config.workflowApiEndpoint() + "/sdk/v1/validate";
        Map<String, Object> body = Map.of("data", data, "rules", rules);
        return httpClient.post(url, body, new TypeReference<ValidationResult>() {});
    }

    // -------------------------------------------------------------------------
    // Webhooks
    // -------------------------------------------------------------------------

    /**
     * Registers a new webhook endpoint.
     *
     * @param config webhook configuration; must not be null
     * @return the registered configuration including the server-assigned {@code webhookId}
     * @throws WorkflowSDKException on errors
     */
    public WebhookConfiguration registerWebhook(WebhookConfiguration config) {
        Objects.requireNonNull(config, "config must not be null");
        String url = this.config.workflowApiEndpoint() + "/sdk/v1/webhooks";
        return httpClient.post(url, config, new TypeReference<WebhookConfiguration>() {});
    }

    /**
     * Unregisters a previously registered webhook.
     *
     * @param webhookId the server-assigned webhook identifier; must not be null
     * @return {@code true} if the webhook was removed
     * @throws WorkflowSDKException on errors
     */
    public boolean unregisterWebhook(String webhookId) {
        Objects.requireNonNull(webhookId, "webhookId must not be null");
        String url = config.workflowApiEndpoint() + "/sdk/v1/webhooks/" + webhookId;
        return httpClient.delete(url);
    }

    // -------------------------------------------------------------------------
    // Health check
    // -------------------------------------------------------------------------

    /**
     * Checks that the workflow API is reachable and healthy.
     *
     * @return {@code true} if the server responded with a 2xx status
     * @throws WorkflowSDKException if the server is unreachable or unhealthy
     */
    public boolean ping() {
        String url = config.workflowApiEndpoint() + "/health";
        httpClient.get(url, new TypeReference<Void>() {});
        return true;
    }

    // -------------------------------------------------------------------------
    // Polling helper
    // -------------------------------------------------------------------------

    /**
     * Blocks until the execution reaches a terminal state or the timeout expires.
     *
     * <p>The method polls {@link #getExecutionStatus(String)} at the interval configured
     * by {@link SDKConfig#pollIntervalMs()}.</p>
     *
     * @param runId   the execution run identifier; must not be null
     * @param maxWait maximum time to wait; must not be null
     * @return the terminal {@link ExecutionStatus}
     * @throws TimeoutException     if {@code maxWait} elapses before a terminal state is reached
     * @throws WorkflowSDKException on other errors
     */
    public ExecutionStatus waitForCompletion(String runId, Duration maxWait) {
        Objects.requireNonNull(runId, "runId must not be null");
        Objects.requireNonNull(maxWait, "maxWait must not be null");

        Instant deadline = Instant.now().plus(maxWait);

        while (Instant.now().isBefore(deadline)) {
            ExecutionStatus status = getExecutionStatus(runId);
            if (status.isTerminal()) {
                return status;
            }
            log.debug("Execution '{}' status='{}' progress={}%; polling again in {}ms",
                    runId, status.status(), status.progress(), config.pollIntervalMs());
            sleep(config.pollIntervalMs());
        }

        throw new TimeoutException(
                "Execution '" + runId + "' did not complete within " + maxWait);
    }

    // -------------------------------------------------------------------------
    // Async methods
    // -------------------------------------------------------------------------

    /**
     * Asynchronously executes a workflow with the given input data.
     *
     * @param workflowId the registered workflow identifier; must not be null
     * @param inputData  key-value input data; must not be null
     * @return a {@link CompletableFuture} that completes with the execution response
     */
    public CompletableFuture<SDKExecuteWorkflowResponse> executeWorkflowAsync(
            String workflowId, Map<String, Object> inputData) {
        Objects.requireNonNull(workflowId, "workflowId must not be null");
        Objects.requireNonNull(inputData, "inputData must not be null");
        return CompletableFuture.supplyAsync(
                () -> executeWorkflow(workflowId, inputData), asyncExecutor);
    }

    /**
     * Asynchronously executes a workflow using the supplied request object.
     *
     * @param workflowId the registered workflow identifier; must not be null
     * @param request    fully constructed execution request; must not be null
     * @return a {@link CompletableFuture} that completes with the execution response
     */
    public CompletableFuture<SDKExecuteWorkflowResponse> executeWorkflowAsync(
            String workflowId, SDKExecuteWorkflowRequest request) {
        Objects.requireNonNull(workflowId, "workflowId must not be null");
        Objects.requireNonNull(request, "request must not be null");
        return CompletableFuture.supplyAsync(
                () -> executeWorkflow(workflowId, request), asyncExecutor);
    }

    /**
     * Asynchronously retrieves the current execution status.
     *
     * @param runId the execution run identifier; must not be null
     * @return a {@link CompletableFuture} that completes with the status
     */
    public CompletableFuture<ExecutionStatus> getExecutionStatusAsync(String runId) {
        Objects.requireNonNull(runId, "runId must not be null");
        return CompletableFuture.supplyAsync(
                () -> getExecutionStatus(runId), asyncExecutor);
    }

    /**
     * Asynchronously retrieves the full execution details.
     *
     * @param runId the execution run identifier; must not be null
     * @return a {@link CompletableFuture} that completes with the execution details
     */
    public CompletableFuture<ExecutionDetails> getExecutionDetailsAsync(String runId) {
        Objects.requireNonNull(runId, "runId must not be null");
        return CompletableFuture.supplyAsync(
                () -> getExecutionDetails(runId), asyncExecutor);
    }

    // -------------------------------------------------------------------------
    // Resource management
    // -------------------------------------------------------------------------

    /**
     * Returns the SDK configuration this client was created with.
     *
     * @return the immutable configuration record
     */
    public SDKConfig getConfig() {
        return config;
    }

    /**
     * Returns the SDK version string.
     *
     * @return version string (e.g. {@code "1.2.0"})
     */
    public String getSdkVersion() {
        return SDK_VERSION;
    }

    /**
     * Releases internal resources held by this client.
     *
     * <p>After calling this method the client must not be used again.</p>
     */
    @Override
    public void close() {
        httpClient.close();
    }

    // -------------------------------------------------------------------------
    // Internal helpers
    // -------------------------------------------------------------------------

    private void sleep(long ms) {
        try {
            Thread.sleep(ms);
        } catch (InterruptedException e) {
            Thread.currentThread().interrupt();
            throw new WorkflowSDKException(ErrorCode.NETWORK_ERROR,
                    "Interrupted during polling wait", e);
        }
    }
}
