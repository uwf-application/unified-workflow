package io.unifiedworkflow.sdk.models;

import com.fasterxml.jackson.annotation.JsonIgnoreProperties;
import com.fasterxml.jackson.annotation.JsonProperty;

import java.time.Instant;
import java.util.Collections;
import java.util.List;
import java.util.Map;

/**
 * Metadata and definition of a registered workflow.
 *
 * @param id          unique workflow identifier
 * @param name        human-readable workflow name
 * @param description optional description of the workflow's purpose
 * @param steps       ordered list of step definition maps; never null, may be empty
 * @param metadata    arbitrary key-value metadata attached to the workflow definition;
 *                    never null, may be empty
 * @param createdAt   time at which the workflow was registered
 * @param updatedAt   time of the most recent definition update; may be null if never updated
 */
@JsonIgnoreProperties(ignoreUnknown = true)
public record WorkflowDefinition(
        @JsonProperty("id") String id,
        @JsonProperty("name") String name,
        @JsonProperty("description") String description,
        @JsonProperty("steps") List<Map<String, Object>> steps,
        @JsonProperty("metadata") Map<String, String> metadata,
        @JsonProperty("created_at") Instant createdAt,
        @JsonProperty("updated_at") Instant updatedAt
) {

    /**
     * Compact constructor that ensures collections are never null.
     */
    public WorkflowDefinition {
        steps = steps != null ? Collections.unmodifiableList(steps) : Collections.emptyList();
        metadata = metadata != null ? Collections.unmodifiableMap(metadata) : Collections.emptyMap();
    }
}
