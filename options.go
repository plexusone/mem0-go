package mem0

import (
	"net/http"
)

// Backend selects which mem0 deployment to use.
type Backend int

const (
	// BackendHosted uses the hosted mem0 platform at api.mem0.ai.
	BackendHosted Backend = iota
	// BackendFOSS uses a self-hosted mem0 deployment.
	BackendFOSS
)

// Option configures the Client.
type Option func(*clientOptions)

type clientOptions struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client
	backend    Backend
}

// WithAPIKey sets the API key for authentication.
func WithAPIKey(apiKey string) Option {
	return func(o *clientOptions) {
		o.apiKey = apiKey
	}
}

// WithBaseURL sets a custom base URL for the API.
func WithBaseURL(baseURL string) Option {
	return func(o *clientOptions) {
		o.baseURL = baseURL
	}
}

// WithHTTPClient sets a custom HTTP client.
func WithHTTPClient(client *http.Client) Option {
	return func(o *clientOptions) {
		o.httpClient = client
	}
}

// WithBackend sets the backend to use (hosted or FOSS).
func WithBackend(backend Backend) Option {
	return func(o *clientOptions) {
		o.backend = backend
	}
}

// AddOption configures the Add operation.
type AddOption func(*addOptions)

type addOptions struct {
	userID   string
	agentID  string
	appID    string
	runID    string
	metadata map[string]interface{}
	infer    *bool
}

// WithUserID sets the user ID for the operation.
func WithUserID(userID string) AddOption {
	return func(o *addOptions) {
		o.userID = userID
	}
}

// WithAgentID sets the agent ID for the operation.
func WithAgentID(agentID string) AddOption {
	return func(o *addOptions) {
		o.agentID = agentID
	}
}

// WithAppID sets the app ID for the operation.
func WithAppID(appID string) AddOption {
	return func(o *addOptions) {
		o.appID = appID
	}
}

// WithRunID sets the run ID for the operation.
func WithRunID(runID string) AddOption {
	return func(o *addOptions) {
		o.runID = runID
	}
}

// WithMetadata sets metadata for the operation.
func WithMetadata(metadata map[string]interface{}) AddOption {
	return func(o *addOptions) {
		o.metadata = metadata
	}
}

// WithInfer sets whether to infer facts from messages.
func WithInfer(infer bool) AddOption {
	return func(o *addOptions) {
		o.infer = &infer
	}
}

// SearchOption configures the Search operation.
type SearchOption func(*searchOptions)

type searchOptions struct {
	filters   *Filters
	topK      int
	rerank    bool
	threshold *float64
}

// WithFilters sets filters for the search.
func WithFilters(filters Filters) SearchOption {
	return func(o *searchOptions) {
		o.filters = &filters
	}
}

// WithTopK sets the number of results to return.
func WithTopK(topK int) SearchOption {
	return func(o *searchOptions) {
		o.topK = topK
	}
}

// WithRerank enables reranking of search results.
func WithRerank(rerank bool) SearchOption {
	return func(o *searchOptions) {
		o.rerank = rerank
	}
}

// WithThreshold sets the minimum score threshold for results.
func WithThreshold(threshold float64) SearchOption {
	return func(o *searchOptions) {
		o.threshold = &threshold
	}
}

// GetAllOption configures the GetAll operation.
type GetAllOption func(*getAllOptions)

type getAllOptions struct {
	filters  *Filters
	page     int
	pageSize int
}

// WithGetAllFilters sets filters for getting all memories.
func WithGetAllFilters(filters Filters) GetAllOption {
	return func(o *getAllOptions) {
		o.filters = &filters
	}
}

// WithPage sets the page number for pagination.
func WithPage(page int) GetAllOption {
	return func(o *getAllOptions) {
		o.page = page
	}
}

// WithPageSize sets the page size for pagination.
func WithPageSize(pageSize int) GetAllOption {
	return func(o *getAllOptions) {
		o.pageSize = pageSize
	}
}

// DeleteAllOption configures the DeleteAll operation.
type DeleteAllOption func(*deleteAllOptions)

type deleteAllOptions struct {
	userID  string
	agentID string
	appID   string
	runID   string
}

// WithDeleteUserID sets the user ID for delete all.
func WithDeleteUserID(userID string) DeleteAllOption {
	return func(o *deleteAllOptions) {
		o.userID = userID
	}
}

// WithDeleteAgentID sets the agent ID for delete all.
func WithDeleteAgentID(agentID string) DeleteAllOption {
	return func(o *deleteAllOptions) {
		o.agentID = agentID
	}
}

// WithDeleteAppID sets the app ID for delete all.
func WithDeleteAppID(appID string) DeleteAllOption {
	return func(o *deleteAllOptions) {
		o.appID = appID
	}
}

// WithDeleteRunID sets the run ID for delete all.
func WithDeleteRunID(runID string) DeleteAllOption {
	return func(o *deleteAllOptions) {
		o.runID = runID
	}
}
