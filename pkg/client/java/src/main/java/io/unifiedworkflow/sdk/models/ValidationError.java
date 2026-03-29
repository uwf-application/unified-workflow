package io.unifiedworkflow.sdk.models;

import com.fasterxml.jackson.annotation.JsonIgnoreProperties;
import com.fasterxml.jackson.annotation.JsonProperty;

/**
 * A single validation error produced when input data fails a {@link ValidationRule}.
 *
 * @param field     the field path that failed validation (e.g. {@code "email"} or {@code "address.zip"})
 * @param ruleType  the rule type that was violated (e.g. {@code "email"}, {@code "required"})
 * @param message   human-readable description of the violation
 * @param value     the actual value that was supplied for the field; may be null
 */
@JsonIgnoreProperties(ignoreUnknown = true)
public record ValidationError(
        @JsonProperty("field") String field,
        @JsonProperty("rule_type") String ruleType,
        @JsonProperty("message") String message,
        @JsonProperty("value") Object value
) {

    /**
     * Compact constructor that enforces non-null field and message.
     */
    public ValidationError {
        if (field == null || field.isBlank()) {
            throw new IllegalArgumentException("field must not be blank");
        }
        if (message == null || message.isBlank()) {
            throw new IllegalArgumentException("message must not be blank");
        }
    }
}
