package handlers

import (
	"encoding/json"
	"net/http"
)

type ChatMessage struct {
	Message string `json:"message"`
	Role    string `json:"role"`
}

type ChatHandler struct {}

func NewChatHandler() *ChatHandler {
	return &ChatHandler{}
}

func (h *ChatHandler) ChatHandler(w http.ResponseWriter, r *http.Request) {
	chat := []ChatMessage{
		{Role: "John", Message: "Hello, how can I help?"},
		{Role: "Jane", Message: "I need assistance with my order."},
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(chat); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
