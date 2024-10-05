package database

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/graphql"
)

const (
	ChatThreadClass  = "ChatThread"
	ChatMessageClass = "ChatMessage"
)

func CreateChatThread(thread ChatThread) (string, error) {
	client, err := GetWeaviateClient()
	if err != nil {
		return "", err
	}

	if thread.ID == "" {
		thread.ID = uuid.New().String()
	}
	thread.CreatedAt = time.Now()
	thread.UpdatedAt = time.Now()

	_, err = client.Data().Creator().
		WithClassName(ChatThreadClass).
		WithID(thread.ID).
		WithProperties(map[string]interface{}{
			"userID":        thread.UserID,
			"botID":         thread.BotID,
			"title":         thread.Title,
			"createdAt":     thread.CreatedAt,
			"updatedAt":     thread.UpdatedAt,
			"isActive":      thread.IsActive,
			"metadata":      thread.Metadata,
			"acceleratorId": thread.AcceleratorId,
		}).
		Do(context.Background())

	if err != nil {
		return "", err
	}

	return thread.ID, nil
}

func GetChatThread(threadID string) (*ChatThread, error) {
	client, err := GetWeaviateClient()
	if err != nil {
		return nil, err
	}

	result, err := client.Data().ObjectsGetter().
		WithClassName(ChatThreadClass).
		WithID(threadID).
		Do(context.Background())

	if err != nil {
		return nil, err
	}

	if len(result) == 0 {
		return nil, fmt.Errorf("chat thread not found")
	}

	thread := &ChatThread{}
	properties := result[0].Properties.(map[string]interface{})

	thread.ID = threadID
	thread.UserID = properties["userID"].(string)
	thread.BotID = properties["botID"].(string)
	thread.Title = properties["title"].(string)
	thread.CreatedAt, _ = time.Parse(time.RFC3339, properties["createdAt"].(string))
	thread.UpdatedAt, _ = time.Parse(time.RFC3339, properties["updatedAt"].(string))
	thread.IsActive = properties["isActive"].(bool)
	thread.Metadata = properties["metadata"].(map[string]interface{})
	thread.AcceleratorId = properties["acceleratorId"].(string)

	thread.Messages, err = GetChatMessages(threadID)
	if err != nil {
		return nil, err
	}

	return thread, nil
}

func UpdateChatThread(thread ChatThread) error {
	client, err := GetWeaviateClient()
	if err != nil {
		return err
	}

	thread.UpdatedAt = time.Now()

	_, err = client.Data().Updater().
		WithClassName(ChatThreadClass).
		WithID(thread.ID).
		WithProperties(map[string]interface{}{
			"userID":        thread.UserID,
			"botID":         thread.BotID,
			"title":         thread.Title,
			"updatedAt":     thread.UpdatedAt,
			"isActive":      thread.IsActive,
			"metadata":      thread.Metadata,
			"acceleratorId": thread.AcceleratorId,
		}).
		Do(context.Background())

	return err
}

func DeleteChatThread(threadID string) error {
	client, err := GetWeaviateClient()
	if err != nil {
		return err
	}

	err = client.Data().Deleter().
		WithClassName(ChatThreadClass).
		WithID(threadID).
		Do(context.Background())

	return err
}

// AddChatMessage adds a new message to a chat thread
func AddChatMessage(threadID string, message ChatMessage) error {
	client, err := GetWeaviateClient()
	if err != nil {
		return err
	}

	if message.ID == "" {
		message.ID = uuid.New().String()
	}
	if message.Timestamp.IsZero() {
		message.Timestamp = time.Now()
	}

	_, err = client.Data().Creator().
		WithClassName(ChatMessageClass).
		WithID(message.ID).
		WithProperties(map[string]interface{}{
			"threadID":  threadID,
			"role":      message.Role,
			"content":   message.Content,
			"timestamp": message.Timestamp,
		}).
		Do(context.Background())

	return err
}

func GetChatMessages(threadID string) ([]ChatMessage, error) {
	client, err := GetWeaviateClient()
	if err != nil {
		return nil, err
	}

	result, err := client.GraphQL().Get().
		WithClassName(ChatMessageClass).
		WithFields("id role content timestamp").
		WithWhere(graphql.Where().
			WithPath([]string{"threadID"}).
			WithOperator(graphql.Equal).
			WithValueString(threadID)).
		Do(context.Background())

	if err != nil {
		return nil, err
	}

	var messages []ChatMessage
	for _, obj := range result.Data["Get"].(map[string]interface{})[ChatMessageClass].([]interface{}) {
		msg := obj.(map[string]interface{})
		timestamp, _ := time.Parse(time.RFC3339, msg["timestamp"].(string))
		messages = append(messages, ChatMessage{
			ID:        msg["id"].(string),
			Role:      msg["role"].(string),
			Content:   msg["content"].(string),
			Timestamp: timestamp,
		})
	}

	return messages, nil
}