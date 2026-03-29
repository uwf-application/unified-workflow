package io.unifiedworkflow.sdk.errors;

import java.util.Collections;
import java.util.Map;
import java.util.Objects;

/**
 * Base runtime exception for all errors produced by the Unified Workflow SDK.
 *
 * <p>All SDK-specific exceptions extend this class. Callers may catch this type
 * to handle any SDK failure, or catch a subclass for fine-grained handling.</p>
 *
 * <p>Example:</p>
 * <pre>{@code
 * try {
 *     sdk.executeWorkflow("antifraud", data);
 * } catch (WorkflowNotFoundException e) {
 *     // workflow not registered
 * } catch (WorkflowSDKException e) {
 *     if (e.isRetryable()) { /* schedule retry *\/ }
 * }
 * }</pre>
 */
public class WorkflowSDKException extends RuntimeException {

    private static final long serialVersionUID = 1L;

    private final ErrorCode code;
    private final int httpStatus;
    private final Map<String, Object> details;

    /**
     * Constructs a new exception with a code and message.
     *
     * @param code    the SDK error code, must not be null
     * @param message human-readable description
     */
    public WorkflowSDKException(ErrorCode code, String message) {
        super(message);
        this.code = Objects.requireNonNull(code, "code must not be null");
        this.httpStatus = 0;
        this.details = Collections.emptyMap();
    }

    /**
     * Constructs a new exception with a code, message, and HTTP status.
     *
     * @param code       the SDK error code, must not be null
     * @param message    human-readable description
     * @param httpStatus the HTTP status code from the server response, or 0 if not applicable
     */
    public WorkflowSDKException(ErrorCode code, String message, int httpStatus) {
        super(message);
        this.code = Objects.requireNonNull(code, "code must not be null");
        this.httpStatus = httpStatus;
        this.details = Collections.emptyMap();
    }

    /**
     * Constructs a new exception with a code, message, HTTP status, and additional detail map.
     *
     * @param code       the SDK error code, must not be null
     * @param message    human-readable description
     * @param httpStatus the HTTP status code from the server response, or 0 if not applicable
     * @param details    optional key-value detail map for structured debugging; may be null
     */
    public WorkflowSDKException(ErrorCode code, String message, int httpStatus, Map<String, Object> details) {
        super(message);
        this.code = Objects.requireNonNull(code, "code must not be null");
        this.httpStatus = httpStatus;
        this.details = details != null ? Collections.unmodifiableMap(details) : Collections.emptyMap();
    }

    /**
     * Constructs a new exception with a code, message, and cause.
     *
     * @param code    the SDK error code, must not be null
     * @param message human-readable description
     * @param cause   the underlying cause
     */
    public WorkflowSDKException(ErrorCode code, String message, Throwable cause) {
        super(message, cause);
        this.code = Objects.requireNonNull(code, "code must not be null");
        this.httpStatus = 0;
        this.details = Collections.emptyMap();
    }

    /**
     * Returns the SDK error code that classifies this failure.
     *
     * @return the error code, never null
     */
    public ErrorCode getCode() {
        return code;
    }

    /**
     * Returns the HTTP status code from the server response that triggered this exception.
     *
     * @return HTTP status code, or {@code 0} if the error did not originate from an HTTP response
     */
    public int getHttpStatus() {
        return httpStatus;
    }

    /**
     * Returns an immutable map of additional structured detail about the error.
     *
     * @return detail map, never null, may be empty
     */
    public Map<String, Object> getDetails() {
        return details;
    }

    /**
     * Returns {@code true} if the operation that caused this exception may succeed on retry.
     *
     * <p>Transient conditions such as network errors, timeouts, and rate-limiting are
     * considered retryable. Client errors (validation failures, authentication errors,
     * not-found responses) are not.</p>
     *
     * @return {@code true} if the caller should consider retrying the operation
     */
    public boolean isRetryable() {
        return switch (code) {
            case NETWORK_ERROR, TIMEOUT, RATE_LIMITED, RETRY_EXHAUSTED -> true;
            default -> false;
        };
    }
}
