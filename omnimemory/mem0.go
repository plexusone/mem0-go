// Package omnimemory provides a mem0 provider for omnimemory.
//
// This package implements the omnimemory Provider interface using the mem0
// API as the backend. It maps omnimemory concepts to mem0:
//
//   - TenantID → mem0 app_id (or ignored if single-tenant)
//   - SubjectID → mem0 user_id
//   - Memory → mem0 MemoryItem
//   - Memory.AgentID → mem0 agent_id
//   - Search/Recall → mem0 Search API
//
// # Usage
//
// Import this package to register the mem0 provider:
//
//	import (
//	    "github.com/plexusone/omnimemory"
//	    "github.com/plexusone/omnimemory/core"
//	    _ "github.com/plexusone/mem0-go/omnimemory" // Register mem0 provider
//	)
//
//	func main() {
//	    client, err := omnimemory.NewClient(core.ClientConfig{
//	        Providers: []core.ProviderConfig{
//	            {
//	                Name:   core.ProviderNameMem0,
//	                APIKey: "your-mem0-api-key",
//	            },
//	        },
//	    })
//	    // ...
//	}
//
// # Configuration
//
// The provider accepts the following options:
//
//   - api_key: mem0 API key (required, or use MEM0_API_KEY env)
//   - base_url: Custom base URL (optional, defaults to api.mem0.ai)
//
// # Mapping
//
// mem0 uses a flat structure with user_id for isolation:
//
//   - user_id: The subject that owns the memories (maps to SubjectID)
//   - agent_id: The agent that created the memory (maps to AgentID)
//   - app_id: Application identifier (maps to TenantID)
//
// The SubjectID in omnimemory is mapped to mem0's user_id.
// The TenantID is mapped to mem0's app_id for multi-tenant isolation.
package omnimemory

import (
	"context"
	"os"
	"time"

	"github.com/plexusone/mem0-go"
	"github.com/plexusone/omnimemory/core"
)

func init() {
	core.RegisterProvider(core.ProviderNameMem0, NewProvider, core.PriorityThick)
}

// Provider implements core.Provider using mem0 API.
type Provider struct {
	client   *mem0.Client
	embedder core.Embedder
	config   core.ProviderConfig
}

// NewProvider creates a new mem0 Provider.
func NewProvider(config core.ProviderConfig, embedder core.Embedder) (core.Provider, error) {
	apiKey := config.APIKey
	if apiKey == "" {
		apiKey = getOption(config.Options, "api_key", os.Getenv("MEM0_API_KEY"))
	}

	if apiKey == "" {
		return nil, core.NewValidationError("api_key", "mem0 API key is required")
	}

	baseURL := config.Endpoint
	if baseURL == "" {
		baseURL = getOption(config.Options, "base_url", "")
	}

	opts := []mem0.Option{mem0.WithAPIKey(apiKey)}
	if baseURL != "" {
		opts = append(opts, mem0.WithBaseURL(baseURL))
	}

	client, err := mem0.NewClient(opts...)
	if err != nil {
		return nil, core.NewProviderError("mem0", "NewClient", err)
	}

	return &Provider{
		client:   client,
		embedder: embedder,
		config:   config,
	}, nil
}

// Name returns the provider name.
func (p *Provider) Name() string {
	return core.ProviderNameMem0.String()
}

// Close closes the provider.
func (p *Provider) Close() error {
	// mem0 client doesn't need explicit closing
	return nil
}

// Add adds a new memory to mem0.
func (p *Provider) Add(ctx context.Context, req *core.AddRequest) (*core.Memory, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	now := time.Now()

	// Build messages for mem0
	messages := []mem0.Message{
		{Role: mem0.RoleUser, Content: req.Content},
	}

	// Build options
	opts := []mem0.AddOption{
		mem0.WithUserID(req.SubjectID),
		mem0.WithAppID(req.TenantID),
	}

	if req.AgentID != "" {
		opts = append(opts, mem0.WithAgentID(req.AgentID))
	}
	if req.SessionID != "" {
		opts = append(opts, mem0.WithRunID(req.SessionID))
	}

	// Add metadata including type and scope
	metadata := make(map[string]interface{})
	if req.Metadata != nil {
		for k, v := range req.Metadata {
			metadata[k] = v
		}
	}
	metadata["omnimemory_type"] = string(req.Type)
	metadata["omnimemory_scope"] = string(req.Scope)
	if req.TTL > 0 {
		expiresAt := now.Add(req.TTL)
		metadata["expires_at"] = expiresAt.Format(time.RFC3339)
	}
	opts = append(opts, mem0.WithMetadata(metadata))

	resp, err := p.client.Memory().Add(ctx, messages, opts...)
	if err != nil {
		return nil, core.NewProviderError(p.Name(), "Add", err)
	}

	// Extract memory ID from results
	var memoryID string
	if len(resp.Results) > 0 {
		memoryID = resp.Results[0].ID
	}

	memory := &core.Memory{
		ID:        memoryID,
		TenantID:  req.TenantID,
		SubjectID: req.SubjectID,
		AgentID:   req.AgentID,
		SessionID: req.SessionID,
		Scope:     req.Scope,
		Type:      req.Type,
		Content:   req.Content,
		Metadata:  req.Metadata,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if req.TTL > 0 {
		expiresAt := now.Add(req.TTL)
		memory.ExpiresAt = &expiresAt
	}

	return memory, nil
}

// Get retrieves a memory by ID.
func (p *Provider) Get(ctx context.Context, req *core.GetRequest) (*core.Memory, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	item, err := p.client.Memory().Get(ctx, req.ID)
	if err != nil {
		if mem0.IsNotFoundError(err) {
			return nil, core.ErrNotFound
		}
		return nil, core.NewProviderError(p.Name(), "Get", err)
	}

	// Verify isolation
	if item.UserID != req.SubjectID {
		return nil, core.ErrNotFound
	}
	if item.AppID != "" && item.AppID != req.TenantID {
		return nil, core.ErrNotFound
	}

	return memoryItemToCore(item, req.TenantID), nil
}

// Update updates an existing memory.
func (p *Provider) Update(ctx context.Context, req *core.UpdateRequest) (*core.Memory, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	// Get existing memory first to verify it exists and check isolation
	existing, err := p.Get(ctx, &core.GetRequest{
		Context: req.Context,
		ID:      req.ID,
	})
	if err != nil {
		return nil, err
	}

	// Determine content to update
	content := existing.Content
	if req.Content != "" {
		content = req.Content
	}

	item, err := p.client.Memory().Update(ctx, req.ID, content)
	if err != nil {
		if mem0.IsNotFoundError(err) {
			return nil, core.ErrNotFound
		}
		return nil, core.NewProviderError(p.Name(), "Update", err)
	}

	memory := memoryItemToCore(item, req.TenantID)

	// Merge metadata
	if req.Metadata != nil {
		if memory.Metadata == nil {
			memory.Metadata = make(map[string]any)
		}
		for k, v := range req.Metadata {
			memory.Metadata[k] = v
		}
	}

	return memory, nil
}

// Delete deletes a memory by ID.
func (p *Provider) Delete(ctx context.Context, req *core.DeleteRequest) error {
	if err := req.Validate(); err != nil {
		return err
	}

	// Verify isolation first
	_, err := p.Get(ctx, &core.GetRequest{
		Context: req.Context,
		ID:      req.ID,
	})
	if err != nil {
		return err
	}

	err = p.client.Memory().Delete(ctx, req.ID)
	if err != nil {
		if mem0.IsNotFoundError(err) {
			return core.ErrNotFound
		}
		return core.NewProviderError(p.Name(), "Delete", err)
	}

	return nil
}

// List lists memories with optional filters.
func (p *Provider) List(ctx context.Context, req *core.ListRequest) (*core.ListResponse, error) {
	if req.TenantID == "" {
		return nil, core.ErrTenantRequired
	}
	if req.SubjectID == "" {
		return nil, core.ErrSubjectRequired
	}

	limit := req.Limit
	if limit <= 0 {
		limit = 100
	}

	page := 1
	if req.Offset > 0 && limit > 0 {
		page = (req.Offset / limit) + 1
	}

	opts := []mem0.GetAllOption{
		mem0.WithGetAllFilters(mem0.Filters{
			UserID: req.SubjectID,
			AppID:  req.TenantID,
		}),
		mem0.WithPage(page),
		mem0.WithPageSize(limit),
	}

	resp, err := p.client.Memory().GetAll(ctx, opts...)
	if err != nil {
		return nil, core.NewProviderError(p.Name(), "List", err)
	}

	memories := make([]*core.Memory, 0, len(resp.Results))
	for i := range resp.Results {
		item := &resp.Results[i]

		// Apply type filter
		if len(req.Types) > 0 {
			memType := extractMemoryType(item.Metadata)
			if !containsType(req.Types, memType) {
				continue
			}
		}

		// Apply scope filter
		if len(req.Scopes) > 0 {
			memScope := extractScope(item.Metadata)
			if !containsScope(req.Scopes, memScope) {
				continue
			}
		}

		memories = append(memories, memoryItemToCore(item, req.TenantID))
	}

	return &core.ListResponse{
		Memories:   memories,
		TotalCount: resp.Count,
		HasMore:    resp.Next != "",
	}, nil
}

// Search performs semantic search using mem0 Search API.
func (p *Provider) Search(ctx context.Context, req *core.SearchRequest) (*core.SearchResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	limit := req.Limit
	if limit <= 0 {
		limit = 10
	}

	opts := []mem0.SearchOption{
		mem0.WithFilters(mem0.Filters{
			UserID: req.SubjectID,
			AppID:  req.TenantID,
		}),
		mem0.WithTopK(limit),
	}

	if req.Threshold > 0 {
		opts = append(opts, mem0.WithThreshold(req.Threshold))
	}

	results, err := p.client.Memory().Search(ctx, req.Query, opts...)
	if err != nil {
		return nil, core.NewProviderError(p.Name(), "Search", err)
	}

	var searchResults []*core.SearchResult
	for i := range results {
		r := &results[i]

		// Apply type filter
		if len(req.Types) > 0 {
			memType := extractMemoryType(r.Metadata)
			if !containsType(req.Types, memType) {
				continue
			}
		}

		// Apply scope filter
		if len(req.Scopes) > 0 {
			memScope := extractScope(r.Metadata)
			if !containsScope(req.Scopes, memScope) {
				continue
			}
		}

		searchResults = append(searchResults, &core.SearchResult{
			Memory: searchResultToCore(r, req.TenantID),
			Score:  r.Score,
		})
	}

	return &core.SearchResponse{
		Results: searchResults,
	}, nil
}

// Recall retrieves relevant memories using mem0 Search API.
func (p *Provider) Recall(ctx context.Context, req *core.RecallRequest) (*core.RecallResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	maxResults := req.MaxResults
	if maxResults <= 0 {
		maxResults = 20
	}

	opts := []mem0.SearchOption{
		mem0.WithFilters(mem0.Filters{
			UserID: req.SubjectID,
			AppID:  req.TenantID,
		}),
		mem0.WithTopK(maxResults),
	}

	results, err := p.client.Memory().Search(ctx, req.Query, opts...)
	if err != nil {
		return nil, core.NewProviderError(p.Name(), "Recall", err)
	}

	var memories []*core.Memory
	for i := range results {
		r := &results[i]

		// Apply type filter
		if len(req.IncludeTypes) > 0 {
			memType := extractMemoryType(r.Metadata)
			if !containsType(req.IncludeTypes, memType) {
				continue
			}
		}

		memories = append(memories, searchResultToCore(r, req.TenantID))
	}

	return &core.RecallResponse{
		Memories: memories,
	}, nil
}

// memoryItemToCore converts a mem0.MemoryItem to core.Memory.
func memoryItemToCore(item *mem0.MemoryItem, tenantID string) *core.Memory {
	m := &core.Memory{
		ID:        item.ID,
		TenantID:  tenantID,
		SubjectID: item.UserID,
		AgentID:   item.AgentID,
		SessionID: item.RunID,
		Type:      extractMemoryType(item.Metadata),
		Scope:     extractScope(item.Metadata),
		Content:   item.Memory,
		Metadata:  item.Metadata,
		CreatedAt: item.CreatedAt,
		UpdatedAt: item.UpdatedAt,
	}

	// Extract expiration from metadata
	if item.Metadata != nil {
		if expiresAtStr, ok := item.Metadata["expires_at"].(string); ok {
			if expiresAt, err := time.Parse(time.RFC3339, expiresAtStr); err == nil {
				m.ExpiresAt = &expiresAt
			}
		}
	}

	return m
}

// searchResultToCore converts a mem0.SearchResult to core.Memory.
func searchResultToCore(r *mem0.SearchResult, tenantID string) *core.Memory {
	m := &core.Memory{
		ID:        r.ID,
		TenantID:  tenantID,
		SubjectID: r.UserID,
		AgentID:   r.AgentID,
		SessionID: r.RunID,
		Type:      extractMemoryType(r.Metadata),
		Scope:     extractScope(r.Metadata),
		Content:   r.Memory,
		Metadata:  r.Metadata,
		CreatedAt: r.CreatedAt,
		UpdatedAt: r.UpdatedAt,
	}

	// Extract expiration from metadata
	if r.Metadata != nil {
		if expiresAtStr, ok := r.Metadata["expires_at"].(string); ok {
			if expiresAt, err := time.Parse(time.RFC3339, expiresAtStr); err == nil {
				m.ExpiresAt = &expiresAt
			}
		}
	}

	return m
}

// extractMemoryType extracts the omnimemory type from metadata.
func extractMemoryType(metadata map[string]interface{}) core.MemoryType {
	if metadata == nil {
		return core.MemoryTypeObservation
	}
	if t, ok := metadata["omnimemory_type"].(string); ok {
		return core.MemoryType(t)
	}
	return core.MemoryTypeObservation
}

// extractScope extracts the omnimemory scope from metadata.
func extractScope(metadata map[string]interface{}) core.Scope {
	if metadata == nil {
		return core.ScopeUser
	}
	if s, ok := metadata["omnimemory_scope"].(string); ok {
		return core.Scope(s)
	}
	return core.ScopeUser
}

// containsType checks if a memory type is in the list.
func containsType(types []core.MemoryType, t core.MemoryType) bool {
	for _, mt := range types {
		if mt == t {
			return true
		}
	}
	return false
}

// containsScope checks if a scope is in the list.
func containsScope(scopes []core.Scope, s core.Scope) bool {
	for _, sc := range scopes {
		if sc == s {
			return true
		}
	}
	return false
}

// getOption retrieves an option from the map with a default value.
func getOption(options map[string]any, key, defaultValue string) string {
	if options == nil {
		return defaultValue
	}
	if v, ok := options[key]; ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return defaultValue
}
