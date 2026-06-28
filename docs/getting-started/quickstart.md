# Quick Start

This guide walks you through the basic memory operations with mem0-go.

## Creating a Client

```go
import "github.com/plexusone/mem0-go"

// From environment variable (MEM0_API_KEY)
client, err := mem0.NewClient()

// Or with explicit API key
client, err := mem0.NewClient(mem0.WithAPIKey("your-api-key"))
```

## Adding Memories

Add memories from conversation messages:

```go
messages := []mem0.Message{
    {Role: mem0.RoleUser, Content: "My favorite color is blue"},
    {Role: mem0.RoleAssistant, Content: "I'll remember that your favorite color is blue!"},
}

result, err := client.Memory().Add(ctx, messages,
    mem0.WithUserID("user-123"),
    mem0.WithMetadata(map[string]any{"source": "onboarding"}),
)
```

## Searching Memories

Find relevant memories using semantic search:

```go
results, err := client.Memory().Search(ctx, "color preferences",
    mem0.WithFilters(mem0.Filters{UserID: "user-123"}),
    mem0.WithTopK(5),
)

for _, r := range results {
    fmt.Printf("Memory: %s (score: %.2f)\n", r.Memory, r.Score)
}
```

## Getting All Memories

Retrieve all memories for a user:

```go
memories, err := client.Memory().GetAll(ctx,
    mem0.WithUserID("user-123"),
    mem0.WithPage(1),
    mem0.WithPageSize(10),
)

for _, m := range memories {
    fmt.Printf("ID: %s, Memory: %s\n", m.ID, m.Memory)
}
```

## Updating a Memory

```go
updated, err := client.Memory().Update(ctx, "memory-id", "Updated memory content")
```

## Deleting Memories

```go
// Delete a single memory
err := client.Memory().Delete(ctx, "memory-id")

// Delete all memories for a user
err := client.Memory().DeleteAll(ctx, mem0.WithUserID("user-123"))
```

## Memory History

Track changes to a memory over time:

```go
history, err := client.Memory().History(ctx, "memory-id")

for _, h := range history {
    fmt.Printf("Version %d: %s -> %s\n", h.ID, h.PrevValue, h.NewValue)
}
```

## Multi-Tenancy

Organize memories by different scopes:

```go
// User-level memories
mem0.WithUserID("user-123")

// Agent-level memories
mem0.WithAgentID("agent-456")

// App-level memories
mem0.WithAppID("app-789")

// Run/session-level memories
mem0.WithRunID("run-abc")

// Combine multiple scopes
client.Memory().Add(ctx, messages,
    mem0.WithUserID("user-123"),
    mem0.WithAgentID("support-agent"),
)
```

## Error Handling

```go
result, err := client.Memory().Get(ctx, "memory-id")
if err != nil {
    if mem0.IsNotFoundError(err) {
        fmt.Println("Memory not found")
    } else if mem0.IsAuthenticationError(err) {
        fmt.Println("Invalid API key")
    } else {
        fmt.Printf("Error: %v\n", err)
    }
}
```

## Next Steps

- [OmniMemory Provider](../omnimemory.md) - Use with vendor-neutral abstraction
- [API Reference](https://pkg.go.dev/github.com/plexusone/mem0-go) - Full API documentation
