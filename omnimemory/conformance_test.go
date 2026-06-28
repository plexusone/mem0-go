package omnimemory

import (
	"os"
	"testing"

	"github.com/plexusone/omnimemory/core"
	"github.com/plexusone/omnimemory/core/providertest"
)

func TestConformance(t *testing.T) {
	apiKey := os.Getenv("MEM0_API_KEY")
	if apiKey == "" {
		t.Skip("MEM0_API_KEY must be set for conformance tests")
	}

	// TenantID maps to app_id in mem0
	// SubjectID maps to user_id in mem0
	tenantID := os.Getenv("MEM0_APP_ID")
	if tenantID == "" {
		tenantID = "omnimemory-test"
	}

	subjectID := os.Getenv("MEM0_USER_ID")
	if subjectID == "" {
		subjectID = "test-user"
	}

	p, err := NewProvider(core.ProviderConfig{
		APIKey: apiKey,
	}, nil)
	if err != nil {
		t.Fatalf("NewProvider() error: %v", err)
	}
	defer func() { _ = p.Close() }()

	providertest.RunAll(t, providertest.Config{
		Provider:  p,
		TenantID:  tenantID,
		SubjectID: subjectID,
	})
}
