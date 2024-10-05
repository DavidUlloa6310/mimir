package database

import (
	"testing"
)

func TestWeaviateOperations(t *testing.T) {
	// Step 1: Initialize Weaviate Client
	_, err := InitWeaviateClient()
	if err != nil {
		t.Fatalf("Failed to initialize Weaviate client: %v", err)
	}
}