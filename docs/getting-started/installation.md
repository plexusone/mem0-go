# Installation

## Requirements

- Go 1.21 or later
- A mem0 API key from [mem0.ai](https://mem0.ai)

## Install

```bash
go get github.com/plexusone/mem0-go
```

## API Key

Get your API key from [mem0.ai](https://mem0.ai) and set it as an environment variable:

```bash
export MEM0_API_KEY="your-api-key"
```

Or pass it directly when creating the client:

```go
client, err := mem0.NewClient(mem0.WithAPIKey("your-api-key"))
```

## Verify Installation

```go
package main

import (
    "fmt"
    "github.com/plexusone/mem0-go"
)

func main() {
    client, err := mem0.NewClient()
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }
    fmt.Println("mem0-go installed successfully!")
    _ = client
}
```

## Next Steps

- [Quick Start](quickstart.md) - Your first memory operation
