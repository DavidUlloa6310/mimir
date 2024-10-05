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
    err := godotenv.Load("../.env")
    if err != nil {
        log.Fatal("Error loading .env variables")
    }

	if weaviateClient != nil {
		return weaviateClient, nil
	}

	cfg := weaviate.Config{
		Host: os.Getenv("WEAVIATE_URL"),
		Scheme: "https",
		AuthConfig: auth.ApiKey{
			Value: os.Getenv("WEAVIATE_API_KEY"),
		},
		Headers: map[string]string{
			"X-OpenAI-Api-Key": os.Getenv("OPENAI_API_KEY"),
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

