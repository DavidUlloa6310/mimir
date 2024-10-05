package database

import (
	"context"
	"fmt"
	"log"

	"github.com/davidulloa/mimir/models"
	"github.com/google/uuid"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/fault"
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
        log.Printf("Expected properties to be a map but got: %v", result[0].Properties)
        return nil, fmt.Errorf("invalid properties for accelerator with ID: %s", acceleratorID)
    }

    if id, ok := properties["_additional"].(map[string]interface{})["id"].(float64); ok {
        accelerator.ID = int(id)
    } else {
        log.Printf("Missing or invalid ID for accelerator with ID: %s", acceleratorID)
        return nil, fmt.Errorf("invalid ID for accelerator with ID: %s", acceleratorID)
    }

    if url, ok := properties["url"].(string); ok {
        accelerator.Url = url
    }
    if title, ok := properties["title"].(string); ok {
        accelerator.Title = title
    }
    if description, ok := properties["description"].(string); ok {
        accelerator.Description = description
    }
    if category, ok := properties["category"].(string); ok {
        accelerator.Category = category
    }

    log.Printf("Retrieved accelerator with ID: %s", acceleratorID)
    return accelerator, nil
}

func CreateAccelerator(accelerator *models.Accelerator) (string, error) {
    client, err := GetWeaviateClient()
    if err != nil {
        log.Printf("Error getting Weaviate client: %v", err)
        return "", err
    }

    response, err := client.Data().Creator().
        WithClassName("Accelerator").
        WithProperties(map[string]interface{}{
            "url":         accelerator.Url,
            "title":       accelerator.Title,
            "description": accelerator.Description,
            "category":    accelerator.Category,
			"iD": 		   uuid.New().String(),
        }).
        Do(context.Background())

    if err != nil {
        log.Printf("Error creating accelerator: %v", err)
        return "", err
    }

	acceleratorID := response.Object.ID
    log.Printf("Created accelerator with ID: %s", acceleratorID)
    return string(acceleratorID, nil
}

func DeleteAccelerator(acceleratorID string) error {
    client, err := GetWeaviateClient()
    if err != nil {
        log.Printf("Error getting Weaviate client: %v", err)
        return err
    }

    err = client.Data().Deleter().
        WithClassName("Accelerator").
        WithID(acceleratorID).
        Do(context.Background())

    if err != nil {
        if clientErr, ok := err.(*fault.WeaviateClientError); ok {
            if clientErr.StatusCode == 404 {
                log.Printf("Accelerator not found with ID: %s", acceleratorID)
                return fmt.Errorf("accelerator not found")
            }
        }
        log.Printf("Error deleting accelerator with ID %s: %v", acceleratorID, err)
        return err
    }

    log.Printf("Deleted accelerator with ID: %s", acceleratorID)
    return nil
}