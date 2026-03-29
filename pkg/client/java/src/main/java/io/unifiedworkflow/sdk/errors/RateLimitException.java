package io.unifiedworkflow.sdk.errors;

import java.time.Duration;

/**
 * Thrown when the server returns HTTP 429 (Too Many Requests), indicating that this
 * client has exceeded the allowed request rate.
 *
 * <p>Use {@link #getRetryAfter()} to determine how long to wait before the next attempt.</p>
 */
public class RateLimitException extends WorkflowSDKException {

    private static final long serialVersionUID = 1L;

    private final Duration retryAfter;

    /**
     * Constructs a new rate-limit exception.
     *
     * @param message    human-readable description
     * @param retryAfter suggested wait duration before retrying; must not be null; use
     *                   {@link Duration#ZERO} when the server did not supply a Retry-After header
     */
    public RateLimitException(String message, Duration retryAfter) {
        super(ErrorCode.RATE_LIMITED, message, 429);
        this.retryAfter = retryAfter != null ? retryAfter : Duration.ZERO;
    }

    /**
     * Returns the suggested duration to wait before retrying the request.
     *
     * <p>When the server did not supply a {@code Retry-After} header this method
     * returns {@link Duration#ZERO}.</p>
     *
     * @return the suggested retry delay, never null
     */
    public Duration getRetryAfter() {
        return retryAfter;
    }
}
