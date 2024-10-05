package models

import (
	"time"
)

type Accelerator struct {
    ID    int    `json:"id"`
    Url   string `json:"url"`
    Title string `json:"title"`
}

type ChatThread struct {
    ID            string                 `json:"id"`
    UserID        string                 `json:"user_id"`
    BotID         string                 `json:"bot_id"`
    Title         string                 `json:"title"`
    Messages      []ChatMessage          `json:"messages"`
    CreatedAt     time.Time              `json:"created_at"`
    UpdatedAt     time.Time              `json:"updated_at"`
    IsActive      bool                   `json:"is_active"`
    Metadata      map[string]interface{} `json:"metadata"`
    AcceleratorId string                 `json:"accelerator_id"`
}

type ChatMessage struct {
    ID        string    `json:"id"`
    Role      string    `json:"role"`
    Content   string    `json:"content"`
    Timestamp time.Time `json:"timestamp"`
}