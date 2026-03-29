/**
 * Configuration management for Unified Workflow SDK
 */
import { SDKConfiguration, AuthType, LogLevel, ValidationRule } from './models';
export declare const DEFAULT_CONFIG: Partial<SDKConfiguration>;
export declare class SDKConfig {
    private config;
    constructor(config: Partial<SDKConfiguration>);
    /**
     * Validate the SDK configuration
     */
    private validate;
    /**
     * Get the complete configuration
     */
    getConfig(): SDKConfiguration;
    /**
     * Update configuration
     */
    updateConfig(updates: Partial<SDKConfiguration>): void;
    /**
     * Get workflow API endpoint
     */
    getWorkflowApiEndpoint(): string;
    /**
     * Get authentication configuration
     */
    getAuthConfig(): {
        type: AuthType;
        token?: string;
    };
    /**
     * Get timeout configuration
     */
    getTimeoutConfig(): {
        timeoutMs: number;
        maxRetries: number;
        retryDelayMs: number;
    };
    /**
     * Get validation configuration
     */
    getValidationConfig(): {
        enableValidation: boolean;
        enableSanitization: boolean;
        strictValidation: boolean;
        defaultValidationRules: ValidationRule[];
    };
    /**
     * Get logging configuration
     */
    getLoggingConfig(): {
        logLevel: LogLevel;
        enableRequestLogging: boolean;
    };
    /**
     * Check if metrics are enabled
     */
    isMetricsEnabled(): boolean;
    /**
     * Check if session extraction is enabled
     */
    isSessionExtractionEnabled(): boolean;
    /**
     * Check if security context is enabled
     */
    isSecurityContextEnabled(): boolean;
    /**
     * Check if full HTTP context should be included
     */
    shouldIncludeFullHttpContext(): boolean;
    /**
     * Get custom validators
     */
    getCustomValidators(): string[];
    /**
     * Create headers for API requests
     */
    createHeaders(): Record<string, string>;
    /**
     * Create request options for fetch/axios
     */
    createRequestOptions(): {
        timeout: number;
        headers: Record<string, string>;
        maxRedirects?: number;
    };
    /**
     * Create a new configuration from environment variables
     * Note: This method requires Node.js environment
     */
    static fromEnvironment(): SDKConfig;
    /**
     * Create a new configuration with overrides
     */
    static merge(base: Partial<SDKConfiguration>, overrides: Partial<SDKConfiguration>): SDKConfig;
    /**
     * Load configuration from a JSON/YAML file
     * Note: This method requires Node.js environment with fs module
     */
    static fromConfigFile(filePath: string): Promise<SDKConfig>;
    /**
     * Load configuration from default locations
     * Looks for config files in the following order:
     * 1. Environment variables
     * 2. sdk-config.json in current directory
     * 3. sdk-config.yaml in current directory
     * 4. .sdkrc in current directory
     * 5. Default configuration
     */
    static loadDefault(): Promise<SDKConfig>;
}
export declare function createConfig(config: Partial<SDKConfiguration>): SDKConfig;
export declare function validateConfig(config: Partial<SDKConfiguration>): string[];
//# sourceMappingURL=config.d.ts.map