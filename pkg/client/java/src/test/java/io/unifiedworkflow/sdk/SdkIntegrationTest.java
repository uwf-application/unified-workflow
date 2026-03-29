package io.unifiedworkflow.sdk;

import io.unifiedworkflow.sdk.errors.WorkflowSDKException;
import io.unifiedworkflow.sdk.models.*;
import org.junit.jupiter.api.*;

import java.util.List;
import java.util.Map;

import static org.junit.jupiter.api.Assertions.*;

/**
 * Live integration tests against locally running services (port 8082).
 *
 * Run with:  mvn test -Dtest=SdkIntegrationTest
 *
 * Services expected:
 *   workflow-api  -> http://localhost:8082
 */
@TestMethodOrder(MethodOrderer.OrderAnnotation.class)
class SdkIntegrationTest {

    private static final String API_BASE     = "http://localhost:8082";
    // Echo Workflow — always registered at startup
    private static final String ECHO_WORKFLOW_NAME = "Echo Workflow";

    private static UnifiedWorkflowSDK sdk;
    private static String echoWorkflowId;
    private static String lastRunId;

    @BeforeAll
    static void setUp() {
        sdk = UnifiedWorkflowSDK.create(
                SDKConfig.builder()
                        .workflowApiEndpoint(API_BASE)
                        .enableRequestLogging(true)
                        .maxRetries(0)           // fail fast in tests
                        .build()
        );
        System.out.println("[SDK] version=" + sdk.getSdkVersion()
                + "  endpoint=" + sdk.getConfig().workflowApiEndpoint());
    }

    @AfterAll
    static void tearDown() {
        if (sdk != null) sdk.close();
    }

    // -------------------------------------------------------------------------
    // 1. Health check
    // -------------------------------------------------------------------------

    @Test
    @Order(1)
    @DisplayName("ping() — service is reachable and healthy")
    void testPing() {
        boolean alive = sdk.ping();
        assertTrue(alive, "ping() must return true for a healthy service");
        System.out.println("[PASS] ping()");
    }

    // -------------------------------------------------------------------------
    // 2. Workflow listing
    // -------------------------------------------------------------------------

    @Test
    @Order(2)
    @DisplayName("listWorkflows() — returns at least one registered workflow")
    void testListWorkflows() {
        List<WorkflowDefinition> workflows = sdk.listWorkflows();

        assertNotNull(workflows, "listWorkflows() must never return null");
        assertFalse(workflows.isEmpty(), "Expected at least one registered workflow");

        System.out.println("[PASS] listWorkflows() — count=" + workflows.size());
        workflows.forEach(w -> System.out.printf("       id=%-45s  name=%s%n", w.id(), w.name()));

        // Capture the Echo Workflow ID for later tests
        echoWorkflowId = workflows.stream()
                .filter(w -> ECHO_WORKFLOW_NAME.equals(w.name()))
                .map(WorkflowDefinition::id)
                .findFirst()
                .orElseGet(() -> workflows.get(0).id());

        System.out.println("[INFO] using workflowId=" + echoWorkflowId);
    }

    // -------------------------------------------------------------------------
    // 3. Get single workflow
    // -------------------------------------------------------------------------

    @Test
    @Order(3)
    @DisplayName("getWorkflow() — returns correct definition for a known ID")
    void testGetWorkflow() {
        assumeWorkflowIdKnown();

        WorkflowDefinition wf = sdk.getWorkflow(echoWorkflowId);

        assertNotNull(wf, "getWorkflow() must not return null");
        assertEquals(echoWorkflowId, wf.id(), "Returned workflow ID must match requested ID");
        assertNotNull(wf.name(),        "Workflow name must not be null");
        assertNotNull(wf.description(), "Workflow description must not be null");

        System.out.printf("[PASS] getWorkflow()  id=%s  name=%s%n", wf.id(), wf.name());
    }

    // -------------------------------------------------------------------------
    // 4. Execute workflow
    // -------------------------------------------------------------------------

    @Test
    @Order(4)
    @DisplayName("executeWorkflow() — returns a runId and pending/running status")
    void testExecuteWorkflow() {
        assumeWorkflowIdKnown();

        Map<String, Object> input = Map.of(
                "message",       "hello from Java SDK integration test",
                "test_run",      true,
                "timestamp",     System.currentTimeMillis()
        );

        SDKExecuteWorkflowResponse response = sdk.executeWorkflow(echoWorkflowId, input);

        assertNotNull(response, "executeWorkflow() must not return null");
        assertNotNull(response.runId(), "runId must not be null");
        assertFalse(response.runId().isBlank(), "runId must not be blank");

        lastRunId = response.runId();

        System.out.printf("[PASS] executeWorkflow()  runId=%s  status=%s  message=%s%n",
                response.runId(), response.status(), response.message());
    }

    // -------------------------------------------------------------------------
    // 5. Execute with full request builder
    // -------------------------------------------------------------------------

    @Test
    @Order(5)
    @DisplayName("executeWorkflow(request) — builder API produces a valid run")
    void testExecuteWorkflowWithRequest() {
        assumeWorkflowIdKnown();

        SDKExecuteWorkflowRequest request = SDKExecuteWorkflowRequest.builder()
                .inputData(Map.of("amount", 100.00, "currency", "KZT"))
                .enableValidation(false)
                .enableSanitization(false)
                .priority(8)
                .metadata(Map.of("source", "java-sdk-integration-test"))
                .build();

        SDKExecuteWorkflowResponse response = sdk.executeWorkflow(echoWorkflowId, request);

        assertNotNull(response);
        assertNotNull(response.runId());
        assertFalse(response.runId().isBlank());

        System.out.printf("[PASS] executeWorkflow(request)  runId=%s%n", response.runId());
    }

    // -------------------------------------------------------------------------
    // 6. Get execution status
    // -------------------------------------------------------------------------

    @Test
    @Order(6)
    @DisplayName("getExecutionStatus() — returns status for a previously submitted run")
    void testGetExecutionStatus() {
        assumeRunIdKnown();

        ExecutionStatus status = sdk.getExecutionStatus(lastRunId);

        assertNotNull(status, "getExecutionStatus() must not return null");
        assertEquals(lastRunId, status.runId(), "Returned runId must match the submitted run");
        assertNotNull(status.status(), "status field must not be null");

        System.out.printf("[PASS] getExecutionStatus()  runId=%s  status=%s  isTerminal=%s  progress=%d%%%n",
                status.runId(), status.status(), status.isTerminal(), status.progress());
    }

    // -------------------------------------------------------------------------
    // 7. Async execute
    // -------------------------------------------------------------------------

    @Test
    @Order(7)
    @DisplayName("executeWorkflowAsync() — CompletableFuture completes successfully")
    void testExecuteWorkflowAsync() throws Exception {
        assumeWorkflowIdKnown();

        SDKExecuteWorkflowResponse response = sdk
                .executeWorkflowAsync(echoWorkflowId, Map.of("async", true))
                .get();  // block for test simplicity

        assertNotNull(response);
        assertNotNull(response.runId());
        System.out.printf("[PASS] executeWorkflowAsync()  runId=%s%n", response.runId());
    }

    // -------------------------------------------------------------------------
    // 8. Async status
    // -------------------------------------------------------------------------

    @Test
    @Order(8)
    @DisplayName("getExecutionStatusAsync() — CompletableFuture completes with valid status")
    void testGetExecutionStatusAsync() throws Exception {
        assumeRunIdKnown();

        ExecutionStatus status = sdk.getExecutionStatusAsync(lastRunId).get();

        assertNotNull(status);
        assertEquals(lastRunId, status.runId());
        System.out.printf("[PASS] getExecutionStatusAsync()  runId=%s  status=%s%n",
                status.runId(), status.status());
    }

    // -------------------------------------------------------------------------
    // 9. Null-guard checks
    // -------------------------------------------------------------------------

    @Test
    @Order(9)
    @DisplayName("Null workflowId throws NullPointerException immediately")
    void testNullWorkflowIdThrows() {
        assertThrows(NullPointerException.class,
                () -> sdk.executeWorkflow(null, Map.of()),
                "Passing null workflowId must throw NullPointerException");
        System.out.println("[PASS] null workflowId guard");
    }

    @Test
    @Order(10)
    @DisplayName("Null runId throws NullPointerException immediately")
    void testNullRunIdThrows() {
        assertThrows(NullPointerException.class,
                () -> sdk.getExecutionStatus(null),
                "Passing null runId must throw NullPointerException");
        System.out.println("[PASS] null runId guard");
    }

    // -------------------------------------------------------------------------
    // 10. Unknown workflow — expect SDK exception
    // -------------------------------------------------------------------------

    @Test
    @Order(11)
    @DisplayName("executeWorkflow() with unknown workflowId throws WorkflowSDKException")
    void testUnknownWorkflowThrowsException() {
        WorkflowSDKException ex = assertThrows(WorkflowSDKException.class,
                () -> sdk.executeWorkflow("non-existent-workflow-id-xyz", Map.of("k", "v")),
                "Unknown workflow should raise WorkflowSDKException");

        System.out.printf("[PASS] unknown workflow exception  code=%s  status=%d  msg=%s%n",
                ex.getCode(), ex.getHttpStatus(), ex.getMessage());
    }

    // -------------------------------------------------------------------------
    // 11. Unknown run ID — API returns pending (not 404), verify graceful handling
    // -------------------------------------------------------------------------

    @Test
    @Order(12)
    @DisplayName("getExecutionStatus() with unknown runId returns a status (API does not 404)")
    void testUnknownRunIdReturnsPending() {
        // The workflow-api returns a synthetic pending record for unknown run IDs
        ExecutionStatus status = sdk.getExecutionStatus("run-does-not-exist-00000");
        assertNotNull(status, "status must not be null even for unknown runs");
        System.out.printf("[PASS] unknown runId returns status=%s  isTerminal=%s%n",
                status.status(), status.isTerminal());
    }

    // -------------------------------------------------------------------------
    // Helpers
    // -------------------------------------------------------------------------

    private static void assumeWorkflowIdKnown() {
        if (echoWorkflowId == null) {
            // Resolve on-demand so individual tests can run standalone
            List<WorkflowDefinition> wfs = sdk.listWorkflows();
            assertFalse(wfs.isEmpty(), "No workflows registered — cannot proceed");
            echoWorkflowId = wfs.stream()
                    .filter(w -> ECHO_WORKFLOW_NAME.equals(w.name()))
                    .map(WorkflowDefinition::id)
                    .findFirst()
                    .orElseGet(() -> wfs.get(0).id());
        }
    }

    private static void assumeRunIdKnown() {
        if (lastRunId == null) {
            assumeWorkflowIdKnown();
            SDKExecuteWorkflowResponse r = sdk.executeWorkflow(
                    echoWorkflowId, Map.of("seed", "status-test"));
            lastRunId = r.runId();
        }
    }
}
