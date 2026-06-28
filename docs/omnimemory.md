# OmniMemory Provider

mem0-go includes an [OmniMemory](https://github.com/plexusone/omnimemory) provider, enabling vendor-neutral memory operations with automatic fallback support.

## Overview

OmniMemory is a unified memory abstraction layer for Go that provides:

- Consistent interface across multiple memory backends
- Automatic failover between providers
- Multi-tenancy with tenant/subject isolation
- Semantic search capabilities

## Installation

```go
import (
    "github.com/plexusone/omnimemory/core"
    _ "github.com/plexusone/mem0-go/omnimemory" // Register mem0 provider
)
```

## Configuration

The mem0 provider is registered automatically when imported. Configure via environment or options:

```go
// Environment variables
// MEM0_API_KEY - API key for mem0.ai

// Or via ProviderConfig
config := core.ProviderConfig{
    APIKey: "your-api-key",
    Options: map[string]any{
        "base_url": "https://api.mem0.ai/v1", // Optional custom URL
    },
}
```

## Usage with OmniMemory

```go
package main

import (
    "context"
    "log"

    "github.com/plexusone/omnimemory/core"
    _ "github.com/plexusone/mem0-go/omnimemory"
)

func main() {
    // Create provider from registry
    provider, err := core.NewProvider("mem0", core.ProviderConfig{
        APIKey: "your-api-key",
    }, nil)
    if err != nil {
        log.Fatal(err)
    }
    defer provider.Close()

    ctx := context.Background()

    // Store a memory
    memory, err := provider.Store(ctx, core.StoreRequest{
        TenantID:  "app-123",
        SubjectID: "user-456",
        Content:   "User prefers dark mode",
        Type:      core.MemoryTypePreference,
        Scope:     core.ScopeUser,
    })
    if err != nil {
        log.Fatal(err)
    }

    // Search memories
    results, err := provider.Search(ctx, core.SearchRequest{
        TenantID:  "app-123",
        SubjectID: "user-456",
        Query:     "theme preferences",
        Limit:     10,
    })
    if err != nil {
        log.Fatal(err)
    }

    for _, r := range results {
        log.Printf("Memory: %s (score: %.2f)", r.Content, r.Score)
    }
}
```

## Tenant/Subject Mapping

The mem0 provider maps OmniMemory's multi-tenancy model:

| OmniMemory | mem0 |
|------------|------|
| TenantID | app_id |
| SubjectID | user_id |

## Memory Type and Scope

Memory type and scope are stored in mem0 metadata with prefixed keys:

- `omnimemory_type` - Memory type (fact, preference, observation, etc.)
- `omnimemory_scope` - Memory scope (user, agent, tenant)

## Provider Priority

The mem0 provider is registered with `PriorityThick` priority, indicating it uses an external service (as opposed to in-process providers).

## Next Steps

- [OmniMemory Documentation](https://github.com/plexusone/omnimemory)
- [API Reference](https://pkg.go.dev/github.com/plexusone/mem0-go/omnimemory)
