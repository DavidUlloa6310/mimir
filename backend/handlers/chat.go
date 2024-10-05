package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/davidulloa/mimir/database"
	"github.com/davidulloa/mimir/models"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

type ChatHandler struct{}

func NewChatHandler() *ChatHandler {
	return &ChatHandler{}
}

func (h *ChatHandler) verifyBasicAuth(r *http.Request, instanceID string) error {
    username, password, ok := r.BasicAuth()
    if !ok {
        return fmt.Errorf("basic authentication required")
    }

    validated, err := database.ValidateAuthentication(instanceID, username, password)
    if err != nil {
        return fmt.Errorf("authentication validation error: %v", err)
    }

    if !validated {
        return fmt.Errorf("invalid credentials")
    }

    return nil
}

// ChatHandler handles all POST requests for chat operations
func (h *ChatHandler) ChatHandler(w http.ResponseWriter, r *http.Request) {
    var body map[string]interface{}
    if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    instanceID, ok := body["instance_id"].(string)
    if !ok {
        http.Error(w, "instance_id is required", http.StatusBadRequest)
        return
    }

    if err := h.verifyBasicAuth(r, instanceID); err != nil {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }

	if createThread, ok := body["createThread"].(bool); ok && createThread {
		h.createChatThread(w, body)
		return
	}

	if threadID, ok := body["threadId"].(string); ok {
		if _, ok := body["message"]; ok {
			h.postNewMessage(w, body)
		} else {
			h.fetchChatThread(w, threadID)
		}
		return
	}

	h.fetchAllChatThreads(w, instanceID)
}

func (h *ChatHandler) createChatThread(w http.ResponseWriter, body map[string]interface{}) {
	acceleratorID, ok := body["acceleratorId"].(string)

	if !ok {
		http.Error(w, "acceleratorId is required", http.StatusBadRequest)
		return
	}

	instanceID, ok := body["instanceId"].(string)

	if !ok {
		http.Error(w, "instanceId is required", http.StatusBadRequest)
		return
	}

	thread := models.ChatThread{
		UserID:        instanceID, 
		Title:         "New Chat Thread",
		IsActive:      true,
		AcceleratorId: acceleratorID,
	}

	threadID, err := database.CreateChatThread(thread)
	if err != nil {
		log.Printf("Error creating chat thread: %v", err)
		http.Error(w, "Error creating chat thread", http.StatusInternalServerError)
		return
	}
	
	go h.generateInitialBotResponse(threadID, acceleratorID)

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{
        "threadId": threadID,
    })
}

func (h *ChatHandler) getBotResponse(systemPrompt string, threadID string, userMessage models.ChatMessage) string {
	client := openai.NewClient(
		option.WithAPIKey(os.Getenv("OPENAI_API_KEY")),
	)
    previousMessages, err := database.GetChatMessages(threadID)
    if err != nil {
        log.Printf("Error fetching previous messages for thread %s: %v", threadID, err)
        return "I'm sorry, I encountered an error while processing your request."
    }

    messages := []openai.ChatCompletionMessageParamUnion{
        openai.SystemMessage(systemPrompt),
    }

    for _, msg := range previousMessages {
        if msg.Role == "user" {
            messages = append(messages, openai.UserMessage(msg.Content))
        } else if msg.Role == "assistant" {
            messages = append(messages, openai.AssistantMessage(msg.Content))
        }
    }

    messages = append(messages, openai.UserMessage(userMessage.Content))

    chat, err := client.Chat.Completions.New(context.TODO(), openai.ChatCompletionNewParams{
        Messages: openai.F(messages),
        Model:    openai.F(openai.ChatModelGPT4o2024_08_06),
    })

    if err != nil {
        log.Printf("Error generating bot response for thread %s: %v", threadID, err)
        return "I apologize, but I'm having trouble generating a response right now. Please try again later."
    }

    if len(chat.Choices) == 0 || chat.Choices[0].Message.Content == "" {
        log.Printf("Received empty response from OpenAI for thread %s", threadID)
        return "I'm sorry, but I couldn't generate a meaningful response. Please rephrase your question or try again later."
    }

    return chat.Choices[0].Message.Content
}

func (h* ChatHandler) generateSystemPrompt(acceleratorID string) (string, error) {
	accelerator, err := database.GetAcceleratorByID(acceleratorID)

	if err != nil {
        log.Printf("Error getting accelerator information for system prompt %v", err)
        return "", err
    }
	
	title := accelerator.Title
	description := accelerator.Description
	category := accelerator.Category

	systemPrompt := fmt.Sprintf(`You are an AI assistant specializing in ServiceNow accelerators. Your role is to provide information and recommendations about a specific ServiceNow accelerator to help companies optimize their ServiceNow implementation.

	You have access to the following information about the accelerator:
	
	Title: %s
	Description: %s
	Category: %s
	
	Your task is to:
	
	1. Understand the accelerator's purpose and benefits based on the provided information.
	2. Explain how this accelerator can help companies improve their ServiceNow implementation.
	3. Provide context on when and why a company might want to use this particular accelerator.
	4. Answer questions about the accelerator's features, implementation process, and potential outcomes.
	5. Relate the accelerator to the broader category it belongs to (%s) and explain its significance within that context.
	
	Remember to:
	- Be informative and professional in your responses.
	- Tailor your explanations to the company's potential needs and challenges.
	- Avoid discussing other accelerators not mentioned in the provided information.
	- If asked about something outside your knowledge scope, politely explain that you can only provide information about the specific accelerator you're trained on.
	
	Engage in a helpful dialogue to assist companies in understanding how this ServiceNow accelerator can benefit their organization.`, title, description, category, category)
	
	return systemPrompt, nil
}

func (h *ChatHandler) generateInitialBotResponse(threadID, acceleratorID string) {
    systemPrompt, err := h.generateSystemPrompt(acceleratorID)
    if err != nil {
        log.Printf("Error adding initial user message to thread %s: %v", threadID, err)
        return
    }

    userMessage := models.ChatMessage{
		Content: "How can I use this accelerator in my service?",
		Role: "user",
	}

	err = database.AddChatMessage(threadID, userMessage)
    if err != nil {
        log.Printf("Error adding initial user message to thread %s: %v", threadID, err)
        return
    }
	
    botResponse := h.getBotResponse(systemPrompt, threadID, userMessage)

    botMessage := models.ChatMessage{
        Content: botResponse,
        Role:    "assistant",
    }

    err = database.AddChatMessage(threadID, botMessage)
    if err != nil {
        log.Printf("Error adding bot response to thread %s: %v", threadID, err)
		return
    }
}

func (h *ChatHandler) fetchChatThread(w http.ResponseWriter, threadID string) {
	chatThread, err := database.GetChatThread(threadID)
	if err != nil {
		log.Printf("Error fetching chat thread: %v", err)
		http.Error(w, "Error fetching chat thread", http.StatusInternalServerError)
		return
	}

	status := "ready"
    if len(chatThread.Messages) == 0 {
        status = "processing"
    }

    response := struct {
        *models.ChatThread
        Status string `json:"status"`
    }{
        ChatThread: chatThread,
        Status:     status,
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}

func (h *ChatHandler) postNewMessage(w http.ResponseWriter, body map[string]interface{}) {
	threadID := body["threadId"].(string)

	messageContent := body["message"].(map[string]interface{})["content"].(string)

	message := models.ChatMessage{
		Content: messageContent,
		Role:    "user",
	}

	err := database.AddChatMessage(threadID, message)
	if err != nil {
		log.Printf("Error adding user message: %v", err)
		http.Error(w, "Error adding message", http.StatusInternalServerError)
		return
	}

	botMessage := models.ChatMessage{
		Content: "This is a bot response", 
		Role:    "bot",
	}

	err = database.AddChatMessage(threadID, botMessage)
	if err != nil {
		log.Printf("Error adding bot message: %v", err)
		http.Error(w, "Error adding bot message", http.StatusInternalServerError)
		return
	}

	chatThread, err := database.GetChatThread(threadID)
	if err != nil {
		log.Printf("Error fetching chat thread: %v", err)
		http.Error(w, "Error fetching chat thread", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(chatThread)
}

func (h *ChatHandler) fetchAllChatThreads(w http.ResponseWriter, instanceID string) {
	chatThreads, err := database.GetChatThreadsByInstanceID(instanceID)
	if err != nil {
		log.Printf("Error fetching chat threads: %v", err)
		http.Error(w, "Error fetching chat threads", http.StatusInternalServerError)
		return
	}

	minimizedThreads := make([]map[string]interface{}, len(chatThreads))
	for i, thread := range chatThreads {
		minimizedThreads[i] = map[string]interface{}{
			"threadId": thread.ID,
			"title":    thread.Title,
			"isActive": thread.IsActive,
			"acceleratorId": thread.AcceleratorId,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(minimizedThreads)
}