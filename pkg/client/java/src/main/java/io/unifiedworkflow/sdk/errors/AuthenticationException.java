package io.unifiedworkflow.sdk.errors;

/**
 * Thrown when the server rejects the client's authentication credentials.
 *
 * <p>This typically corresponds to HTTP 401 responses. Verify that the
 * {@code authToken} in {@link io.unifiedworkflow.sdk.SDKConfig} is correct
 * and has not expired.</p>
 */
public class AuthenticationException extends WorkflowSDKException {

    private static final long serialVersionUID = 1L;

    /**
     * Constructs a new authentication exception with the given message.
     *
     * @param message human-readable description of the authentication failure
     */
    public AuthenticationException(String message) {
        super(ErrorCode.AUTHENTICATION_FAILED, message, 401);
    }
}
