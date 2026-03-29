package io.unifiedworkflow.sdk.models;

import com.fasterxml.jackson.annotation.JsonIgnoreProperties;
import com.fasterxml.jackson.annotation.JsonProperty;

import java.time.Instant;
import java.util.Collections;
import java.util.List;
import java.util.Map;

/**
 * Represents the authenticated user session associated with a workflow execution request.
 *
 * @param userId      unique identifier of the authenticated user; may be null for anonymous requests
 * @param sessionId   session or token identifier; may be null
 * @param roles       list of roles assigned to the user; never null, may be empty
 * @param permissions list of granular permissions granted to the user; never null, may be empty
 * @param authMethod  authentication method used (e.g. {@code "bearer_token"}, {@code "api_key"});
 *                    may be null
 * @param expiresAt   session expiry time; may be null for non-expiring sessions
 * @param attributes  additional arbitrary session attributes; never null, may be empty
 */
@JsonIgnoreProperties(ignoreUnknown = true)
public record SessionContext(
        @JsonProperty("user_id") String userId,
        @JsonProperty("session_id") String sessionId,
        @JsonProperty("roles") List<String> roles,
        @JsonProperty("permissions") List<String> permissions,
        @JsonProperty("auth_method") String authMethod,
        @JsonProperty("expires_at") Instant expiresAt,
        @JsonProperty("attributes") Map<String, String> attributes
) {

    /**
     * Compact constructor that ensures collections are never null.
     */
    public SessionContext {
        roles = roles != null ? Collections.unmodifiableList(roles) : Collections.emptyList();
        permissions = permissions != null
                ? Collections.unmodifiableList(permissions)
                : Collections.emptyList();
        attributes = attributes != null
                ? Collections.unmodifiableMap(attributes)
                : Collections.emptyMap();
    }
}
