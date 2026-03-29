package io.unifiedworkflow.sdk;

import io.unifiedworkflow.sdk.errors.*;
import io.unifiedworkflow.sdk.internal.HttpClientWrapper;
import io.unifiedworkflow.sdk.models.*;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.ExtendWith;
import org.mockito.Mock;
import org.mockito.MockitoAnnotations;

import java.time.Duration;
import java.time.Instant;
import java.util.List;
import java.util.Map;

import static org.junit.jupiter.api.Assertions.*;

/**
 * Unit tests for {@link UnifiedWorkflowSDK}.
 *
 * <p>HTTP interactions are verified via Mockito. The SDK uses an internal
 * {@link HttpClientWrapper} which is replaced with a mock in tests that need
 * to control responses. Tests that only exercise pure Java logic (builders,
 * records, config) do not require mocking.</p>
 */
@DisplayName("UnifiedWorkflowSDK")
class UnifiedWorkflowSDKTest {

    // -------------------------------------------------------------------------
    // Helpers for creating test fixtures
    // -------------------------------------------------------------------------

    private static SDKConfig defaultConfig() {
        return SDKConfig.builder()
                .workflowApiEndpoint("http://localhost:8080")
                .timeout(Duration.ofSeconds(5))
                .maxRetries(0)
                .build();
    }

    private static ExecutionStatus pendingStatus(String runId) {
        return new ExecutionStatus(runId, "pending", 0, null,
                false, null, Instant.now(), null);
    }

    private static ExecutionStatus completedStatus(String runId) {
        return new ExecutionStatus(runId, "completed", 100, null,
                true, null, Instant.now(), Instant.now());
    }

    private static ExecutionStatus runningStatus(String runId) {
        return new ExecutionStatus(runId, "running", 50, "step-1",
                false, null, Instant.now(), null);
    }

    // =========================================================================
    // SDKConfig builder tests
    // =========================================================================

    @Test
    @DisplayName("SDKConfig builder applies defaults")
    void sdkConfig_defaults() {
        SDKConfig config = SDKConfig.builder().build();

        assertEquals("http://localhost:8080", config.workflowApiEndpoint());
        assertEquals(Duration.ofSeconds(30), config.timeout());
        assertEquals(3, config.maxRetries());
        assertEquals(Duration.ofSeconds(1), config.retryDelay());
        assertEquals(SDKConfig.AuthType.NONE, config.authType());
        assertTrue(config.enableValidation());
        assertTrue(config.enableSanitization());
        assertFalse(config.strictValidation());
        assertTrue(config.asyncExecution());
        assertEquals(2000L, config.pollIntervalMs());
        assertEquals(5, config.defaultPriority());
        assertTrue(config.enableCircuitBreaker());
        assertEquals(5, config.circuitBreakerThreshold());
        assertEquals(Duration.ofSeconds(60), config.circuitBreakerTimeout());
        assertFalse(config.enableRequestLogging());
    }

    @Test
    @DisplayName("SDKConfig builder respects custom values")
    void sdkConfig_customValues() {
        SDKConfig config = SDKConfig.builder()
                .workflowApiEndpoint("https://workflow.example.com")
                .timeout(Duration.ofSeconds(60))
                .maxRetries(5)
                .authToken("my-token")
                .authType(SDKConfig.AuthType.BEARER_TOKEN)
                .enableRequestLogging(true)
                .defaultPriority(8)
                .build();

        assertEquals("https://workflow.example.com", config.workflowApiEndpoint());
        assertEquals(Duration.ofSeconds(60), config.timeout());
        assertEquals(5, config.maxRetries());
        assertEquals("my-token", config.authToken());
        assertEquals(SDKConfig.AuthType.BEARER_TOKEN, config.authType());
        assertTrue(config.enableRequestLogging());
        assertEquals(8, config.defaultPriority());
    }

    @Test
    @DisplayName("SDKConfig strips trailing slash from endpoint")
    void sdkConfig_stripsTrailingSlash() {
        SDKConfig config = SDKConfig.builder()
                .workflowApiEndpoint("http://localhost:8080/")
                .build();
        assertEquals("http://localhost:8080", config.workflowApiEndpoint());
    }

    @Test
    @DisplayName("SDKConfigBuilder rejects invalid defaultPriority")
    void sdkConfigBuilder_rejectsInvalidPriority() {
        assertThrows(IllegalArgumentException.class,
                () -> SDKConfig.builder().defaultPriority(0).build());
        assertThrows(IllegalArgumentException.class,
                () -> SDKConfig.builder().defaultPriority(11).build());
    }

    // =========================================================================
    // SDK static factory tests
    // =========================================================================

    @Test
    @DisplayName("create(endpoint) produces a usable client")
    void create_withEndpoint() {
        try (UnifiedWorkflowSDK sdk = UnifiedWorkflowSDK.create("http://localhost:8080")) {
            assertNotNull(sdk);
            assertEquals("http://localhost:8080", sdk.getConfig().workflowApiEndpoint());
            assertEquals("1.2.0", sdk.getSdkVersion());
        }
    }

    @Test
    @DisplayName("create(SDKConfig) produces a usable client")
    void create_withConfig() {
        SDKConfig config = SDKConfig.builder()
                .workflowApiEndpoint("http://test-host:9090")
                .authToken("tok")
                .authType(SDKConfig.AuthType.API_KEY)
                .build();

        try (UnifiedWorkflowSDK sdk = UnifiedWorkflowSDK.create(config)) {
            assertNotNull(sdk);
            assertEquals("http://test-host:9090", sdk.getConfig().workflowApiEndpoint());
            assertEquals(SDKConfig.AuthType.API_KEY, sdk.getConfig().authType());
        }
    }

    @Test
    @DisplayName("create(null) throws NullPointerException")
    void create_nullEndpointThrows() {
        assertThrows(NullPointerException.class,
                () -> UnifiedWorkflowSDK.create((String) null));
        assertThrows(NullPointerException.class,
                () -> UnifiedWorkflowSDK.create((SDKConfig) null));
    }

    // =========================================================================
    // ValidationRule factory method tests
    // =========================================================================

    @Test
    @DisplayName("ValidationRule.required creates required string rule")
    void validationRule_required() {
        ValidationRule rule = ValidationRule.required("transactionId");

        assertEquals("transactionId", rule.field());
        assertEquals("string", rule.ruleType());
        assertTrue(rule.required());
        assertNull(rule.minLength());
        assertNull(rule.maxLength());
        assertNull(rule.pattern());
        assertNotNull(rule.allowedValues());
        assertTrue(rule.allowedValues().isEmpty());
    }

    @Test
    @DisplayName("ValidationRule.email creates email rule")
    void validationRule_email() {
        ValidationRule rule = ValidationRule.email("userEmail");

        assertEquals("userEmail", rule.field());
        assertEquals("email", rule.ruleType());
        assertTrue(rule.required());
    }

    @Test
    @DisplayName("ValidationRule.number creates number rule")
    void validationRule_number() {
        ValidationRule rule = ValidationRule.number("amount");

        assertEquals("amount", rule.field());
        assertEquals("number", rule.ruleType());
        assertTrue(rule.required());
    }

    @Test
    @DisplayName("ValidationRule.optionalString creates non-required rule")
    void validationRule_optionalString() {
        ValidationRule rule = ValidationRule.optionalString("notes");

        assertEquals("notes", rule.field());
        assertFalse(rule.required());
    }

    // =========================================================================
    // SDKExecuteWorkflowRequest builder tests
    // =========================================================================

    @Test
    @DisplayName("SDKExecuteWorkflowRequest builder sets all fields")
    void executeWorkflowRequest_builder() {
        Map<String, Object> data = Map.of("key", "value");
        List<ValidationRule> rules = List.of(ValidationRule.required("key"));

        SDKExecuteWorkflowRequest request = SDKExecuteWorkflowRequest.builder()
                .inputData(data)
                .callbackUrl("https://example.com/hook")
                .timeoutMs(10_000L)
                .waitForCompletion(true)
                .metadata(Map.of("env", "test"))
                .validationRules(rules)
                .enableValidation(true)
                .enableSanitization(false)
                .priority(7)
                .build();

        assertEquals(data, request.getInputData());
        assertEquals("https://example.com/hook", request.getCallbackUrl());
        assertEquals(10_000L, request.getTimeoutMs());
        assertTrue(request.isWaitForCompletion());
        assertEquals("test", request.getMetadata().get("env"));
        assertEquals(1, request.getValidationRules().size());
        assertTrue(request.isEnableValidation());
        assertFalse(request.isEnableSanitization());
        assertEquals(7, request.getPriority());
    }

    @Test
    @DisplayName("SDKExecuteWorkflowRequest defaults are sensible")
    void executeWorkflowRequest_defaults() {
        SDKExecuteWorkflowRequest request = SDKExecuteWorkflowRequest.builder().build();

        assertNotNull(request.getInputData());
        assertTrue(request.getInputData().isEmpty());
        assertNull(request.getCallbackUrl());
        assertFalse(request.isWaitForCompletion());
        assertTrue(request.isEnableValidation());
        assertTrue(request.isEnableSanitization());
        assertEquals(5, request.getPriority());
        assertNotNull(request.getValidationRules());
    }

    // =========================================================================
    // ExecutionStatus record tests
    // =========================================================================

    @Test
    @DisplayName("ExecutionStatus.isCompleted returns true for completed status")
    void executionStatus_isCompleted() {
        ExecutionStatus completed = completedStatus("run-123");
        assertTrue(completed.isCompleted());
        assertFalse(completed.isFailed());
        assertTrue(completed.isTerminal());
    }

    @Test
    @DisplayName("ExecutionStatus.isFailed returns true for failed status")
    void executionStatus_isFailed() {
        ExecutionStatus failed = new ExecutionStatus(
                "run-456", "failed", 30, "step-2",
                true, "Step raised exception", Instant.now(), Instant.now());

        assertTrue(failed.isFailed());
        assertFalse(failed.isCompleted());
        assertTrue(failed.isTerminal());
        assertEquals("Step raised exception", failed.errorMessage());
    }

    @Test
    @DisplayName("ExecutionStatus pending is not terminal")
    void executionStatus_pendingNotTerminal() {
        ExecutionStatus pending = pendingStatus("run-789");
        assertFalse(pending.isTerminal());
        assertFalse(pending.isCompleted());
    }

    // =========================================================================
    // BatchExecuteWorkflowsRequest builder tests
    // =========================================================================

    @Test
    @DisplayName("BatchExecuteWorkflowsRequest builder adds items correctly")
    void batchRequest_builder() {
        SDKExecuteWorkflowRequest req1 = SDKExecuteWorkflowRequest.builder()
                .inputData(Map.of("a", 1)).build();
        SDKExecuteWorkflowRequest req2 = SDKExecuteWorkflowRequest.builder()
                .inputData(Map.of("b", 2)).build();

        BatchExecuteWorkflowsRequest batch = BatchExecuteWorkflowsRequest.builder()
                .addExecution(new BatchExecutionItem("workflow-a", req1, 5))
                .addExecution(new BatchExecutionItem("workflow-b", req2, 3))
                .parallel(true)
                .maxConcurrent(4)
                .stopOnFirstFailure(true)
                .build();

        assertEquals(2, batch.getExecutions().size());
        assertTrue(batch.isParallel());
        assertEquals(4, batch.getMaxConcurrent());
        assertTrue(batch.isStopOnFirstFailure());
        assertEquals("workflow-a", batch.getExecutions().get(0).workflowId());
        assertEquals("workflow-b", batch.getExecutions().get(1).workflowId());
    }

    @Test
    @DisplayName("BatchExecutionItem enforces non-null workflowId")
    void batchExecutionItem_nullWorkflowIdThrows() {
        assertThrows(NullPointerException.class,
                () -> new BatchExecutionItem(null, null, 5));
    }

    // =========================================================================
    // WebhookConfiguration builder tests
    // =========================================================================

    @Test
    @DisplayName("WebhookConfiguration builder sets required fields")
    void webhookConfig_builder() {
        WebhookConfiguration config = WebhookConfiguration.builder()
                .url("https://example.com/hooks")
                .events(List.of("workflow_completed", "workflow_failed"))
                .secret("s3cr3t")
                .retryCount(5)
                .enabled(true)
                .build();

        assertEquals("https://example.com/hooks", config.getUrl());
        assertEquals(2, config.getEvents().size());
        assertEquals("s3cr3t", config.getSecret());
        assertEquals(5, config.getRetryCount());
        assertTrue(config.isEnabled());
    }

    @Test
    @DisplayName("WebhookConfiguration build() throws when url is missing")
    void webhookConfig_missingUrlThrows() {
        assertThrows(NullPointerException.class,
                () -> WebhookConfiguration.builder()
                        .events(List.of("workflow_completed"))
                        .build());
    }

    // =========================================================================
    // Error class tests
    // =========================================================================

    @Test
    @DisplayName("WorkflowSDKException.isRetryable returns true for transient errors")
    void sdkException_isRetryable_transient() {
        assertTrue(new NetworkException("network failure").isRetryable());
        assertTrue(new TimeoutException("timed out").isRetryable());
        assertTrue(new RateLimitException("limited", Duration.ofSeconds(5)).isRetryable());
    }

    @Test
    @DisplayName("WorkflowSDKException.isRetryable returns false for client errors")
    void sdkException_isRetryable_clientErrors() {
        assertFalse(new WorkflowNotFoundException("wf-1").isRetryable());
        assertFalse(new ExecutionNotFoundException("run-1").isRetryable());
        assertFalse(new AuthenticationException("bad token").isRetryable());
        assertFalse(new ValidationException("invalid", List.of()).isRetryable());
    }

    @Test
    @DisplayName("RateLimitException exposes retryAfter")
    void rateLimitException_retryAfter() {
        Duration retryAfter = Duration.ofSeconds(30);
        RateLimitException ex = new RateLimitException("too many requests", retryAfter);

        assertEquals(ErrorCode.RATE_LIMITED, ex.getCode());
        assertEquals(429, ex.getHttpStatus());
        assertEquals(retryAfter, ex.getRetryAfter());
    }

    @Test
    @DisplayName("ValidationException exposes errors list")
    void validationException_errors() {
        ValidationError ve = new ValidationError("email", "email", "Invalid email format", "bad@");
        ValidationException ex = new ValidationException("Validation failed", List.of(ve));

        assertEquals(1, ex.getValidationErrors().size());
        assertEquals("email", ex.getValidationErrors().get(0).field());
        assertEquals(ErrorCode.VALIDATION_FAILED, ex.getCode());
    }

    @Test
    @DisplayName("WorkflowNotFoundException stores workflowId")
    void workflowNotFoundException_storesId() {
        WorkflowNotFoundException ex = new WorkflowNotFoundException("antifraud");
        assertEquals("antifraud", ex.getWorkflowId());
        assertEquals(404, ex.getHttpStatus());
        assertEquals(ErrorCode.WORKFLOW_NOT_FOUND, ex.getCode());
    }

    @Test
    @DisplayName("ExecutionNotFoundException stores runId")
    void executionNotFoundException_storesRunId() {
        ExecutionNotFoundException ex = new ExecutionNotFoundException("run-abc");
        assertEquals("run-abc", ex.getRunId());
        assertEquals(404, ex.getHttpStatus());
        assertEquals(ErrorCode.EXECUTION_NOT_FOUND, ex.getCode());
    }

    // =========================================================================
    // Model immutability tests
    // =========================================================================

    @Test
    @DisplayName("ExecutionDetails collections are unmodifiable")
    void executionDetails_collectionsImmutable() {
        ExecutionDetails details = new ExecutionDetails(
                "run-1", completedStatus("run-1"),
                List.of(), List.of(), 1000L, Map.of());

        assertThrows(UnsupportedOperationException.class,
                () -> details.steps().add(
                        new ExecutionStep("s1", "step", "completed", 100L, null, Map.of())));
        assertThrows(UnsupportedOperationException.class,
                () -> details.primitivesUsed().add("extra"));
    }

    @Test
    @DisplayName("ValidationResult.passed() returns valid result with empty collections")
    void validationResult_passed() {
        ValidationResult result = ValidationResult.passed();

        assertTrue(result.valid());
        assertTrue(result.errors().isEmpty());
        assertTrue(result.warnings().isEmpty());
        assertTrue(result.sanitizedData().isEmpty());
    }

    // =========================================================================
    // SDKConfig environment factory test
    // =========================================================================

    @Test
    @DisplayName("SDKConfig.fromEnvironment() returns a config with sensible defaults")
    void sdkConfig_fromEnvironment_defaults() {
        // Environment variables are not set in a unit test context, so we expect all defaults
        SDKConfig config = SDKConfig.fromEnvironment();
        assertNotNull(config);
        assertEquals("http://localhost:8080", config.workflowApiEndpoint());
        assertNull(config.authToken());
    }

    // =========================================================================
    // Text block usage (Java 21 feature verification)
    // =========================================================================

    @Test
    @DisplayName("Text blocks can be used for JSON test fixtures")
    void textBlocks_jsonFixture() {
        String json = """
                {
                  "run_id": "run-001",
                  "status": "completed",
                  "progress": 100,
                  "is_terminal": true
                }
                """;

        assertNotNull(json);
        assertTrue(json.contains("run-001"));
        assertTrue(json.contains("completed"));
    }

    // =========================================================================
    // Switch expression (Java 21 feature verification)
    // =========================================================================

    @Test
    @DisplayName("ErrorCode switch expression maps correctly")
    void errorCode_switchExpression() {
        ErrorCode code = ErrorCode.NETWORK_ERROR;
        boolean retryable = switch (code) {
            case NETWORK_ERROR, TIMEOUT, RATE_LIMITED, RETRY_EXHAUSTED -> true;
            default -> false;
        };
        assertTrue(retryable);
    }

    // =========================================================================
    // Pattern matching (Java 21 feature verification)
    // =========================================================================

    @Test
    @DisplayName("Pattern matching instanceof works for SDK exceptions")
    void patternMatching_instanceof() {
        WorkflowSDKException ex = new WorkflowNotFoundException("wf-x");

        String result;
        if (ex instanceof WorkflowNotFoundException notFound) {
            result = "not-found:" + notFound.getWorkflowId();
        } else if (ex instanceof ValidationException validation) {
            result = "validation:" + validation.getValidationErrors().size();
        } else {
            result = "other";
        }

        assertEquals("not-found:wf-x", result);
    }
}
