package io.unifiedworkflow.sdk.errors;

/**
 * Enumeration of all error codes produced by the Unified Workflow SDK.
 *
 * <p>Use these codes to programmatically distinguish error conditions without
 * relying on message-string parsing.</p>
 */
public enum ErrorCode {

    /** SDK or client configuration is invalid. */
    INVALID_CONFIG,

    /** Request payload failed validation rules. */
    VALIDATION_FAILED,

    /** The requested workflow definition was not found. */
    WORKFLOW_NOT_FOUND,

    /** The requested workflow execution run was not found. */
    EXECUTION_NOT_FOUND,

    /** Authentication credentials are missing or invalid. */
    AUTHENTICATION_FAILED,

    /** A network-level error occurred while communicating with the server. */
    NETWORK_ERROR,

    /** The request timed out before a response was received. */
    TIMEOUT,

    /** The server has applied rate-limiting to this client. */
    RATE_LIMITED,

    /** One or more executions in a batch request failed. */
    BATCH_EXECUTION_FAILED,

    /** A webhook registration or delivery operation failed. */
    WEBHOOK_ERROR,

    /** The circuit breaker is open; requests are being rejected to protect the server. */
    CIRCUIT_BREAKER_OPEN,

    /** All retry attempts were exhausted without a successful response. */
    RETRY_EXHAUSTED,

    /** An unexpected or unclassified error occurred. */
    UNKNOWN
}
