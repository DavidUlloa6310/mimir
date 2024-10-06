package database

import (
	"context"
	"fmt"
	"log"

	"github.com/davidulloa/mimir/models"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/fault"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/graphql"
)

func GetAcceleratorByID(acceleratorID string) (*models.Accelerator, error) {

    client, err := GetWeaviateClient()
    if err != nil {
        log.Printf("Error getting Weaviate client: %v", err)
        return nil, err
    }

    result, err := client.Data().ObjectsGetter().
        WithClassName("Accelerator").
        WithID(acceleratorID).
        Do(context.Background())

    if err != nil {
        if clientErr, ok := err.(*fault.WeaviateClientError); ok {
            if clientErr.StatusCode == 404 {
                log.Printf("Accelerator not found with ID: %s", acceleratorID)
                return nil, fmt.Errorf("accelerator not found")
            }
        }
        log.Printf("Error retrieving accelerator with ID %s: %v", acceleratorID, err)
        return nil, err
    }

    if len(result) == 0 {
        log.Printf("Accelerator not found with ID: %s", acceleratorID)
        return nil, fmt.Errorf("accelerator not found")
    }

    accelerator := &models.Accelerator{}
    properties, ok := result[0].Properties.(map[string]interface{})
    if !ok {
        log.Printf("Expected properties to be a map but got: %T %+v", result[0].Properties, result[0].Properties)
        return nil, fmt.Errorf("invalid properties for accelerator with ID: %s", acceleratorID)
    }

    if url, ok := properties["url"].(string); ok {
        accelerator.Url = url
    } else {
        log.Printf("Missing or invalid URL for ID %s", acceleratorID)
    }

    if title, ok := properties["title"].(string); ok {
        accelerator.Title = title
    } else {
        log.Printf("Missing or invalid Title for ID %s", acceleratorID)
    }

    if description, ok := properties["description"].(string); ok {
        accelerator.Description = description
    } else {
        log.Printf("Missing or invalid Description for ID %s", acceleratorID)
    }

    if category, ok := properties["category"].(string); ok {
        accelerator.Category = category
    } else {
        log.Printf("Missing or invalid Category for ID %s", acceleratorID)
    }

    return accelerator, nil
}

func GetAllAccelerators() ([]string, error) {
    client, err := GetWeaviateClient()
    if err != nil {
        log.Printf("Error getting Weaviate client: %v", err)
        return nil, err
    }

    fields := []graphql.Field{
        {Name: "_additional { id }"},
    }

    result, err := client.GraphQL().Get().
        WithClassName("Accelerator").
        WithFields(fields...).
        Do(context.Background())

    if err != nil {
        log.Printf("Error retrieving all accelerators: %v", err)
        return nil, err
    }

    if result.Errors != nil {
        for _, err := range result.Errors {
            log.Printf("GraphQL error: %v", err)
        }
        return nil, fmt.Errorf("graphQL errors: %v", result.Errors)
    }

    var acceleratorIDs []string
    accelerators, ok := result.Data["Get"].(map[string]interface{})["Accelerator"]
    if !ok {
        log.Printf("Unexpected data structure: %v", result.Data)
        return nil, fmt.Errorf("unexpected data structure: %v", result.Data)
    }

    acceleratorsSlice, ok := accelerators.([]interface{})
    if !ok {
        log.Printf("Expected slice of accelerators but got: %v", accelerators)
        return nil, fmt.Errorf("expected slice of accelerators but got: %v", accelerators)
    }

    for _, obj := range acceleratorsSlice {
        acc, ok := obj.(map[string]interface{})
        if !ok {
            log.Printf("Expected accelerator object but got: %v", obj)
            return nil, fmt.Errorf("expected accelerator object but got: %v", obj)
        }

        additional, ok := acc["_additional"].(map[string]interface{})
        if !ok {
            log.Printf("Missing _additional field for accelerator: %v", acc)
            return nil, fmt.Errorf("missing _additional field for accelerator: %v", acc)
        }

        id, ok := additional["id"].(string)
        if !ok {
            log.Printf("Invalid or missing ID in _additional field: %v", additional)
            return nil, fmt.Errorf("invalid or missing ID in _additional field: %v", additional)
        }

        acceleratorIDs = append(acceleratorIDs, id)
    }

    return acceleratorIDs, nil
}