package io.unifiedworkflow.sdk.models;

import com.fasterxml.jackson.annotation.JsonIgnoreProperties;
import com.fasterxml.jackson.annotation.JsonProperty;

import java.util.Collections;
import java.util.List;
import java.util.Map;

/**
 * Security metadata associated with a workflow execution request.
 *
 * <p>Used by the workflow engine for audit logging and policy enforcement.</p>
 *
 * @param authenticated         {@code true} if the request carries valid credentials
 * @param authenticationMethod  authentication scheme used (e.g. {@code "bearer_token"},
 *                              {@code "api_key"}); may be null
 * @param scopes                OAuth/JWT scopes granted to the caller; never null, may be empty
 * @param claims                decoded JWT claims or equivalent key-value assertions;
 *                              never null, may be empty
 * @param ipAddress             originating IP address of the caller; may be null
 * @param userAgent             value of the {@code User-Agent} header; may be null
 * @param geoLocation           ISO 3166-1 alpha-2 country code or similar geo tag; may be null
 */
@JsonIgnoreProperties(ignoreUnknown = true)
public record SecurityContext(
        @JsonProperty("authenticated") boolean authenticated,
        @JsonProperty("authentication_method") String authenticationMethod,
        @JsonProperty("scopes") List<String> scopes,
        @JsonProperty("claims") Map<String, String> claims,
        @JsonProperty("ip_address") String ipAddress,
        @JsonProperty("user_agent") String userAgent,
        @JsonProperty("geo_location") String geoLocation
) {

    /**
     * Compact constructor that ensures collections are never null.
     */
    public SecurityContext {
        scopes = scopes != null ? Collections.unmodifiableList(scopes) : Collections.emptyList();
        claims = claims != null ? Collections.unmodifiableMap(claims) : Collections.emptyMap();
    }
}
