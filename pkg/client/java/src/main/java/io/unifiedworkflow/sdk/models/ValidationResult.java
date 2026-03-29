package io.unifiedworkflow.sdk.models;

import com.fasterxml.jackson.annotation.JsonIgnoreProperties;
import com.fasterxml.jackson.annotation.JsonProperty;

import java.util.Collections;
import java.util.List;
import java.util.Map;

/**
 * The outcome of running a set of {@link ValidationRule} rules against input data.
 *
 * @param valid          {@code true} if all rules passed and no errors were reported
 * @param errors         list of field-level validation failures; never null, empty when valid
 * @param warnings       list of non-fatal advisory messages; never null, may be empty
 * @param sanitizedData  the input data after sanitization (e.g. trimmed strings, normalised values);
 *                       never null, may be empty if sanitization was not requested
 */
@JsonIgnoreProperties(ignoreUnknown = true)
public record ValidationResult(
        @JsonProperty("valid") boolean valid,
        @JsonProperty("errors") List<ValidationError> errors,
        @JsonProperty("warnings") List<String> warnings,
        @JsonProperty("sanitized_data") Map<String, Object> sanitizedData
) {

    /**
     * Compact constructor that guarantees collections are never null.
     */
    public ValidationResult {
        errors = errors != null ? Collections.unmodifiableList(errors) : Collections.emptyList();
        warnings = warnings != null ? Collections.unmodifiableList(warnings) : Collections.emptyList();
        sanitizedData = sanitizedData != null
                ? Collections.unmodifiableMap(sanitizedData)
                : Collections.emptyMap();
    }

    /**
     * Convenience factory for a passing result with no warnings and no sanitized data.
     *
     * @return a {@code ValidationResult} representing a clean pass
     */
    public static ValidationResult passed() {
        return new ValidationResult(true, Collections.emptyList(),
                Collections.emptyList(), Collections.emptyMap());
    }
}
