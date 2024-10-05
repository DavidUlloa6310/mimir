package database

import (
	"encoding/json"
	"testing"

	"github.com/davidulloa/mimir/models"
	"github.com/stretchr/testify/assert"
)

func TestWeaviateOperations(t *testing.T) {
	// Step 1: Initialize Weaviate Client
	client, err := InitWeaviateClient()
	if err != nil {
		t.Fatalf("Failed to initialize Weaviate client: %v", err)
	}
	assert.NotNil(t, client)

	metadata := map[string]interface{}{
		"test": "metadata",
	}

	jsonData, _ := json.Marshal(metadata)
	

	// Step 2: Test CreateChatThread
	thread := models.ChatThread{
		UserID:        "user123",
		Title:         "Test Chat Thread",
		IsActive:      true,
		Metadata:      string(jsonData),
		AcceleratorId: "acc123",
	}
	threadID, err := CreateChatThread(thread)
	assert.NoError(t, err)
	assert.NotEmpty(t, threadID)

	// Step 3: Test GetChatThread
	retrievedThread, err := GetChatThread(threadID)
	assert.NoError(t, err)
	assert.Equal(t, thread.UserID, retrievedThread.UserID)
	assert.Equal(t, thread.Title, retrievedThread.Title)
	assert.True(t, retrievedThread.IsActive)
	assert.Equal(t, thread.Metadata, retrievedThread.Metadata)
	assert.Equal(t, thread.AcceleratorId, retrievedThread.AcceleratorId)

	// Step 4: Test UpdateChatThread
	updatedThread := *retrievedThread
	updatedThread.Title = "Updated Test Chat Thread"
	updatedThread.IsActive = false
	err = UpdateChatThread(updatedThread)
	assert.NoError(t, err)

	// Verify the update
	retrievedUpdatedThread, err := GetChatThread(threadID)
	assert.NoError(t, err)
	assert.Equal(t, "Updated Test Chat Thread", retrievedUpdatedThread.Title)
	assert.False(t, retrievedUpdatedThread.IsActive)

	// Step 5: Test AddChatMessage
	userMessage := models.ChatMessage{
		Role:    "user",
		Content: "Hello, GPT!",
	}
	err = AddChatMessage(threadID, userMessage)
	assert.NoError(t, err)

	// Simulate GPT response (you would typically call the OpenAI API here)
	gptMessage := models.ChatMessage{
		Role:    "assistant",
		Content: "Hello! How can I assist you today?",
	}
	err = AddChatMessage(threadID, gptMessage)
	assert.NoError(t, err)

	// Step 6: Test GetChatMessages
	messages, err := GetChatMessages(threadID)
	assert.NoError(t, err)
	assert.Len(t, messages, 2)
	assert.Equal(t, "user", messages[0].Role)
	assert.Equal(t, "Hello, GPT!", messages[0].Content)
	assert.Equal(t, "assistant", messages[1].Role)
	assert.Equal(t, "Hello! How can I assist you today?", messages[1].Content)

	// Step 7: Test DeleteChatThread
	err = DeleteChatThread(threadID)
	assert.NoError(t, err)

	// Verify the deletion
	_, err = GetChatThread(threadID)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "chat thread not found")
}

func TestChatThreadCRUD(t *testing.T) {
	// Create a chat thread

	metadata := map[string]interface{}{
		"test": "metadata",
	}

	jsonData, _ := json.Marshal(metadata)
	

	thread := models.ChatThread{
		UserID:        "user456",
		Title:         "CRUD Test Thread",
		IsActive:      true,
		Metadata:      string(jsonData),
		AcceleratorId: "acc456",
	}
	threadID, err := CreateChatThread(thread)
	if err != nil {
        t.Fatalf("Failed to create chat thread: %v", err)
    }
    assert.NotEmpty(t, threadID)

	// Read the chat thread
	retrievedThread, err := GetChatThread(threadID)
	assert.NoError(t, err)
	assert.Equal(t, thread.UserID, retrievedThread.UserID)
	assert.Equal(t, thread.Title, retrievedThread.Title)

	// Update the chat thread
	retrievedThread.Title = "Updated CRUD Test Thread"
	err = UpdateChatThread(*retrievedThread)
	assert.NoError(t, err)

	// Verify the update
	updatedThread, err := GetChatThread(threadID)
	assert.NoError(t, err)
	assert.Equal(t, "Updated CRUD Test Thread", updatedThread.Title)

	// Delete the chat thread
	err = DeleteChatThread(threadID)
	assert.NoError(t, err)

	// Verify the deletion
	_, err = GetChatThread(threadID)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "chat thread not found")
}

func TestChatMessageOperations(t *testing.T) {
	// Create a chat thread for message testing
	thread := models.ChatThread{
		UserID: "user789",
		Title:  "Message Test Thread",
	}
	threadID, err := CreateChatThread(thread)
	assert.NoError(t, err)

	// Add multiple messages
	messages := []models.ChatMessage{
		{Role: "user", Content: "What's the weather like?"},
		{Role: "assistant", Content: "I'm sorry, but I don't have access to real-time weather information. You would need to check a weather service or app for current conditions."},
		{Role: "user", Content: "Thanks for letting me know."},
	}

	for _, msg := range messages {
		err = AddChatMessage(threadID, msg)
		assert.NoError(t, err)
	}

	// Retrieve and verify messages
	retrievedMessages, err := GetChatMessages(threadID)
	assert.NoError(t, err)
	assert.Len(t, retrievedMessages, len(messages))

	for i, msg := range retrievedMessages {
		assert.Equal(t, messages[i].Role, msg.Role)
		assert.Equal(t, messages[i].Content, msg.Content)
		assert.False(t, msg.Timestamp.IsZero())
	}

	// Clean up
	err = DeleteChatThread(threadID)
	assert.NoError(t, err)
}