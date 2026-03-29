package io.unifiedworkflow.sdk.models;

import com.fasterxml.jackson.annotation.JsonIgnoreProperties;
import com.fasterxml.jackson.annotation.JsonProperty;

import java.util.Collections;
import java.util.List;

/**
 * Defines a validation constraint for a single field in a workflow request payload.
 *
 * <p>Use the static factory methods for common rule types:</p>
 * <pre>{@code
 * List<ValidationRule> rules = List.of(
 *     ValidationRule.required("transactionId"),
 *     ValidationRule.email("contactEmail"),
 *     ValidationRule.number("amount")
 * );
 * }</pre>
 *
 * @param field         the field path this rule applies to (e.g. {@code "email"})
 * @param ruleType      the type of validation: {@code "string"}, {@code "number"},
 *                      {@code "boolean"}, {@code "email"}, {@code "url"}, {@code "uuid"}
 * @param required      {@code true} if the field must be present and non-null
 * @param minLength     minimum string length; null means no minimum
 * @param maxLength     maximum string length; null means no maximum
 * @param pattern       regular-expression pattern the value must match; null means no pattern check
 * @param minValue      minimum numeric value (inclusive); null means no minimum
 * @param maxValue      maximum numeric value (inclusive); null means no maximum
 * @param allowedValues whitelist of allowed values; empty list means no whitelist restriction
 */
@JsonIgnoreProperties(ignoreUnknown = true)
public record ValidationRule(
        @JsonProperty("field") String field,
        @JsonProperty("rule_type") String ruleType,
        @JsonProperty("required") boolean required,
        @JsonProperty("min_length") Integer minLength,
        @JsonProperty("max_length") Integer maxLength,
        @JsonProperty("pattern") String pattern,
        @JsonProperty("min_value") Double minValue,
        @JsonProperty("max_value") Double maxValue,
        @JsonProperty("allowed_values") List<String> allowedValues
) {

    /**
     * Compact constructor that ensures allowedValues is never null.
     */
    public ValidationRule {
        allowedValues = allowedValues != null
                ? Collections.unmodifiableList(allowedValues)
                : Collections.emptyList();
    }

    // -------------------------------------------------------------------------
    // Static factory methods
    // -------------------------------------------------------------------------

    /**
     * Creates a required-field rule with no additional constraints.
     *
     * @param field the field path; must not be null or blank
     * @return a new rule that marks the field as required
     */
    public static ValidationRule required(String field) {
        return new ValidationRule(field, "string", true,
                null, null, null, null, null, Collections.emptyList());
    }

    /**
     * Creates a required string-type rule.
     *
     * @param field the field path; must not be null or blank
     * @return a new rule that validates the field as a string
     */
    public static ValidationRule string(String field) {
        return new ValidationRule(field, "string", true,
                null, null, null, null, null, Collections.emptyList());
    }

    /**
     * Creates a required email-format rule.
     *
     * @param field the field path; must not be null or blank
     * @return a new rule that validates the field as an email address
     */
    public static ValidationRule email(String field) {
        return new ValidationRule(field, "email", true,
                null, null, null, null, null, Collections.emptyList());
    }

    /**
     * Creates a required numeric rule.
     *
     * @param field the field path; must not be null or blank
     * @return a new rule that validates the field as a number
     */
    public static ValidationRule number(String field) {
        return new ValidationRule(field, "number", true,
                null, null, null, null, null, Collections.emptyList());
    }

    /**
     * Creates an optional string-type rule (field is not required to be present).
     *
     * @param field the field path; must not be null or blank
     * @return a new optional string rule
     */
    public static ValidationRule optionalString(String field) {
        return new ValidationRule(field, "string", false,
                null, null, null, null, null, Collections.emptyList());
    }
}
