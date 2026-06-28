package mem0

import "time"

// Role represents the role of a message sender.
type Role string

const (
	RoleUser      Role = "user"
	RoleAssistant Role = "assistant"
	RoleSystem    Role = "system"
)

// Message represents a conversation message.
type Message struct {
	Role    Role   `json:"role"`
	Content string `json:"content"`
}

// MemoryItem represents a single memory record.
type MemoryItem struct {
	ID         string                 `json:"id"`
	Memory     string                 `json:"memory"`
	UserID     string                 `json:"user_id,omitempty"`
	AgentID    string                 `json:"agent_id,omitempty"`
	AppID      string                 `json:"app_id,omitempty"`
	RunID      string                 `json:"run_id,omitempty"`
	Hash       string                 `json:"hash,omitempty"`
	Categories []string               `json:"categories,omitempty"`
	Metadata   map[string]interface{} `json:"metadata,omitempty"`
	CreatedAt  time.Time              `json:"created_at"`
	UpdatedAt  time.Time              `json:"updated_at"`
}

// SearchResult represents a memory search result with a relevance score.
type SearchResult struct {
	MemoryItem
	Score float64 `json:"score"`
}

// MemoryEvent represents the type of operation performed on a memory.
type MemoryEvent string

const (
	MemoryEventAdd    MemoryEvent = "ADD"
	MemoryEventUpdate MemoryEvent = "UPDATE"
	MemoryEventDelete MemoryEvent = "DELETE"
	MemoryEventNoop   MemoryEvent = "NOOP"
)

// AddResult represents the result of adding a memory.
type AddResult struct {
	ID      string      `json:"id"`
	Memory  string      `json:"memory"`
	Event   MemoryEvent `json:"event"`
	EventID string      `json:"event_id,omitempty"`
}

// AddResponse represents the full response from adding memories.
type AddResponse struct {
	EventID string      `json:"event_id,omitempty"`
	Results []AddResult `json:"results,omitempty"`
}

// HistoryItem represents a single history entry for a memory.
type HistoryItem struct {
	ID        string      `json:"id"`
	MemoryID  string      `json:"memory_id"`
	OldMemory string      `json:"old_memory,omitempty"`
	NewMemory string      `json:"new_memory"`
	Event     MemoryEvent `json:"event"`
	CreatedAt time.Time   `json:"created_at"`
}

// EventStatus represents the status of an asynchronous operation.
type EventStatus struct {
	EventID     string      `json:"event_id"`
	Status      string      `json:"status"`
	Results     []AddResult `json:"results,omitempty"`
	Error       string      `json:"error,omitempty"`
	CreatedAt   time.Time   `json:"created_at"`
	CompletedAt *time.Time  `json:"completed_at,omitempty"`
}

// EntityType represents the type of entity.
type EntityType string

const (
	EntityTypeUser  EntityType = "user"
	EntityTypeAgent EntityType = "agent"
	EntityTypeApp   EntityType = "app"
	EntityTypeRun   EntityType = "run"
)

// Entity represents an entity (user, agent, app, or run).
type Entity struct {
	Type        EntityType `json:"type"`
	ID          string     `json:"id"`
	MemoryCount int        `json:"memory_count"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// Filters represents filter conditions for memory queries.
type Filters struct {
	UserID  string `json:"user_id,omitempty"`
	AgentID string `json:"agent_id,omitempty"`
	AppID   string `json:"app_id,omitempty"`
	RunID   string `json:"run_id,omitempty"`
}

// GetAllResult represents the result of getting all memories with pagination.
type GetAllResult struct {
	Count    int          `json:"count"`
	Next     string       `json:"next,omitempty"`
	Previous string       `json:"previous,omitempty"`
	Results  []MemoryItem `json:"results"`
}

// DeleteAllResult represents the result of deleting all memories.
type DeleteAllResult struct {
	Message      string `json:"message"`
	DeletedCount int    `json:"deleted_count"`
}
