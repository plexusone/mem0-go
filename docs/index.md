# mem0-go

**Go client SDK for [mem0](https://mem0.ai) - the AI memory layer for agents and apps.**

mem0-go provides a type-safe Go interface for the mem0 API, enabling persistent memory for AI applications. Store user preferences, facts, and context that persists across sessions.

## Features

- **Memory Operations**: Add, search, update, and delete memories
- **Multi-Tenancy**: Organize memories by user, agent, app, or run
- **Semantic Search**: Find relevant memories using natural language queries
- **Memory History**: Track changes to memories over time
- **Type Safe**: Full Go type safety with comprehensive error handling
- **OmniMemory Provider**: Integrates with [OmniMemory](https://github.com/plexusone/omnimemory) abstraction layer

## Quick Example

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/plexusone/mem0-go"
)

func main() {
    // Create a client (API key from MEM0_API_KEY env var or WithAPIKey option)
    client, err := mem0.NewClient(mem0.WithAPIKey("your-api-key"))
    if err != nil {
        log.Fatal(err)
    }

    ctx := context.Background()

    // Add a memory
    messages := []mem0.Message{
        {Role: mem0.RoleUser, Content: "I prefer dark mode in all my applications"},
    }

    result, err := client.Memory().Add(ctx, messages, mem0.WithUserID("user-123"))
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Added memory: %s\n", result.ID)

    // Search memories
    results, err := client.Memory().Search(ctx, "dark mode preferences",
        mem0.WithFilters(mem0.Filters{UserID: "user-123"}),
        mem0.WithTopK(5),
    )
    if err != nil {
        log.Fatal(err)
    }

    for _, r := range results {
        fmt.Printf("Memory: %s (score: %.2f)\n", r.Memory, r.Score)
    }
}
```

## Next Steps

- [Installation](getting-started/installation.md) - Get mem0-go set up
- [Quick Start](getting-started/quickstart.md) - Your first memory operation
- [OmniMemory Provider](omnimemory.md) - Use with vendor-neutral abstraction
- [API Reference](https://pkg.go.dev/github.com/plexusone/mem0-go) - Full API documentation
