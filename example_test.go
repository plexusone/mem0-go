package mem0_test

import (
	"context"
	"fmt"
	"log"

	"github.com/plexusone/mem0-go"
)

func Example_basic() {
	// Create a client with API key
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

	fmt.Printf("Added memory with event ID: %s\n", result.EventID)
}

func Example_search() {
	client, err := mem0.NewClient(mem0.WithAPIKey("your-api-key"))
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()

	// Search for memories
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

func Example_getAll() {
	client, err := mem0.NewClient(mem0.WithAPIKey("your-api-key"))
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()

	// Get all memories for a user with pagination
	result, err := client.Memory().GetAll(ctx,
		mem0.WithGetAllFilters(mem0.Filters{UserID: "user-123"}),
		mem0.WithPage(1),
		mem0.WithPageSize(50),
	)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Total memories: %d\n", result.Count)
	for _, m := range result.Results {
		fmt.Printf("- %s\n", m.Memory)
	}
}

func Example_update() {
	client, err := mem0.NewClient(mem0.WithAPIKey("your-api-key"))
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()

	// Update an existing memory
	updated, err := client.Memory().Update(ctx, "memory-id-123", "I now prefer light mode")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Updated memory: %s\n", updated.Memory)
}

func Example_history() {
	client, err := mem0.NewClient(mem0.WithAPIKey("your-api-key"))
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()

	// Get history for a memory
	history, err := client.Memory().History(ctx, "memory-id-123")
	if err != nil {
		log.Fatal(err)
	}

	for _, h := range history {
		fmt.Printf("Event: %s, Old: %s, New: %s\n", h.Event, h.OldMemory, h.NewMemory)
	}
}

func Example_withMetadata() {
	client, err := mem0.NewClient(mem0.WithAPIKey("your-api-key"))
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()

	// Add a memory with metadata
	messages := []mem0.Message{
		{Role: mem0.RoleUser, Content: "My favorite programming language is Go"},
	}

	metadata := map[string]interface{}{
		"source":   "user_profile",
		"verified": true,
	}

	result, err := client.Memory().Add(ctx, messages,
		mem0.WithUserID("user-123"),
		mem0.WithMetadata(metadata),
	)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Added memory with metadata: %v\n", result.EventID)
}

func Example_errorHandling() {
	client, err := mem0.NewClient(mem0.WithAPIKey("your-api-key"))
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()

	// Try to get a non-existent memory
	_, err = client.Memory().Get(ctx, "non-existent-id")
	if err != nil {
		if mem0.IsNotFoundError(err) {
			fmt.Println("Memory not found")
		} else if mem0.IsUnauthorizedError(err) {
			fmt.Println("Authentication failed")
		} else {
			log.Fatal(err)
		}
	}
}
