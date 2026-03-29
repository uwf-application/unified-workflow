package io.unifiedworkflow.sdk.errors;

/**
 * Thrown when a network-level error prevents communication with the workflow server.
 *
 * <p>This includes connection refused, DNS resolution failures, and unexpected
 * I/O errors. {@link #isRetryable()} always returns {@code true} for this exception.</p>
 */
public class NetworkException extends WorkflowSDKException {

    private static final long serialVersionUID = 1L;

    /**
     * Constructs a new network exception with the given message.
     *
     * @param message human-readable description of the network failure
     */
    public NetworkException(String message) {
        super(ErrorCode.NETWORK_ERROR, message);
    }

    /**
     * Constructs a new network exception wrapping an underlying cause.
     *
     * @param message human-readable description of the network failure
     * @param cause   the underlying I/O or connectivity exception
     */
    public NetworkException(String message, Throwable cause) {
        super(ErrorCode.NETWORK_ERROR, message, cause);
    }
}
