package mem0

import (
	"net/http"
)

// Client is the main entry point for the mem0 SDK.
type Client struct {
	memory Memory
	config *Config
}

// NewClient creates a new mem0 client with the given options.
// If no API key is provided via options, it will be loaded from the
// MEM0_API_KEY environment variable.
func NewClient(opts ...Option) (*Client, error) {
	// Load defaults from environment
	envConfig := LoadConfigFromEnv()

	// Apply options
	options := &clientOptions{
		apiKey:     envConfig.APIKey,
		baseURL:    envConfig.BaseURL,
		httpClient: http.DefaultClient,
		backend:    BackendHosted,
	}
	for _, opt := range opts {
		opt(options)
	}

	// Set default base URL based on backend if not specified
	if options.baseURL == "" {
		switch options.backend {
		case BackendHosted:
			options.baseURL = DefaultHostedBaseURL
		case BackendFOSS:
			options.baseURL = DefaultFOSSBaseURL
		default:
			return nil, ErrInvalidBackend
		}
	}

	// Validate API key for hosted backend
	if options.backend == BackendHosted && options.apiKey == "" {
		return nil, ErrNoAPIKey
	}

	// Build config
	config := &Config{
		APIKey:  options.apiKey,
		BaseURL: options.baseURL,
		Backend: options.backend,
	}

	// Create backend-specific memory implementation
	var memory Memory
	switch options.backend {
	case BackendHosted:
		m, err := newHostedMemory(options.baseURL, options.apiKey, options.httpClient)
		if err != nil {
			return nil, err
		}
		memory = m
	case BackendFOSS:
		return nil, ErrInvalidBackend // FOSS not yet implemented
	default:
		return nil, ErrInvalidBackend
	}

	return &Client{
		memory: memory,
		config: config,
	}, nil
}

// Memory returns the Memory interface for performing memory operations.
func (c *Client) Memory() Memory {
	return c.memory
}

// Config returns the client configuration.
func (c *Client) Config() *Config {
	return c.config
}

// authHTTPClient wraps an HTTP client to add authentication headers.
type authHTTPClient struct {
	client *http.Client
	apiKey string
}

// Do implements the ogen HTTPClient interface.
func (c *authHTTPClient) Do(req *http.Request) (*http.Response, error) {
	if c.apiKey != "" {
		req.Header.Set("Authorization", "Token "+c.apiKey)
	}
	req.Header.Set("X-Mem0-SDK-Version", Version)
	req.Header.Set("X-Mem0-SDK-Lang", "go")
	//nolint:gosec // G704: Internal wrapper, URLs controlled by ogen-generated client
	return c.client.Do(req)
}
