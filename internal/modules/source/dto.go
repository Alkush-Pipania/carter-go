package source

import (
	"time"
)

// CreateSourceRequest represents the incoming request body for creating a source
type CreateSourceRequest struct {
	CollectionID string `json:"collection_id" validate:"required,uuid"`
	Type         string `json:"type" validate:"required,oneof=link note"`
	Title        string `json:"title"`
	OriginalUrl  string `json:"original_url,omitempty" validate:"omitempty,url"`
	Content      string `json:"content,omitempty"` // For notes - the text content
}

// SourceResponse represents the API response for a source
type SourceResponse struct {
	ID           string    `json:"id"`
	UserID       string    `json:"user_id"`
	CollectionID string    `json:"collection_id"`
	Type         string    `json:"type"`
	Status       string    `json:"status"`
	Title        string    `json:"title"`
	OriginalUrl  string    `json:"original_url,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
}

// SourceProcessingMessage is the message sent to RabbitMQ for processing
type SourceProcessingMessage struct {
	SourceID string `json:"source_id"`
	Type     string `json:"type"`
	UserID   string `json:"user_id"`
}
