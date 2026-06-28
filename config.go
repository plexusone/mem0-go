package mem0

import "os"

const (
	// DefaultHostedBaseURL is the default base URL for the hosted mem0 platform.
	DefaultHostedBaseURL = "https://api.mem0.ai"

	// DefaultFOSSBaseURL is the default base URL for a self-hosted mem0 deployment.
	DefaultFOSSBaseURL = "http://localhost:8888"

	// EnvAPIKey is the environment variable name for the API key.
	EnvAPIKey = "MEM0_API_KEY"

	// EnvBaseURL is the environment variable name for the base URL.
	EnvBaseURL = "MEM0_BASE_URL"
)

// Config holds the client configuration.
type Config struct {
	APIKey  string
	BaseURL string
	Backend Backend
}

// LoadConfigFromEnv loads configuration from environment variables.
func LoadConfigFromEnv() *Config {
	return &Config{
		APIKey:  os.Getenv(EnvAPIKey),
		BaseURL: os.Getenv(EnvBaseURL),
	}
}
