package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/davidulloa/mimir/models"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/fault"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/filters"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/graphql"
)

const (
	ChatThreadClass  = "ChatThread"
	ChatMessageClass = "ChatMessage"
)

func CreateChatThread(thread models.ChatThread) (string, error) {
	client, err := GetWeaviateClient()
	if err != nil {
		log.Printf("Error getting Weaviate client: %v", err)
		return "", err
	}

	thread.CreatedAt = time.Now()
	thread.UpdatedAt = time.Now()

	response, err := client.Data().Creator().
		WithClassName(ChatThreadClass).
		WithProperties(map[string]interface{}{
			"userID":        thread.UserID,
			"title":         thread.Title,
			"createdAt":     thread.CreatedAt,
			"updatedAt":     thread.UpdatedAt,
			"isActive":      thread.IsActive,
			"metadata":      thread.Metadata,
			"acceleratorID": thread.AcceleratorId,
		}).
		Do(context.Background())

	if err != nil {
		log.Printf("Error creating chat thread: %v", err)
		return "", err
	}

	// Retrieve the generated ID from the Additional field
	threadID := response.Object.ID

	log.Printf("Chat thread created successfully with ID: %s", threadID)
	return string(threadID), nil
}

func GetChatThread(threadID string) (*models.ChatThread, error) {
	client, err := GetWeaviateClient()
	if err != nil {
		log.Printf("Error getting Weaviate client: %v", err)
		return nil, err
	}

	result, err := client.Data().ObjectsGetter().
		WithClassName(ChatThreadClass).
		WithID(threadID).
		Do(context.Background())

	if err != nil {
		if clientErr, ok := err.(*fault.WeaviateClientError); ok {
			if clientErr.StatusCode == 404 {
				log.Printf("Chat thread not found with ID: %s", threadID)
				return nil, fmt.Errorf("chat thread not found")
			}
		}
		log.Printf("Error retrieving chat thread with ID %s: %v", threadID, err)
		return nil, err
	}

	if len(result) == 0 {
		log.Printf("Chat thread not found with ID: %s", threadID)
		return nil, fmt.Errorf("chat thread not found")
	}

	thread := &models.ChatThread{}
	properties, ok := result[0].Properties.(map[string]interface{})
	if !ok {
		log.Printf("Expected properties to be a map but got: %v", result[0].Properties)
		return nil, fmt.Errorf("invalid properties for chat thread with ID: %s", threadID)
	}

	thread.ID = threadID
	thread.UserID = properties["userID"].(string)
	thread.Title = properties["title"].(string)

	if createdAt, ok := properties["createdAt"].(string); ok {
		thread.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
	}
	if updatedAt, ok := properties["updatedAt"].(string); ok {
		thread.UpdatedAt, _ = time.Parse(time.RFC3339, updatedAt)
	}

	thread.IsActive = properties["isActive"].(bool)
	thread.Metadata = properties["metadata"].(string)
	thread.AcceleratorId = properties["acceleratorID"].(string)

	thread.Messages, err = GetChatMessages(threadID)
	if err != nil {
		log.Printf("Error retrieving messages for chat thread ID %s: %v", threadID, err)
		return nil, err
	}

	log.Printf("Retrieved chat thread with ID: %s", threadID)
	return thread, nil
}

func UpdateChatThread(thread models.ChatThread) error {
	client, err := GetWeaviateClient()
	if err != nil {
		log.Printf("Error getting Weaviate client: %v", err)
		return err
	}

	thread.UpdatedAt = time.Now()

	log.Printf("Updating chat thread with ID: %s", thread.ID)

	err = client.Data().Updater().
		WithClassName(ChatThreadClass).
		WithID(thread.ID).
		WithProperties(map[string]interface{}{
			"userID":        thread.UserID,
			"title":         thread.Title,
			"updatedAt":     thread.UpdatedAt,
			"createdAt":     thread.CreatedAt,
			"isActive":      thread.IsActive,
			"metadata":      thread.Metadata,
			"acceleratorID": thread.AcceleratorId,
		}).
		Do(context.Background())

	if err != nil {
		if clientErr, ok := err.(*fault.WeaviateClientError); ok {
			if clientErr.StatusCode == 404 {
				log.Printf("Chat thread not found for update with ID: %s", thread.ID)
				return fmt.Errorf("chat thread not found")
			}
		}
		log.Printf("Error updating chat thread with ID %s: %v", thread.ID, err)
		return err
	}

	log.Printf("Chat thread updated successfully with ID: %s", thread.ID)
	return nil
}

func EditChatThreadTitle(threadID string) error {
	thread, err := GetChatThread(threadID)
	if err != nil {
		return err
	}

	newTitle := GenerateTitle(thread.Messages)
	thread.Title = newTitle

	return UpdateChatThread(*thread)
}

func GenerateTitle(messages []models.ChatMessage) string {
	apiKey := os.Getenv("OPENAI_API_KEY")

	client := openai.NewClient(
		option.WithAPIKey(apiKey),
	)

	prompt := "Based on the following chat messages, generate a short title (5 words maximum) that summarizes the conversation:"

	var conversationContent string
	for _, message := range messages {
		conversationContent += message.Content + "\n"
	}

	content := fmt.Sprintf("%s\n%s", prompt, conversationContent)

	chat, err := client.Chat.Completions.New(context.TODO(), openai.ChatCompletionNewParams{
		Messages: openai.F([]openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage("You are a helpful assistant."),
			openai.UserMessage(content),
		}),
		Model: openai.F(openai.ChatModelGPT4o2024_08_06),
	})

	if err != nil {
		return "Unnamed Chat"
	}

	if len(chat.Choices) == 0 {
		return "Unnamed Chat"
	}

	return chat.Choices[0].Message.Content
}

func DeleteChatThread(threadID string) error {
	client, err := GetWeaviateClient()
	if err != nil {
		log.Printf("Error getting Weaviate client: %v", err)
		return err
	}

	log.Printf("Deleting chat thread with ID: %s", threadID)
	err = client.Data().Deleter().
		WithClassName(ChatThreadClass).
		WithID(threadID).
		Do(context.Background())

	if err != nil {
		if clientErr, ok := err.(*fault.WeaviateClientError); ok {
			if clientErr.StatusCode == 404 {
				log.Printf("Chat thread not found for update with ID: %s", threadID)
				return fmt.Errorf("chat thread not found")
			}
		}
		log.Printf("Error updating chat thread with ID %s: %v", threadID, err)
		return err
	}

	log.Printf("Chat thread deleted successfully with ID: %s", threadID)
	return nil
}

func AddChatMessage(threadID string, message models.ChatMessage) error {
	client, err := GetWeaviateClient()
	if err != nil {
		log.Printf("Error getting Weaviate client: %v", err)
		return err
	}

	if message.Timestamp.IsZero() {
		message.Timestamp = time.Now()
		log.Printf("Setting current timestamp for message: %s", message.Content)
	}

	response, err := client.Data().Creator().
		WithClassName(ChatMessageClass).
		WithProperties(map[string]interface{}{
			"threadID":  threadID,
			"role":      message.Role,
			"content":   message.Content,
			"timestamp": message.Timestamp,
		}).
		Do(context.Background())

	if err != nil {
		log.Printf("Error adding chat message to thread ID %s: %v", threadID, err)
		return err
	}

	messageID := response.Object.ID

	log.Printf("Chat message added successfully to thread ID: %s with message ID: %s", threadID, messageID)

	return nil
}

func GetChatMessages(threadID string) ([]models.ChatMessage, error) {
	client, err := GetWeaviateClient()
	if err != nil {
		log.Printf("Error getting Weaviate client: %v", err)
		return nil, err
	}

	fields := []string{"role", "content", "timestamp", "_additional{id}"}
	graphqlFields := make([]graphql.Field, len(fields))
	for i, field := range fields {
		graphqlFields[i] = graphql.Field{Name: field}
	}

	whereFilter := filters.Where().
		WithPath([]string{"threadID"}).
		WithOperator(filters.Equal).
		WithValueString(threadID)

	log.Printf("Fetching chat messages for thread ID: %s", threadID)
	result, err := client.GraphQL().Get().
		WithClassName(ChatMessageClass).
		WithFields(graphqlFields...).
		WithWhere(whereFilter).
		Do(context.Background())

	if err != nil {
		log.Printf("Error retrieving chat messages for thread ID %s: %v", threadID, err)
		return nil, err
	}

	if result.Errors != nil {
		for _, err := range result.Errors {
			log.Printf("GraphQL error: %v", err)
		}
		return nil, fmt.Errorf("graphQL errors: %v", result.Errors)
	}

	var messages []models.ChatMessage
	chatMessages, ok := result.Data["Get"].(map[string]interface{})[ChatMessageClass]
	if !ok {
		log.Printf("Unexpected data structure: %v", result.Data)
		return nil, fmt.Errorf("unexpected data structure: %v", result.Data)
	}

	chatMessagesSlice, ok := chatMessages.([]interface{})
	if !ok {
		log.Printf("Expected slice of messages but got: %v", chatMessages)
		return nil, fmt.Errorf("expected slice of messages but got: %v", chatMessages)
	}

	for _, obj := range chatMessagesSlice {
		msg, ok := obj.(map[string]interface{})
		if !ok {
			log.Printf("Expected message object but got: %v", obj)
			return nil, fmt.Errorf("expected message object but got: %v", obj)
		}

		timestampStr, exists := msg["timestamp"].(string)
		if !exists || timestampStr == "" {
			log.Printf("Timestamp field is missing or empty for message ID: %v", msg["id"])
			return nil, fmt.Errorf("timestamp field is missing or empty for message ID: %v", msg["id"])
		}

		timestamp, err := time.Parse(time.RFC3339, timestampStr)
		if err != nil {
			log.Printf("Error parsing timestamp for message ID %v: %v", msg["id"], err)
			return nil, fmt.Errorf("error parsing timestamp for message ID %v: %v", msg["id"], err)
		}

		messages = append(messages, models.ChatMessage{
			ID:        msg["_additional"].(map[string]interface{})["id"].(string),
			Role:      msg["role"].(string),
			Content:   msg["content"].(string),
			Timestamp: timestamp,
		})
	}

	log.Printf("Retrieved %d chat messages for thread ID: %s", len(messages), threadID)
	return messages, nil
}

func GetChatThreadsByInstanceID(instanceID string) ([]models.ChatThread, error) {
	client, err := GetWeaviateClient()
	if err != nil {
		log.Printf("Error getting Weaviate client: %v", err)
		return nil, err
	}

	fields := []string{"userID", "title", "createdAt", "updatedAt", "isActive", "metadata", "acceleratorID", "_additional{id}"}
	graphqlFields := make([]graphql.Field, len(fields))
	for i, field := range fields {
		graphqlFields[i] = graphql.Field{Name: field}
	}

	whereFilter := filters.Where().
		WithPath([]string{"userID"}). // userID = instanceID
		WithOperator(filters.Equal).
		WithValueString(instanceID)

	log.Printf("Fetching chat threads for accelerator ID: %s", instanceID)
	result, err := client.GraphQL().Get().
		WithClassName(ChatThreadClass).
		WithFields(graphqlFields...).
		WithWhere(whereFilter).
		Do(context.Background())

	if err != nil {
		log.Printf("Error retrieving chat threads for accelerator ID %s: %v", instanceID, err)
		return nil, err
	}

	if result.Errors != nil {
		for _, err := range result.Errors {
			log.Printf("GraphQL error: %v", err)
		}
		return nil, fmt.Errorf("graphQL errors: %v", result.Errors)
	}
	var threads []models.ChatThread
	chatThreads, ok := result.Data["Get"].(map[string]interface{})[ChatThreadClass]
	if !ok {
		log.Printf("Unexpected data structure: %v", result.Data)
		return nil, fmt.Errorf("unexpected data structure: %v", result.Data)
	}

	chatThreadsSlice, ok := chatThreads.([]interface{})
	if !ok {
		log.Printf("Expected slice of threads but got: %v", chatThreads)
		return nil, fmt.Errorf("expected slice of threads but got: %v", chatThreads)
	}

	for _, obj := range chatThreadsSlice {
		thread, ok := obj.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("expected thread object but got: %v", obj)
		}

		if thread["createdAt"] == nil || thread["updatedAt"] == nil {
			log.Printf("Skipping thread due to nil createdAt or updatedAt: %v", thread)
			continue
		}

		createdAt, err := time.Parse(time.RFC3339, thread["createdAt"].(string))
		if err != nil {
			return nil, fmt.Errorf("error parsing createdAt: %v", err)
		}

		updatedAt, err := time.Parse(time.RFC3339, thread["updatedAt"].(string))
		if err != nil {
			return nil, fmt.Errorf("error parsing updatedAt: %v", err)
		}

		if thread["_additional"] == nil {
			log.Printf("Skipping thread due to nil _additional field: %v", thread)
			continue
		}

		additional := thread["_additional"].(map[string]interface{})
		if additional["id"] == nil || thread["userID"] == nil || thread["title"] == nil {
			log.Printf("Skipping thread due to nil id, userID, or title: %v", thread)
			continue // Skip if necessary fields are nil
		}

		threads = append(threads, models.ChatThread{
			ID:            additional["id"].(string),
			UserID:        thread["userID"].(string),
			Title:         thread["title"].(string),
			CreatedAt:     createdAt,
			UpdatedAt:     updatedAt,
			IsActive:      thread["isActive"].(bool),
			Metadata:      thread["metadata"].(string),
			AcceleratorId: thread["acceleratorID"].(string),
		})
	}

	log.Printf("Retrieved %d chat threads for accelerator ID: %s", len(threads), instanceID)
	return threads, nil
}
