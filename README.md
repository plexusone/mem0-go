# mem0-go

[![Go CI][go-ci-svg]][go-ci-url]
[![Go Lint][go-lint-svg]][go-lint-url]
[![Go SAST][go-sast-svg]][go-sast-url]
[![Go Report Card][goreport-svg]][goreport-url]
[![Docs][docs-godoc-svg]][docs-godoc-url]
[![Docs][docs-mkdoc-svg]][docs-mkdoc-url]
[![Visualization][viz-svg]][viz-url]
[![License][license-svg]][license-url]

 [go-ci-svg]: https://github.com/plexusone/mem0-go/actions/workflows/go-ci.yaml/badge.svg?branch=main
 [go-ci-url]: https://github.com/plexusone/mem0-go/actions/workflows/go-ci.yaml
 [go-lint-svg]: https://github.com/plexusone/mem0-go/actions/workflows/go-lint.yaml/badge.svg?branch=main
 [go-lint-url]: https://github.com/plexusone/mem0-go/actions/workflows/go-lint.yaml
 [go-sast-svg]: https://github.com/plexusone/mem0-go/actions/workflows/go-sast-codeql.yaml/badge.svg?branch=main
 [go-sast-url]: https://github.com/plexusone/mem0-go/actions/workflows/go-sast-codeql.yaml
 [goreport-svg]: https://goreportcard.com/badge/github.com/plexusone/mem0-go
 [goreport-url]: https://goreportcard.com/report/github.com/plexusone/mem0-go
 [docs-godoc-svg]: https://pkg.go.dev/badge/github.com/plexusone/mem0-go
 [docs-godoc-url]: https://pkg.go.dev/github.com/plexusone/mem0-go
 [docs-mkdoc-svg]: https://img.shields.io/badge/Go-dev%20guide-blue.svg
 [docs-mkdoc-url]: https://plexusone.dev/mem0-go
 [viz-svg]: https://img.shields.io/badge/Go-visualizaton-blue.svg
 [viz-url]: https://mango-dune-07a8b7110.1.azurestaticapps.net/?repo=plexusone%2Fmem0-go
 [loc-svg]: https://tokei.rs/b1/github/plexusone/mem0-go
 [repo-url]: https://github.com/plexusone/mem0-go
 [license-svg]: https://img.shields.io/badge/license-MIT-blue.svg
 [license-url]: https://github.com/plexusone/mem0-go/blob/main/LICENSE

Go client SDK for [mem0](https://mem0.ai) - the AI memory layer for agents and apps.

## Installation

```bash
go get github.com/plexusone/mem0-go
```

## Quick Start

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
    fmt.Printf("Added memory: %v\n", result)

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

## Configuration

### Environment Variables

| Variable | Description |
|----------|-------------|
| `MEM0_API_KEY` | API key for authentication |
| `MEM0_BASE_URL` | Custom base URL (default: `https://api.mem0.ai`) |

### Client Options

```go
client, err := mem0.NewClient(
    mem0.WithAPIKey("your-api-key"),
    mem0.WithBaseURL("https://custom.example.com"),
    mem0.WithHTTPClient(customHTTPClient),
    mem0.WithBackend(mem0.BackendHosted),
)
```

## API Reference

### Memory Interface

```go
type Memory interface {
    Add(ctx context.Context, messages []Message, opts ...AddOption) (*AddResponse, error)
    Get(ctx context.Context, memoryID string) (*MemoryItem, error)
    GetAll(ctx context.Context, opts ...GetAllOption) (*GetAllResult, error)
    Search(ctx context.Context, query string, opts ...SearchOption) ([]SearchResult, error)
    Update(ctx context.Context, memoryID string, text string) (*MemoryItem, error)
    Delete(ctx context.Context, memoryID string) error
    DeleteAll(ctx context.Context, opts ...DeleteAllOption) (*DeleteAllResult, error)
    History(ctx context.Context, memoryID string) ([]HistoryItem, error)
    GetEventStatus(ctx context.Context, eventID string) (*EventStatus, error)
}
```

### Add Options

```go
mem0.WithUserID("user-123")
mem0.WithAgentID("agent-456")
mem0.WithAppID("app-789")
mem0.WithRunID("run-012")
mem0.WithMetadata(map[string]interface{}{"key": "value"})
mem0.WithInfer(true)
```

### Search Options

```go
mem0.WithFilters(mem0.Filters{UserID: "user-123"})
mem0.WithTopK(10)
mem0.WithRerank(true)
mem0.WithThreshold(0.5)
```

### GetAll Options

```go
mem0.WithGetAllFilters(mem0.Filters{UserID: "user-123"})
mem0.WithPage(1)
mem0.WithPageSize(50)
```

### DeleteAll Options

```go
mem0.WithDeleteUserID("user-123")
mem0.WithDeleteAgentID("agent-456")
```

## Error Handling

```go
_, err := client.Memory().Get(ctx, "non-existent-id")
if err != nil {
    if mem0.IsNotFoundError(err) {
        fmt.Println("Memory not found")
    } else if mem0.IsUnauthorizedError(err) {
        fmt.Println("Authentication failed")
    } else if mem0.IsBadRequestError(err) {
        fmt.Println("Invalid request")
    } else if mem0.IsRateLimitedError(err) {
        fmt.Println("Rate limit exceeded")
    }
}
```

## Code Generation

This SDK uses [ogen](https://github.com/ogen-go/ogen) for OpenAPI code generation.

To regenerate the client:

```bash
./generate.sh
```

## License

MIT
