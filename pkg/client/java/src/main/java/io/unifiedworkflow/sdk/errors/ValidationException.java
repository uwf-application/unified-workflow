package io.unifiedworkflow.sdk.errors;

import io.unifiedworkflow.sdk.models.ValidationError;

import java.util.Collections;
import java.util.List;
import java.util.Objects;

/**
 * Thrown when request data fails one or more validation rules.
 *
 * <p>The list of individual violations is available via {@link #getValidationErrors()}.</p>
 */
public class ValidationException extends WorkflowSDKException {

    private static final long serialVersionUID = 1L;

    private final List<ValidationError> validationErrors;

    /**
     * Constructs a new validation exception.
     *
     * @param message          human-readable summary of the validation failure
     * @param validationErrors list of individual field-level errors; must not be null
     */
    public ValidationException(String message, List<ValidationError> validationErrors) {
        super(ErrorCode.VALIDATION_FAILED, message, 400);
        this.validationErrors = Collections.unmodifiableList(
                Objects.requireNonNull(validationErrors, "validationErrors must not be null"));
    }

    /**
     * Returns the list of individual field-level validation errors.
     *
     * @return immutable list of errors, never null, may be empty
     */
    public List<ValidationError> getValidationErrors() {
        return validationErrors;
    }
}
