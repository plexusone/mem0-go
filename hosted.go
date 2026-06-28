package mem0

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/plexusone/mem0-go/internal/ogenhosted"
)

// hostedMemory implements the Memory interface for the hosted mem0 platform.
type hostedMemory struct {
	client *ogenhosted.Client
}

// noopSecuritySource is a no-op security source since we handle auth via HTTP client.
type noopSecuritySource struct{}

func (noopSecuritySource) TokenAuth(context.Context, string) (ogenhosted.TokenAuth, error) {
	return ogenhosted.TokenAuth{}, nil
}

// newHostedMemory creates a new hosted memory implementation.
func newHostedMemory(baseURL, apiKey string, httpClient *http.Client) (*hostedMemory, error) {
	authClient := &authHTTPClient{
		client: httpClient,
		apiKey: apiKey,
	}

	client, err := ogenhosted.NewClient(baseURL, noopSecuritySource{},
		ogenhosted.WithClient(authClient),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create ogen client: %w", err)
	}

	return &hostedMemory{client: client}, nil
}

// Add implements Memory.Add.
func (m *hostedMemory) Add(ctx context.Context, messages []Message, opts ...AddOption) (*AddResponse, error) {
	options := &addOptions{}
	for _, opt := range opts {
		opt(options)
	}

	// Convert messages to ogen types
	ogenMessages := make([]ogenhosted.Message, len(messages))
	for i, msg := range messages {
		ogenMessages[i] = ogenhosted.Message{
			Role:    ogenhosted.MessageRole(msg.Role),
			Content: msg.Content,
		}
	}

	req := &ogenhosted.AddMemoryRequest{
		Messages: ogenMessages,
	}

	if options.userID != "" {
		req.UserID = ogenhosted.NewOptString(options.userID)
	}
	if options.agentID != "" {
		req.AgentID = ogenhosted.NewOptString(options.agentID)
	}
	if options.appID != "" {
		req.AppID = ogenhosted.NewOptString(options.appID)
	}
	if options.runID != "" {
		req.RunID = ogenhosted.NewOptString(options.runID)
	}
	if options.infer != nil {
		req.Infer = ogenhosted.NewOptBool(*options.infer)
	}
	if options.metadata != nil {
		meta := make(ogenhosted.AddMemoryRequestMetadata)
		for k, v := range options.metadata {
			data, err := json.Marshal(v)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal metadata value for key %s: %w", k, err)
			}
			meta[k] = data
		}
		req.Metadata = ogenhosted.NewOptAddMemoryRequestMetadata(meta)
	}

	res, err := m.client.AddMemory(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to add memory: %w", err)
	}

	switch r := res.(type) {
	case *ogenhosted.AddMemoryResponse:
		return convertAddMemoryResponse(r), nil
	case *ogenhosted.AddMemoryBadRequest:
		return nil, &APIError{StatusCode: http.StatusBadRequest, Detail: r.Detail.Or("")}
	case *ogenhosted.AddMemoryUnauthorized:
		return nil, &APIError{StatusCode: http.StatusUnauthorized, Detail: r.Detail.Or("")}
	default:
		return nil, fmt.Errorf("unexpected response type: %T", res)
	}
}

// Get implements Memory.Get.
func (m *hostedMemory) Get(ctx context.Context, memoryID string) (*MemoryItem, error) {
	res, err := m.client.GetMemory(ctx, ogenhosted.GetMemoryParams{MemoryID: memoryID})
	if err != nil {
		return nil, fmt.Errorf("failed to get memory: %w", err)
	}

	switch r := res.(type) {
	case *ogenhosted.MemoryItem:
		return convertMemoryItem(r), nil
	case *ogenhosted.ErrorResponse:
		return nil, &APIError{StatusCode: http.StatusNotFound, Detail: r.Detail.Or("")}
	default:
		return nil, fmt.Errorf("unexpected response type: %T", res)
	}
}

// GetAll implements Memory.GetAll.
func (m *hostedMemory) GetAll(ctx context.Context, opts ...GetAllOption) (*GetAllResult, error) {
	options := &getAllOptions{}
	for _, opt := range opts {
		opt(options)
	}

	var req ogenhosted.OptGetMemoriesRequest
	if options.filters != nil || options.page > 0 || options.pageSize > 0 {
		r := ogenhosted.GetMemoriesRequest{}
		if options.filters != nil {
			r.Filters = ogenhosted.NewOptFilters(convertFiltersToOgen(*options.filters))
		}
		if options.page > 0 {
			r.Page = ogenhosted.NewOptInt(options.page)
		}
		if options.pageSize > 0 {
			r.PageSize = ogenhosted.NewOptInt(options.pageSize)
		}
		req = ogenhosted.NewOptGetMemoriesRequest(r)
	}

	res, err := m.client.GetMemories(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get memories: %w", err)
	}

	switch r := res.(type) {
	case *ogenhosted.GetMemoriesResponse:
		return convertGetMemoriesResponse(r), nil
	case *ogenhosted.ErrorResponse:
		return nil, &APIError{StatusCode: http.StatusBadRequest, Detail: r.Detail.Or("")}
	default:
		return nil, fmt.Errorf("unexpected response type: %T", res)
	}
}

// Search implements Memory.Search.
func (m *hostedMemory) Search(ctx context.Context, query string, opts ...SearchOption) ([]SearchResult, error) {
	options := &searchOptions{}
	for _, opt := range opts {
		opt(options)
	}

	req := &ogenhosted.SearchMemoriesRequest{
		Query: query,
	}

	if options.filters != nil {
		req.Filters = ogenhosted.NewOptFilters(convertFiltersToOgen(*options.filters))
	}
	if options.topK > 0 {
		req.TopK = ogenhosted.NewOptInt(options.topK)
	}
	if options.rerank {
		req.Rerank = ogenhosted.NewOptBool(true)
	}
	if options.threshold != nil {
		req.Threshold = ogenhosted.NewOptFloat32(float32(*options.threshold))
	}

	res, err := m.client.SearchMemories(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to search memories: %w", err)
	}

	switch r := res.(type) {
	case *ogenhosted.SearchMemoriesResponse:
		return convertSearchResults(r.Results), nil
	case *ogenhosted.SearchMemoriesBadRequest:
		return nil, &APIError{StatusCode: http.StatusBadRequest, Detail: r.Detail.Or("")}
	case *ogenhosted.SearchMemoriesUnauthorized:
		return nil, &APIError{StatusCode: http.StatusUnauthorized, Detail: r.Detail.Or("")}
	default:
		return nil, fmt.Errorf("unexpected response type: %T", res)
	}
}

// Update implements Memory.Update.
func (m *hostedMemory) Update(ctx context.Context, memoryID string, text string) (*MemoryItem, error) {
	req := &ogenhosted.UpdateMemoryRequest{
		Text: text,
	}
	params := ogenhosted.UpdateMemoryParams{
		MemoryID: memoryID,
	}

	res, err := m.client.UpdateMemory(ctx, req, params)
	if err != nil {
		return nil, fmt.Errorf("failed to update memory: %w", err)
	}

	switch r := res.(type) {
	case *ogenhosted.MemoryItem:
		return convertMemoryItem(r), nil
	case *ogenhosted.ErrorResponse:
		return nil, &APIError{StatusCode: http.StatusNotFound, Detail: r.Detail.Or("")}
	default:
		return nil, fmt.Errorf("unexpected response type: %T", res)
	}
}

// Delete implements Memory.Delete.
func (m *hostedMemory) Delete(ctx context.Context, memoryID string) error {
	res, err := m.client.DeleteMemory(ctx, ogenhosted.DeleteMemoryParams{MemoryID: memoryID})
	if err != nil {
		return fmt.Errorf("failed to delete memory: %w", err)
	}

	switch r := res.(type) {
	case *ogenhosted.DeleteResponse:
		return nil
	case *ogenhosted.ErrorResponse:
		return &APIError{StatusCode: http.StatusNotFound, Detail: r.Detail.Or("")}
	default:
		return fmt.Errorf("unexpected response type: %T", res)
	}
}

// DeleteAll implements Memory.DeleteAll.
func (m *hostedMemory) DeleteAll(ctx context.Context, opts ...DeleteAllOption) (*DeleteAllResult, error) {
	options := &deleteAllOptions{}
	for _, opt := range opts {
		opt(options)
	}

	params := ogenhosted.DeleteAllMemoriesParams{}
	if options.userID != "" {
		params.UserID = ogenhosted.NewOptString(options.userID)
	}
	if options.agentID != "" {
		params.AgentID = ogenhosted.NewOptString(options.agentID)
	}
	if options.appID != "" {
		params.AppID = ogenhosted.NewOptString(options.appID)
	}
	if options.runID != "" {
		params.RunID = ogenhosted.NewOptString(options.runID)
	}

	res, err := m.client.DeleteAllMemories(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to delete all memories: %w", err)
	}

	return &DeleteAllResult{
		Message:      res.Message.Or(""),
		DeletedCount: res.DeletedCount.Or(0),
	}, nil
}

// History implements Memory.History.
func (m *hostedMemory) History(ctx context.Context, memoryID string) ([]HistoryItem, error) {
	res, err := m.client.GetMemoryHistory(ctx, ogenhosted.GetMemoryHistoryParams{MemoryID: memoryID})
	if err != nil {
		return nil, fmt.Errorf("failed to get memory history: %w", err)
	}

	switch r := res.(type) {
	case *ogenhosted.GetMemoryHistoryOKApplicationJSON:
		return convertHistoryItems(*r), nil
	case *ogenhosted.ErrorResponse:
		return nil, &APIError{StatusCode: http.StatusNotFound, Detail: r.Detail.Or("")}
	default:
		return nil, fmt.Errorf("unexpected response type: %T", res)
	}
}

// GetEventStatus implements Memory.GetEventStatus.
func (m *hostedMemory) GetEventStatus(ctx context.Context, eventID string) (*EventStatus, error) {
	res, err := m.client.GetEventStatus(ctx, ogenhosted.GetEventStatusParams{EventID: eventID})
	if err != nil {
		return nil, fmt.Errorf("failed to get event status: %w", err)
	}

	switch r := res.(type) {
	case *ogenhosted.EventStatus:
		return convertEventStatus(r), nil
	case *ogenhosted.ErrorResponse:
		return nil, &APIError{StatusCode: http.StatusNotFound, Detail: r.Detail.Or("")}
	default:
		return nil, fmt.Errorf("unexpected response type: %T", res)
	}
}

// Conversion helpers

func convertAddMemoryResponse(r *ogenhosted.AddMemoryResponse) *AddResponse {
	resp := &AddResponse{
		EventID: r.EventID.Or(""),
	}

	for _, result := range r.Results {
		resp.Results = append(resp.Results, AddResult{
			ID:     result.ID.Or(""),
			Memory: result.Memory.Or(""),
			Event:  MemoryEvent(result.Event.Or("")),
		})
	}

	return resp
}

func convertMemoryItem(r *ogenhosted.MemoryItem) *MemoryItem {
	item := &MemoryItem{
		ID:         r.ID.Or(""),
		Memory:     r.Memory.Or(""),
		UserID:     r.UserID.Or(""),
		AgentID:    r.AgentID.Or(""),
		AppID:      r.AppID.Or(""),
		RunID:      r.RunID.Or(""),
		Hash:       r.Hash.Or(""),
		Categories: r.Categories,
		CreatedAt:  r.CreatedAt.Or(time.Time{}),
		UpdatedAt:  r.UpdatedAt.Or(time.Time{}),
	}

	if meta, ok := r.Metadata.Get(); ok {
		item.Metadata = make(map[string]interface{})
		for k, v := range meta {
			var decoded interface{}
			if err := json.Unmarshal(v, &decoded); err == nil {
				item.Metadata[k] = decoded
			}
		}
	}

	return item
}

func convertGetMemoriesResponse(r *ogenhosted.GetMemoriesResponse) *GetAllResult {
	result := &GetAllResult{
		Count:    r.Count.Or(0),
		Next:     r.Next.Or(""),
		Previous: r.Previous.Or(""),
	}

	for _, item := range r.Results {
		result.Results = append(result.Results, *convertMemoryItem(&item))
	}

	return result
}

func convertSearchResults(results []ogenhosted.SearchResult) []SearchResult {
	out := make([]SearchResult, len(results))
	for i, r := range results {
		out[i] = SearchResult{
			MemoryItem: MemoryItem{
				ID:         r.ID.Or(""),
				Memory:     r.Memory.Or(""),
				UserID:     r.UserID.Or(""),
				AgentID:    r.AgentID.Or(""),
				AppID:      r.AppID.Or(""),
				RunID:      r.RunID.Or(""),
				Hash:       r.Hash.Or(""),
				Categories: r.Categories,
				CreatedAt:  r.CreatedAt.Or(time.Time{}),
				UpdatedAt:  r.UpdatedAt.Or(time.Time{}),
			},
			Score: float64(r.Score.Or(0)),
		}

		if meta, ok := r.Metadata.Get(); ok {
			out[i].Metadata = make(map[string]interface{})
			for k, v := range meta {
				var decoded interface{}
				if err := json.Unmarshal(v, &decoded); err == nil {
					out[i].Metadata[k] = decoded
				}
			}
		}
	}
	return out
}

func convertHistoryItems(items []ogenhosted.HistoryItem) []HistoryItem {
	out := make([]HistoryItem, len(items))
	for i, item := range items {
		out[i] = HistoryItem{
			ID:        item.ID.Or(""),
			MemoryID:  item.MemoryID.Or(""),
			OldMemory: item.OldMemory.Or(""),
			NewMemory: item.NewMemory.Or(""),
			Event:     MemoryEvent(item.Event.Or("")),
			CreatedAt: item.CreatedAt.Or(time.Time{}),
		}
	}
	return out
}

func convertEventStatus(r *ogenhosted.EventStatus) *EventStatus {
	status := &EventStatus{
		EventID:   r.EventID.Or(""),
		Status:    string(r.Status.Or("")),
		Error:     r.Error.Or(""),
		CreatedAt: r.CreatedAt.Or(time.Time{}),
	}

	if completed, ok := r.CompletedAt.Get(); ok {
		status.CompletedAt = &completed
	}

	for _, result := range r.Results {
		status.Results = append(status.Results, AddResult{
			ID:     result.ID.Or(""),
			Memory: result.Memory.Or(""),
			Event:  MemoryEvent(result.Event.Or("")),
		})
	}

	return status
}

func convertFiltersToOgen(f Filters) ogenhosted.Filters {
	filters := ogenhosted.Filters{}
	if f.UserID != "" {
		filters.UserID = ogenhosted.NewOptString(f.UserID)
	}
	if f.AgentID != "" {
		filters.AgentID = ogenhosted.NewOptString(f.AgentID)
	}
	if f.AppID != "" {
		filters.AppID = ogenhosted.NewOptString(f.AppID)
	}
	if f.RunID != "" {
		filters.RunID = ogenhosted.NewOptString(f.RunID)
	}
	return filters
}
