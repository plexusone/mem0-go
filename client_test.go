package mem0

import (
	"testing"
)

func TestNewClient_NoAPIKey(t *testing.T) {
	// Ensure env var is not set by using t.Setenv with empty value
	t.Setenv(EnvAPIKey, "")

	_, err := NewClient()
	if err == nil {
		t.Fatal("expected error when no API key is provided")
	}
	if err != ErrNoAPIKey {
		t.Fatalf("expected ErrNoAPIKey, got: %v", err)
	}
}

func TestNewClient_WithAPIKey(t *testing.T) {
	client, err := NewClient(WithAPIKey("test-api-key"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if client == nil {
		t.Fatal("expected client to be non-nil")
	}
	if client.Config().APIKey != "test-api-key" {
		t.Errorf("expected API key 'test-api-key', got: %s", client.Config().APIKey)
	}
	if client.Config().BaseURL != DefaultHostedBaseURL {
		t.Errorf("expected base URL %s, got: %s", DefaultHostedBaseURL, client.Config().BaseURL)
	}
}

func TestNewClient_WithBaseURL(t *testing.T) {
	client, err := NewClient(
		WithAPIKey("test-api-key"),
		WithBaseURL("https://custom.example.com"),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if client.Config().BaseURL != "https://custom.example.com" {
		t.Errorf("expected base URL 'https://custom.example.com', got: %s", client.Config().BaseURL)
	}
}

func TestNewClient_FromEnv(t *testing.T) {
	t.Setenv(EnvAPIKey, "env-api-key")

	client, err := NewClient()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if client.Config().APIKey != "env-api-key" {
		t.Errorf("expected API key 'env-api-key', got: %s", client.Config().APIKey)
	}
}

func TestNewClient_InvalidBackend(t *testing.T) {
	_, err := NewClient(
		WithAPIKey("test-api-key"),
		WithBackend(BackendFOSS),
	)
	if err != ErrInvalidBackend {
		t.Fatalf("expected ErrInvalidBackend, got: %v", err)
	}
}

func TestNewClient_Memory(t *testing.T) {
	client, err := NewClient(WithAPIKey("test-api-key"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if client.Memory() == nil {
		t.Fatal("expected Memory() to return non-nil")
	}
}
