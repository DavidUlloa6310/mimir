package models

import (
	"time"
)

type ChatThread struct {
	ID            string        `json:"id"`
	UserID        string        `json:"user_id"`
	Title         string        `json:"title"`
	Messages      []ChatMessage `json:"messages"`
	CreatedAt     time.Time     `json:"created_at"`
	UpdatedAt     time.Time     `json:"updated_at"`
	IsActive      bool          `json:"is_active"`
	Metadata      string        `json:"metadata"`
	AcceleratorId string        `json:"accelerator_id"`
}

type ChatMessage struct {
	ID        string    `json:"id"`
	Role      string    `json:"role"`
	Content   string    `json:"content"`
	Timestamp time.Time `json:"timestamp"`
}
