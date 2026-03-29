package io.unifiedworkflow.sdk.errors;

/**
 * Thrown when a request to the workflow server times out.
 *
 * <p>The timeout duration is controlled by {@link io.unifiedworkflow.sdk.SDKConfig#timeout()}.
 * {@link #isRetryable()} returns {@code true} for this exception.</p>
 */
public class TimeoutException extends WorkflowSDKException {

    private static final long serialVersionUID = 1L;

    /**
     * Constructs a new timeout exception with the given message.
     *
     * @param message human-readable description of the timeout condition
     */
    public TimeoutException(String message) {
        super(ErrorCode.TIMEOUT, message);
    }

    /**
     * Constructs a new timeout exception wrapping an underlying cause.
     *
     * @param message human-readable description of the timeout condition
     * @param cause   the underlying exception (e.g. {@link java.net.http.HttpTimeoutException})
     */
    public TimeoutException(String message, Throwable cause) {
        super(ErrorCode.TIMEOUT, message, cause);
    }
}
