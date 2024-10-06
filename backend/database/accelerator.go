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

func GetAllAccelerators() ([]models.Accelerator, error) {
    client, err := GetWeaviateClient()
    if err != nil {
        log.Printf("Error getting Weaviate client: %v", err)
        return nil, err
    }

    fields := []graphql.Field{
        {Name: "url"},
        {Name: "title"},
        {Name: "description"},
        {Name: "category"},
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

    var acceleratorsList []models.Accelerator

    getData, ok := result.Data["Get"].(map[string]interface{})
    if !ok {
        log.Printf("Unexpected data structure: %v", result.Data)
        return nil, fmt.Errorf("unexpected data structure: %v", result.Data)
    }

    accelerators, ok := getData["Accelerator"]
    if !ok {
        log.Printf("No 'Accelerator' key in data: %v", getData)
        return nil, fmt.Errorf("no 'Accelerator' key in data: %v", getData)
    }

    acceleratorsSlice, ok := accelerators.([]interface{})
    if !ok {
        log.Printf("Expected slice of accelerators but got: %v", accelerators)
        return nil, fmt.Errorf("expected slice of accelerators but got: %v", accelerators)
    }

    for _, obj := range acceleratorsSlice {
        accMap, ok := obj.(map[string]interface{})
        if !ok {
            log.Printf("Expected accelerator object but got: %v", obj)
            return nil, fmt.Errorf("expected accelerator object but got: %v", obj)
        }

        var accelerator models.Accelerator

        if url, ok := accMap["url"].(string); ok {
            accelerator.Url = url
        }

        if title, ok := accMap["title"].(string); ok {
            accelerator.Title = title
        }

        if description, ok := accMap["description"].(string); ok {
            accelerator.Description = description
        }

        if category, ok := accMap["category"].(string); ok {
            accelerator.Category = category
        }

        // if additional, ok := accMap["_additional"].(map[string]interface{}); ok {
        //     if id, ok := additional["id"].(string); ok {
        //         // Optionally, you can add an ID field to your Accelerator struct to store this value
        //         accelerator.Id = id
        //     }
        // }

        acceleratorsList = append(acceleratorsList, accelerator)
    }

    return acceleratorsList, nil
}
