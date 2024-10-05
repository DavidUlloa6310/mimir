package database

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"os"

	"github.com/weaviate/weaviate-go-client/v4/weaviate/filters"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/graphql"
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

	instanceIDFilter := filters.Where().
        WithPath([]string{"instanceID"}).
        WithOperator(filters.Equal).
        WithValueString(instanceID)
	
	authHashFilter := filters.Where().
		WithPath([]string{"authHash"}).
		WithOperator(filters.Equal).
		WithValueString(hash)
	
    response, err := client.GraphQL().Get().
        WithClassName("Authentication").
        WithFields(
            graphql.Field{Name: "instanceID"},
            graphql.Field{Name: "authHash"},
        ).
        WithWhere(instanceIDFilter).
		WithWhere(authHashFilter).
        WithLimit(1).
        Do(context.Background())
	
	if err != nil {
		return false, err
	}
	
	if response.Data == nil {
		return false, nil
	}

	return true, nil
}

func RegisterAuthentication(instanceID string, username string, password string) error {
	hash := CreateHash(username, password)
	client, err := GetWeaviateClient()

	if err != nil {
		return err
	}

	_, err = client.Data().Creator().
		WithClassName(ChatMessageClass).
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
