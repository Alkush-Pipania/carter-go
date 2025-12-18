package source

import (
	"time"

	"github.com/Alkush-Pipania/carter-go/pkg/db"
)

// CreateSourceRequest represents the incoming request body for creating a source
type CreateSourceRequest struct {
	CollectionID string `json:"collection_id"`
	Type         string `json:"type"` // "link", "pdf", "ppt", "doc", "note"
	Title        string `json:"title"`
	OriginalUrl  string `json:"original_url,omitempty"` // Required for link type
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

// ToSourceResponse converts a db.Source to SourceResponse
func ToSourceResponse(s db.Source) SourceResponse {
	return SourceResponse{
		ID:           s.ID.String(),
		UserID:       s.UserID.String(),
		CollectionID: s.CollectionID.String(),
		Type:         string(s.Type),
		Status:       string(s.Status),
		Title:        s.Title,
		OriginalUrl:  s.OriginalUrl.String,
		CreatedAt:    s.CreatedAt.Time,
	}
}

// ToSourceResponses converts a slice of db.Source to []SourceResponse
func ToSourceResponses(sources []db.Source) []SourceResponse {
	responses := make([]SourceResponse, len(sources))
	for i, s := range sources {
		responses[i] = ToSourceResponse(s)
	}
	return responses
}
