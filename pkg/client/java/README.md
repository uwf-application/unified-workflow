# Unified Workflow Java SDK

Java 21 client SDK for the [Unified Workflow System](../../README.md) — a high-performance workflow orchestration platform for banking environments.

## Requirements

- Java 21 or later
- Maven 3.8+ (or Gradle 8+)

## Installation

### Maven

```xml
<dependency>
    <groupId>io.unifiedworkflow</groupId>
    <artifactId>unified-workflow-sdk</artifactId>
    <version>1.2.0</version>
</dependency>
```

### Gradle

```groovy
implementation 'io.unifiedworkflow:unified-workflow-sdk:1.2.0'
```

---

## Quick Start

```java
// 1. Create the client
try (UnifiedWorkflowSDK sdk = UnifiedWorkflowSDK.create("http://workflow-api:8080")) {

    // 2. Execute a workflow
    SDKExecuteWorkflowResponse response = sdk.executeWorkflow(
        "antifraud",
        Map.of("transactionId", "txn-001", "amount", 1500.0, "currency", "USD")
    );

    // 3. Wait for completion
    ExecutionStatus status = sdk.waitForCompletion(response.runId(), Duration.ofSeconds(30));
    System.out.println("Result: " + status.status());
}
```

---

## Configuration

Create an `SDKConfig` using the builder and pass it to `UnifiedWorkflowSDK.create()`:

```java
SDKConfig config = SDKConfig.builder()
    .workflowApiEndpoint("https://workflow.example.com")
    .authToken(System.getenv("WORKFLOW_API_TOKEN"))
    .authType(SDKConfig.AuthType.BEARER_TOKEN)
    .timeout(Duration.ofSeconds(60))
    .maxRetries(3)
    .retryDelay(Duration.ofMillis(500))
    .enableRequestLogging(true)
    .build();

UnifiedWorkflowSDK sdk = UnifiedWorkflowSDK.create(config);
```

### Configuration from environment variables

```java
SDKConfig config = SDKConfig.fromEnvironment();
```

| Variable | Description | Default |
|---|---|---|
| `WORKFLOW_API_ENDPOINT` | Base URL of the workflow API | `http://localhost:8080` |
| `WORKFLOW_AUTH_TOKEN` | Bearer token or API key | (none) |
| `WORKFLOW_AUTH_TYPE` | One of `BEARER_TOKEN`, `API_KEY`, `BASIC_AUTH`, `NONE` | `NONE` |
| `WORKFLOW_TIMEOUT_MS` | Per-request timeout in milliseconds | `30000` |
| `WORKFLOW_MAX_RETRIES` | Retry attempts for transient errors | `3` |
| `WORKFLOW_ENABLE_REQUEST_LOGGING` | Log requests at DEBUG level | `false` |

### Configuration from YAML file

```java
SDKConfig config = SDKConfig.fromYaml(Path.of("/etc/workflow/sdk.yaml"));
```

Example YAML:

```yaml
workflowApiEndpoint: "http://workflow-api:8080"
authToken: "my-bearer-token"
authType: "BEARER_TOKEN"
timeoutMs: 30000
maxRetries: 3
retryDelayMs: 1000
enableValidation: true
enableSanitization: true
strictValidation: false
pollIntervalMs: 2000
defaultPriority: 5
enableCircuitBreaker: true
circuitBreakerThreshold: 5
circuitBreakerTimeoutMs: 60000
enableRequestLogging: false
```

### Full configuration reference

| Option | Type | Default | Description |
|---|---|---|---|
| `workflowApiEndpoint` | `String` | `http://localhost:8080` | Base URL of the workflow API |
| `timeout` | `Duration` | `PT30S` | Per-request HTTP timeout |
| `maxRetries` | `int` | `3` | Retry attempts for retryable errors |
| `retryDelay` | `Duration` | `PT1S` | Base retry delay (exponential back-off) |
| `authToken` | `String` | `null` | Bearer token or API key value |
| `authType` | `AuthType` | `NONE` | Authentication scheme |
| `enableValidation` | `boolean` | `true` | Enable server-side input validation |
| `enableSanitization` | `boolean` | `true` | Enable server-side input sanitization |
| `strictValidation` | `boolean` | `false` | Treat warnings as errors |
| `asyncExecution` | `boolean` | `true` | Submit executions asynchronously |
| `pollIntervalMs` | `long` | `2000` | Polling interval for `waitForCompletion` |
| `defaultPriority` | `int` | `5` | Default execution priority [1–10] |
| `enableCircuitBreaker` | `boolean` | `true` | Enable client-side circuit breaker |
| `circuitBreakerThreshold` | `int` | `5` | Failure count that trips the breaker |
| `circuitBreakerTimeout` | `Duration` | `PT60S` | Open duration before recovery attempt |
| `enableRequestLogging` | `boolean` | `false` | Log requests/responses at DEBUG level |

---

## API Reference

### Workflow Execution

| Method | Description |
|---|---|
| `executeWorkflow(workflowId, inputData)` | Execute a workflow with a plain data map |
| `executeWorkflow(workflowId, request)` | Execute a workflow with a full request object |
| `executeWorkflowAsync(workflowId, inputData)` | Async version, returns `CompletableFuture` |
| `executeWorkflowAsync(workflowId, request)` | Async version with full request object |

### Execution Status and Control

| Method | Description |
|---|---|
| `getExecutionStatus(runId)` | Get current status snapshot |
| `getExecutionDetails(runId)` | Get full execution details with per-step results |
| `waitForCompletion(runId, maxWait)` | Poll until terminal state or timeout |
| `cancelExecution(runId)` | Request cancellation |
| `pauseExecution(runId)` | Pause a running execution |
| `resumeExecution(runId)` | Resume a paused execution |
| `getExecutionStatusAsync(runId)` | Async status check |
| `getExecutionDetailsAsync(runId)` | Async details retrieval |

### Workflow Registry

| Method | Description |
|---|---|
| `listWorkflows()` | List all registered workflow definitions |
| `getWorkflow(workflowId)` | Get a single workflow definition |

### Batch Execution

| Method | Description |
|---|---|
| `batchExecuteWorkflows(request)` | Execute multiple workflows in one call |

### Validation

| Method | Description |
|---|---|
| `validate(data, rules)` | Validate data against rules without executing |

### Webhooks

| Method | Description |
|---|---|
| `registerWebhook(config)` | Register a webhook endpoint |
| `unregisterWebhook(webhookId)` | Remove a registered webhook |

### Utility

| Method | Description |
|---|---|
| `ping()` | Check server health |
| `getSdkVersion()` | Returns the SDK version string |
| `close()` | Release internal resources |

---

## Async Usage with CompletableFuture

```java
UnifiedWorkflowSDK sdk = UnifiedWorkflowSDK.create(config);

sdk.executeWorkflowAsync("antifraud", Map.of("transactionId", "txn-002", "amount", 500.0))
    .thenCompose(response -> sdk.getExecutionStatusAsync(response.runId()))
    .thenAccept(status -> System.out.println("Status: " + status.status()))
    .exceptionally(ex -> {
        System.err.println("Failed: " + ex.getMessage());
        return null;
    });
```

---

## Error Handling

All SDK exceptions extend `WorkflowSDKException` (a `RuntimeException`). Catch the base type
for a catch-all handler, or catch specific subtypes for fine-grained handling:

```java
try {
    SDKExecuteWorkflowResponse response = sdk.executeWorkflow("antifraud", data);
} catch (WorkflowNotFoundException e) {
    // The workflow is not registered — check the workflow ID
    log.error("Workflow not found: {}", e.getWorkflowId());
} catch (ValidationException e) {
    // Input data failed validation rules
    e.getValidationErrors().forEach(err ->
        log.warn("Validation: field={} message={}", err.field(), err.message()));
} catch (AuthenticationException e) {
    // Token is invalid or expired
    log.error("Authentication failed: {}", e.getMessage());
} catch (RateLimitException e) {
    // Server is rate-limiting this client
    log.warn("Rate limited; retry after {}", e.getRetryAfter());
    Thread.sleep(e.getRetryAfter().toMillis());
} catch (WorkflowSDKException e) {
    // Any other SDK error
    if (e.isRetryable()) {
        // safe to retry
    }
    log.error("SDK error [{}]: {}", e.getCode(), e.getMessage());
}
```

### Exception hierarchy

```
WorkflowSDKException (base)
├── ValidationException         — input failed validation rules
├── WorkflowNotFoundException   — workflow ID not registered
├── ExecutionNotFoundException  — run ID not found
├── AuthenticationException     — 401 Unauthorized
├── NetworkException            — I/O / connectivity failure
├── TimeoutException            — request timed out
└── RateLimitException          — 429 Too Many Requests
```

---

## Validation

Validate input data before or during execution:

```java
List<ValidationRule> rules = List.of(
    ValidationRule.required("transactionId"),
    ValidationRule.email("contactEmail"),
    ValidationRule.number("amount"),
    new ValidationRule("amount", "number", true, null, null, null, 0.01, 1_000_000.0, List.of())
);

ValidationResult result = sdk.validate(data, rules);
if (!result.valid()) {
    result.errors().forEach(err ->
        System.err.printf("Field '%s': %s%n", err.field(), err.message()));
}

// Or include validation in the execution request:
SDKExecuteWorkflowRequest request = SDKExecuteWorkflowRequest.builder()
    .inputData(data)
    .validationRules(rules)
    .enableValidation(true)
    .enableSanitization(true)
    .build();
sdk.executeWorkflow("antifraud", request);
```

---

## Batch Execution

```java
BatchExecuteWorkflowsRequest batch = BatchExecuteWorkflowsRequest.builder()
    .addExecution(new BatchExecutionItem(
        "antifraud",
        SDKExecuteWorkflowRequest.builder()
            .inputData(Map.of("transactionId", "txn-A", "amount", 100.0))
            .build(),
        5
    ))
    .addExecution(new BatchExecutionItem(
        "kyc-check",
        SDKExecuteWorkflowRequest.builder()
            .inputData(Map.of("customerId", "cust-001"))
            .build(),
        7
    ))
    .parallel(true)
    .maxConcurrent(4)
    .stopOnFirstFailure(false)
    .build();

BatchExecuteWorkflowsResponse response = sdk.batchExecuteWorkflows(batch);
System.out.printf("Batch %s: %d/%d succeeded%n",
    response.batchId(), response.successful(), response.total());

response.executions().forEach(result -> {
    if (result.success()) {
        System.out.println("  OK  " + result.workflowId() + " -> " + result.runId());
    } else {
        System.out.println("  ERR " + result.workflowId() + " -> " + result.error());
    }
});
```

---

## Polling / waitForCompletion

```java
// Submit asynchronously, then block with a timeout
SDKExecuteWorkflowResponse submitted = sdk.executeWorkflow(
    "antifraud",
    Map.of("transactionId", "txn-003", "amount", 750.0)
);

ExecutionStatus terminal = sdk.waitForCompletion(
    submitted.runId(),
    Duration.ofSeconds(60)
);

if (terminal.isCompleted()) {
    System.out.println("Workflow completed successfully");
} else if (terminal.isFailed()) {
    System.err.println("Workflow failed: " + terminal.errorMessage());
}
```

The poll interval is controlled by `SDKConfig.pollIntervalMs()` (default: 2 000 ms).

---

## Webhook Registration

```java
WebhookConfiguration webhook = WebhookConfiguration.builder()
    .url("https://your-service.example.com/workflow-events")
    .events(List.of("workflow_completed", "workflow_failed"))
    .secret("hmac-signing-secret")
    .retryCount(3)
    .enabled(true)
    .build();

WebhookConfiguration registered = sdk.registerWebhook(webhook);
System.out.println("Registered webhook ID: " + registered.getWebhookId());

// Later, to remove:
sdk.unregisterWebhook(registered.getWebhookId());
```

---

## Java 21 Features

This SDK makes use of the following Java 21 language features:

- **Records** — all immutable model classes (`ExecutionStatus`, `ExecutionStep`, `ValidationRule`, etc.)
- **Virtual threads** — `HttpClientWrapper` uses `Executors.newVirtualThreadPerTaskExecutor()` for
  scalable async I/O without blocking platform threads
- **Pattern matching** — `instanceof` patterns in exception handling code paths
- **Switch expressions** — used in circuit-breaker state machine and error-code classification
- **Text blocks** — used in tests for inline JSON fixtures
- **Sealed-class-friendly** — the exception hierarchy is designed for exhaustive `switch` matching

---

## Building from Source

```bash
cd pkg/client/java
mvn clean package
mvn test
```
