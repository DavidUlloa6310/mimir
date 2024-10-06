package database

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"os"

	"github.com/weaviate/weaviate-go-client/v4/weaviate/filters"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/graphql"
)

const (
	AuthorizationClass = "Authorization"
)

func CreateHash(username string, password string) string {
	combined := username + password + os.Getenv("AUTHORIZATION_SALT")
	h := sha256.New()
	h.Write([]byte(combined))
    hashed := h.Sum(nil)
    return hex.EncodeToString(hashed)
}

func ValidateAuthentication(instanceID string, username string, password string) (bool, error) {
	hash := CreateHash(username, password)
	client, err := GetWeaviateClient()

	if err != nil {
		return false, err
	}

	fields := []string{"instanceID"}
	graphqlFields := make([]graphql.Field, len(fields))
	for i, field := range fields {
		graphqlFields[i] = graphql.Field{Name: field}
	}

    response, err := client.GraphQL().Get().
        WithClassName(AuthorizationClass).
		WithFields(graphqlFields...).
		WithWhere(filters.Where().WithOperator(filters.And).WithOperands([]*filters.WhereBuilder{
			filters.Where().WithPath([]string{"instanceID"}).WithOperator(filters.Equal).WithValueString(instanceID),
			filters.Where().WithPath([]string{"authHash"}).WithOperator(filters.Equal).WithValueString(hash),
		})).
        WithLimit(1).
        Do(context.Background())
	
	if err != nil {
		return false, err
	}
	
	getObject, ok := response.Data["Get"].(map[string]interface{})
	if !ok {
		return false, errors.New("unable to parse 'Get' from response data")
	}

	classObjects, ok := getObject[AuthorizationClass].([]interface{})
	if !ok {
		return false, errors.New("unable to parse 'Authorization' class from class object")
	}

	// Returns true if object is found
	return len(classObjects) > 0, nil

}

func RegisterAuthentication(instanceID string, username string, password string) error {
	hash := CreateHash(username, password)
	client, err := GetWeaviateClient()

	if err != nil {
		return err
	}

	_, err = client.Data().Creator().
		WithClassName(AuthorizationClass).
		WithProperties(map[string]interface{}{
			"instanceID":  instanceID,
			"authHash": hash,
		}).
		Do(context.Background())
	
	if err != nil {
		return err
	}

	return nil
}
