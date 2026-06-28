package mem0

import "context"

// Memory is the core interface for memory operations.
type Memory interface {
	// Add extracts facts from conversation messages and stores them as memories.
	Add(ctx context.Context, messages []Message, opts ...AddOption) (*AddResponse, error)

	// Get retrieves a specific memory by ID.
	Get(ctx context.Context, memoryID string) (*MemoryItem, error)

	// GetAll retrieves all memories with optional filtering and pagination.
	GetAll(ctx context.Context, opts ...GetAllOption) (*GetAllResult, error)

	// Search performs a hybrid search combining semantic, BM25, and entity matching.
	Search(ctx context.Context, query string, opts ...SearchOption) ([]SearchResult, error)

	// Update modifies an existing memory.
	Update(ctx context.Context, memoryID string, text string) (*MemoryItem, error)

	// Delete removes a single memory.
	Delete(ctx context.Context, memoryID string) error

	// DeleteAll removes all memories matching the given filters.
	DeleteAll(ctx context.Context, opts ...DeleteAllOption) (*DeleteAllResult, error)

	// History retrieves the change history for a memory.
	History(ctx context.Context, memoryID string) ([]HistoryItem, error)

	// GetEventStatus retrieves the status of an asynchronous operation.
	GetEventStatus(ctx context.Context, eventID string) (*EventStatus, error)
}
