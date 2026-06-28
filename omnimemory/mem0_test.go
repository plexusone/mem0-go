package omnimemory

import (
	"testing"

	"github.com/plexusone/omnimemory/core"
)

func TestNewProvider_NoAPIKey(t *testing.T) {
	t.Setenv("MEM0_API_KEY", "")

	_, err := NewProvider(core.ProviderConfig{}, nil)
	if err == nil {
		t.Fatal("expected error when no API key is provided")
	}

	valErr, ok := err.(*core.ValidationError)
	if !ok {
		t.Fatalf("expected ValidationError, got: %T", err)
	}
	if valErr.Field != "api_key" {
		t.Errorf("expected field 'api_key', got: %s", valErr.Field)
	}
}

func TestNewProvider_WithAPIKey(t *testing.T) {
	p, err := NewProvider(core.ProviderConfig{
		APIKey: "test-api-key",
	}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p == nil {
		t.Fatal("expected provider to be non-nil")
	}
	defer func() { _ = p.Close() }()

	if p.Name() != "mem0" {
		t.Errorf("expected name 'mem0', got: %s", p.Name())
	}
}

func TestNewProvider_WithOptions(t *testing.T) {
	p, err := NewProvider(core.ProviderConfig{
		Options: map[string]any{
			"api_key":  "test-api-key",
			"base_url": "https://custom.example.com",
		},
	}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p == nil {
		t.Fatal("expected provider to be non-nil")
	}
	defer func() { _ = p.Close() }()
}

func TestNewProvider_FromEnv(t *testing.T) {
	t.Setenv("MEM0_API_KEY", "env-api-key")

	p, err := NewProvider(core.ProviderConfig{}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p == nil {
		t.Fatal("expected provider to be non-nil")
	}
	defer func() { _ = p.Close() }()
}

func TestProvider_Name(t *testing.T) {
	p, _ := NewProvider(core.ProviderConfig{APIKey: "test"}, nil)
	defer func() { _ = p.Close() }()

	name := p.Name()
	if name != "mem0" {
		t.Errorf("expected 'mem0', got: %s", name)
	}

	// Verify name format (lowercase alphanumeric with hyphens/underscores)
	for _, r := range name {
		isLowerAlpha := r >= 'a' && r <= 'z'
		isDigit := r >= '0' && r <= '9'
		isSpecial := r == '-' || r == '_'
		if !isLowerAlpha && !isDigit && !isSpecial {
			t.Errorf("Name() contains invalid character %q", r)
		}
	}
}

func TestProvider_Close(t *testing.T) {
	p, _ := NewProvider(core.ProviderConfig{APIKey: "test"}, nil)

	err := p.Close()
	if err != nil {
		t.Errorf("Close() error: %v", err)
	}
}

func TestExtractMemoryType(t *testing.T) {
	tests := []struct {
		name     string
		metadata map[string]interface{}
		want     core.MemoryType
	}{
		{
			name:     "nil metadata",
			metadata: nil,
			want:     core.MemoryTypeObservation,
		},
		{
			name:     "empty metadata",
			metadata: map[string]interface{}{},
			want:     core.MemoryTypeObservation,
		},
		{
			name: "with type",
			metadata: map[string]interface{}{
				"omnimemory_type": "fact",
			},
			want: core.MemoryTypeFact,
		},
		{
			name: "with preference type",
			metadata: map[string]interface{}{
				"omnimemory_type": "preference",
			},
			want: core.MemoryTypePreference,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := extractMemoryType(tc.metadata)
			if got != tc.want {
				t.Errorf("extractMemoryType() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestExtractScope(t *testing.T) {
	tests := []struct {
		name     string
		metadata map[string]interface{}
		want     core.Scope
	}{
		{
			name:     "nil metadata",
			metadata: nil,
			want:     core.ScopeUser,
		},
		{
			name:     "empty metadata",
			metadata: map[string]interface{}{},
			want:     core.ScopeUser,
		},
		{
			name: "with scope",
			metadata: map[string]interface{}{
				"omnimemory_scope": "agent",
			},
			want: core.ScopeAgent,
		},
		{
			name: "with tenant scope",
			metadata: map[string]interface{}{
				"omnimemory_scope": "tenant",
			},
			want: core.ScopeTenant,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := extractScope(tc.metadata)
			if got != tc.want {
				t.Errorf("extractScope() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestContainsType(t *testing.T) {
	types := []core.MemoryType{core.MemoryTypeFact, core.MemoryTypePreference}

	if !containsType(types, core.MemoryTypeFact) {
		t.Error("expected containsType to return true for fact")
	}
	if containsType(types, core.MemoryTypeObservation) {
		t.Error("expected containsType to return false for observation")
	}
}

func TestContainsScope(t *testing.T) {
	scopes := []core.Scope{core.ScopeUser, core.ScopeAgent}

	if !containsScope(scopes, core.ScopeUser) {
		t.Error("expected containsScope to return true for user")
	}
	if containsScope(scopes, core.ScopeTenant) {
		t.Error("expected containsScope to return false for tenant")
	}
}
