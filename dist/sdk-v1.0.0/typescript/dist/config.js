/**
 * Configuration management for Unified Workflow SDK
 */
import { AuthType, LogLevel } from './models';
import { SDKConfigurationError } from './errors';
export const DEFAULT_CONFIG = {
    sdkVersion: '1.0.0',
    timeoutMs: 30000,
    maxRetries: 3,
    retryDelayMs: 1000,
    authType: AuthType.BEARER_TOKEN,
    enableValidation: true,
    enableSanitization: true,
    strictValidation: false,
    enableSessionExtraction: true,
    enableSecurityContext: true,
    includeFullHttpContext: true,
    logLevel: LogLevel.INFO,
    enableRequestLogging: true,
    enableMetrics: true,
    defaultValidationRules: [],
    customValidators: [],
    // Execution configuration
    asyncExecution: true,
    defaultPriority: 5,
    pollIntervalMs: 1000,
    estimatedCompletionMs: 5000,
    executionExpiryDurationMs: 3600000, // 1 hour in milliseconds
    // Circuit breaker configuration
    enableCircuitBreaker: true,
    circuitBreakerThreshold: 5,
    circuitBreakerTimeoutMs: 60000, // 60 seconds in milliseconds
    // HTTP client configuration
    maxRedirects: 5
};
export class SDKConfig {
    constructor(config) {
        // Merge with defaults
        this.config = {
            ...DEFAULT_CONFIG,
            ...config
        };
        // Validate configuration
        this.validate();
    }
    /**
     * Validate the SDK configuration
     */
    validate() {
        const errors = [];
        // Check required fields
        if (!this.config.workflowApiEndpoint) {
            errors.push('workflowApiEndpoint is required');
        }
        if (!this.config.sdkVersion) {
            errors.push('sdkVersion is required');
        }
        // Validate URL format
        if (this.config.workflowApiEndpoint) {
            try {
                new URL(this.config.workflowApiEndpoint);
            }
            catch (error) {
                errors.push('workflowApiEndpoint must be a valid URL');
            }
        }
        // Validate numeric fields
        if (this.config.timeoutMs !== undefined && this.config.timeoutMs <= 0) {
            errors.push('timeoutMs must be positive');
        }
        if (this.config.maxRetries !== undefined && this.config.maxRetries < 0) {
            errors.push('maxRetries cannot be negative');
        }
        if (this.config.retryDelayMs !== undefined && this.config.retryDelayMs < 0) {
            errors.push('retryDelayMs cannot be negative');
        }
        // Validate auth configuration
        if (this.config.authType === AuthType.BEARER_TOKEN && !this.config.authToken) {
            errors.push('authToken is required for bearer token authentication');
        }
        if (this.config.authType === AuthType.API_KEY && !this.config.authToken) {
            errors.push('authToken is required for API key authentication');
        }
        // Throw if there are validation errors
        if (errors.length > 0) {
            throw new SDKConfigurationError('Invalid SDK configuration', undefined, { errors });
        }
    }
    /**
     * Get the complete configuration
     */
    getConfig() {
        return { ...this.config };
    }
    /**
     * Update configuration
     */
    updateConfig(updates) {
        this.config = {
            ...this.config,
            ...updates
        };
        this.validate();
    }
    /**
     * Get workflow API endpoint
     */
    getWorkflowApiEndpoint() {
        return this.config.workflowApiEndpoint;
    }
    /**
     * Get authentication configuration
     */
    getAuthConfig() {
        return {
            type: this.config.authType || AuthType.BEARER_TOKEN,
            token: this.config.authToken
        };
    }
    /**
     * Get timeout configuration
     */
    getTimeoutConfig() {
        return {
            timeoutMs: this.config.timeoutMs || 30000,
            maxRetries: this.config.maxRetries || 3,
            retryDelayMs: this.config.retryDelayMs || 1000
        };
    }
    /**
     * Get validation configuration
     */
    getValidationConfig() {
        return {
            enableValidation: this.config.enableValidation || true,
            enableSanitization: this.config.enableSanitization || true,
            strictValidation: this.config.strictValidation || false,
            defaultValidationRules: this.config.defaultValidationRules || []
        };
    }
    /**
     * Get logging configuration
     */
    getLoggingConfig() {
        return {
            logLevel: this.config.logLevel || LogLevel.INFO,
            enableRequestLogging: this.config.enableRequestLogging || true
        };
    }
    /**
     * Check if metrics are enabled
     */
    isMetricsEnabled() {
        return this.config.enableMetrics || true;
    }
    /**
     * Check if session extraction is enabled
     */
    isSessionExtractionEnabled() {
        return this.config.enableSessionExtraction || true;
    }
    /**
     * Check if security context is enabled
     */
    isSecurityContextEnabled() {
        return this.config.enableSecurityContext || true;
    }
    /**
     * Check if full HTTP context should be included
     */
    shouldIncludeFullHttpContext() {
        return this.config.includeFullHttpContext || true;
    }
    /**
     * Get custom validators
     */
    getCustomValidators() {
        return this.config.customValidators || [];
    }
    /**
     * Create headers for API requests
     */
    createHeaders() {
        const headers = {
            'Content-Type': 'application/json',
            'User-Agent': `UnifiedWorkflowSDK/${this.config.sdkVersion}`,
            'X-SDK-Version': this.config.sdkVersion
        };
        // Add authentication header
        const authConfig = this.getAuthConfig();
        switch (authConfig.type) {
            case AuthType.BEARER_TOKEN:
                if (authConfig.token) {
                    headers['Authorization'] = `Bearer ${authConfig.token}`;
                }
                break;
            case AuthType.API_KEY:
                if (authConfig.token) {
                    headers['X-API-Key'] = authConfig.token;
                }
                break;
            case AuthType.BASIC_AUTH:
                if (authConfig.token) {
                    headers['Authorization'] = `Basic ${authConfig.token}`;
                }
                break;
            // Note: OAUTH2 and AWS_SIGV4 would require more complex handling
        }
        return headers;
    }
    /**
     * Create request options for fetch/axios
     */
    createRequestOptions() {
        return {
            timeout: this.config.timeoutMs || 30000,
            headers: this.createHeaders(),
            maxRedirects: this.config.maxRedirects || 5
        };
    }
    /**
     * Create a new configuration from environment variables
     * Note: This method requires Node.js environment
     */
    static fromEnvironment() {
        // Check if we're in a Node.js environment
        if (typeof process === 'undefined' || !process.env) {
            throw new Error('Environment variables are not available in this runtime');
        }
        const envConfig = {
            workflowApiEndpoint: process.env.SDK_WORKFLOW_API_ENDPOINT,
            sdkVersion: process.env.SDK_VERSION || '1.0.0',
            timeoutMs: process.env.SDK_TIMEOUT_MS ? parseInt(process.env.SDK_TIMEOUT_MS, 10) : undefined,
            maxRetries: process.env.SDK_MAX_RETRIES ? parseInt(process.env.SDK_MAX_RETRIES, 10) : undefined,
            retryDelayMs: process.env.SDK_RETRY_DELAY_MS ? parseInt(process.env.SDK_RETRY_DELAY_MS, 10) : undefined,
            authToken: process.env.SDK_AUTH_TOKEN,
            authType: process.env.SDK_AUTH_TYPE,
            enableValidation: process.env.SDK_ENABLE_VALIDATION !== 'false',
            enableSanitization: process.env.SDK_ENABLE_SANITIZATION !== 'false',
            strictValidation: process.env.SDK_STRICT_VALIDATION === 'true',
            logLevel: process.env.SDK_LOG_LEVEL,
            enableRequestLogging: process.env.SDK_ENABLE_REQUEST_LOGGING !== 'false',
            enableMetrics: process.env.SDK_ENABLE_METRICS !== 'false',
            asyncExecution: process.env.SDK_ASYNC_EXECUTION !== 'false',
            defaultPriority: process.env.SDK_DEFAULT_PRIORITY ? parseInt(process.env.SDK_DEFAULT_PRIORITY, 10) : undefined,
            pollIntervalMs: process.env.SDK_POLL_INTERVAL_MS ? parseInt(process.env.SDK_POLL_INTERVAL_MS, 10) : undefined,
            estimatedCompletionMs: process.env.SDK_ESTIMATED_COMPLETION_MS ? parseInt(process.env.SDK_ESTIMATED_COMPLETION_MS, 10) : undefined,
            executionExpiryDurationMs: process.env.SDK_EXECUTION_EXPIRY_DURATION_MS ? parseInt(process.env.SDK_EXECUTION_EXPIRY_DURATION_MS, 10) : undefined,
            enableCircuitBreaker: process.env.SDK_ENABLE_CIRCUIT_BREAKER !== 'false',
            circuitBreakerThreshold: process.env.SDK_CIRCUIT_BREAKER_THRESHOLD ? parseInt(process.env.SDK_CIRCUIT_BREAKER_THRESHOLD, 10) : undefined,
            circuitBreakerTimeoutMs: process.env.SDK_CIRCUIT_BREAKER_TIMEOUT_MS ? parseInt(process.env.SDK_CIRCUIT_BREAKER_TIMEOUT_MS, 10) : undefined,
            maxRedirects: process.env.SDK_MAX_REDIRECTS ? parseInt(process.env.SDK_MAX_REDIRECTS, 10) : undefined
        };
        // Filter out undefined values
        const filteredConfig = Object.fromEntries(Object.entries(envConfig).filter(([_, value]) => value !== undefined));
        return new SDKConfig(filteredConfig);
    }
    /**
     * Create a new configuration with overrides
     */
    static merge(base, overrides) {
        return new SDKConfig({
            ...base,
            ...overrides
        });
    }
    /**
     * Load configuration from a JSON/YAML file
     * Note: This method requires Node.js environment with fs module
     */
    static async fromConfigFile(filePath) {
        // Check if we're in a Node.js environment
        if (typeof require === 'undefined') {
            throw new Error('File system access is not available in this runtime');
        }
        const fs = require('fs');
        const path = require('path');
        // Read and parse the config file
        const fileContent = fs.readFileSync(filePath, 'utf8');
        const ext = path.extname(filePath).toLowerCase();
        let configData;
        if (ext === '.json') {
            configData = JSON.parse(fileContent);
        }
        else if (ext === '.yaml' || ext === '.yml') {
            // Try to load yaml if available
            try {
                const yaml = require('yaml');
                configData = yaml.parse(fileContent);
            }
            catch (error) {
                throw new Error('YAML parsing requires the "yaml" package. Install it with: npm install yaml');
            }
        }
        else {
            throw new Error(`Unsupported config file format: ${ext}. Use .json, .yaml, or .yml`);
        }
        return new SDKConfig(configData);
    }
    /**
     * Load configuration from default locations
     * Looks for config files in the following order:
     * 1. Environment variables
     * 2. sdk-config.json in current directory
     * 3. sdk-config.yaml in current directory
     * 4. .sdkrc in current directory
     * 5. Default configuration
     */
    static async loadDefault() {
        // Check if we're in a Node.js environment
        if (typeof process === 'undefined' || !process.env) {
            return new SDKConfig({});
        }
        const fs = require('fs');
        const path = require('path');
        // Try to load from environment variables first
        try {
            const envConfig = SDKConfig.fromEnvironment();
            return envConfig;
        }
        catch (error) {
            // Environment variables not available, try config files
        }
        // Try config files in order
        const configFiles = [
            'sdk-config.json',
            'sdk-config.yaml',
            'sdk-config.yml',
            '.sdkrc'
        ];
        for (const configFile of configFiles) {
            const filePath = path.join(process.cwd(), configFile);
            if (fs.existsSync(filePath)) {
                try {
                    return await SDKConfig.fromConfigFile(filePath);
                }
                catch (error) {
                    console.warn(`Failed to load config from ${configFile}:`, error);
                    continue;
                }
            }
        }
        // Fall back to default configuration
        return new SDKConfig({});
    }
}
// Helper function to create configuration
export function createConfig(config) {
    return new SDKConfig(config);
}
// Helper function to validate configuration without creating instance
export function validateConfig(config) {
    const errors = [];
    if (!config.workflowApiEndpoint) {
        errors.push('workflowApiEndpoint is required');
    }
    if (!config.sdkVersion) {
        errors.push('sdkVersion is required');
    }
    if (config.workflowApiEndpoint) {
        try {
            new URL(config.workflowApiEndpoint);
        }
        catch (error) {
            errors.push('workflowApiEndpoint must be a valid URL');
        }
    }
    if (config.timeoutMs !== undefined && config.timeoutMs <= 0) {
        errors.push('timeoutMs must be positive');
    }
    if (config.maxRetries !== undefined && config.maxRetries < 0) {
        errors.push('maxRetries cannot be negative');
    }
    if (config.retryDelayMs !== undefined && config.retryDelayMs < 0) {
        errors.push('retryDelayMs cannot be negative');
    }
    return errors;
}
//# sourceMappingURL=config.js.map