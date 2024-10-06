package database

import (
	"context"
	"errors"

	"github.com/davidulloa/mimir/models"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/filters"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/graphql"
	WeaviateModels "github.com/weaviate/weaviate/entities/models"
)

const (
	TicketClass = "Ticket"
)

func RetrieveTickets(ids []string) ([]models.Ticket, error) {
	client, err := GetWeaviateClient()
	if err != nil {
		return []models.Ticket{}, err
	}

	fields := []string{"_additional { id }", "shortDescription", "state", "priority", "number"}
	graphqlFields := make([]graphql.Field, len(fields))
	for i, field := range fields {
		graphqlFields[i] = graphql.Field{Name: field}
	}

	inIds := filters.Where().
		WithPath([]string{"_additional { id }"}).
		WithOperator(filters.ContainsAny).
		WithValueText(ids...)

	response, err := client.GraphQL().Get().WithClassName(TicketClass).WithFields(graphqlFields...).
		WithWhere(inIds).
		Do(context.Background())
	
	if err != nil {
		return []models.Ticket{}, err
	}

	getObject, ok := response.Data["Get"].(map[string]interface{})
	if !ok {
		return []models.Ticket{}, errors.New("unable to parse 'Get' from response data")
	}

	classObjects, ok := getObject[AuthorizationClass].([]interface{})
	if !ok {
		return []models.Ticket{}, errors.New("unable to parse 'Authorization' class from class object")
	}

	tickets := []models.Ticket{}

    for _, obj := range classObjects {
        objMap, ok := obj.(map[string]interface{})
        if !ok {
            return []models.Ticket{}, errors.New("unable to parse ticket object")
        }

        var ticket models.Ticket

        // Map the fields from objMap to the ticket struct
        if id, ok := objMap["id"].(string); ok {
            ticket.ID = id
        }

        if priority, ok := objMap["priority"].(int); ok {
            ticket.Priority = string(priority)
        }

        if shortDescription, ok := objMap["shortDescription"].(string); ok {
            ticket.ShortDescription = shortDescription
        }

        if state, ok := objMap["state"].(int); ok {
            ticket.State = string(state)
		}

		if number, ok := objMap["number"].(int); ok {
			ticket.Number = string(number)
        }

        tickets = append(tickets, ticket)
    }

	return tickets, nil
	
}

func StoreTickets(tickets []models.Ticket) error {
	client, err := GetWeaviateClient()
	if err != nil {
		return err
	}

	batcher := client.Batch().ObjectsBatcher()
	dataObjs := []WeaviateModels.PropertySchema{}
	for _, ticket := range tickets {
		dataObjs = append(dataObjs, map[string]interface{}{
			"priority": ticket.Priority,
			"shortDescription": ticket.ShortDescription,
			"state": ticket.State,
			"number": ticket.Number,
		})
	}

	for _, dataObj := range dataObjs {
		batcher.WithObjects(&WeaviateModels.Object {
			Class: TicketClass,
			Properties: dataObj,
		})
	}

	_, err = batcher.Do(context.Background())
	return err
}