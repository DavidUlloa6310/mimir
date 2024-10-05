package database

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/weaviate/weaviate-go-client/v4/weaviate"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/auth"
)

var weaviateClient *weaviate.Client

func InitWeaviateClient() (*weaviate.Client, error) {
    if weaviateClient != nil {
        return weaviateClient, nil
    }

    // Load .env file
    err := godotenv.Load("../../.env")
    if err != nil {
        log.Printf("Warning: Error loading .env file: %v", err)
    }

    weaviateURL := os.Getenv("WEAVIATE_URL")
    weaviateAPIKey := os.Getenv("WEAVIATE_API_KEY")
    openAIAPIKey := os.Getenv("OPENAI_API_KEY")

    if weaviateURL == "" || weaviateAPIKey == "" || openAIAPIKey == "" {
        return nil, fmt.Errorf("missing required environment variables")
    }

    cfg := weaviate.Config{
        Host: weaviateURL,
        Scheme: "https",
        AuthConfig: auth.ApiKey{
            Value: weaviateAPIKey,
        },
        Headers: map[string]string{
            "X-OpenAI-Api-Key": openAIAPIKey,
        },
    }

    client, err := weaviate.NewClient(cfg)
    if err != nil {
        return nil, fmt.Errorf("error creating Weaviate client: %v", err)
    }

    weaviateClient = client
    return client, nil
}

func GetWeaviateClient() (*weaviate.Client, error) {
    if weaviateClient == nil {
        return InitWeaviateClient()
    }
    return weaviateClient, nil
}