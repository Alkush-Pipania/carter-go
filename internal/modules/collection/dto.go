package collection

import (
	"time"
)

// CreateCollectionRequest represents the request to create a collection
type CreateCollectionRequest struct {
	Name string `json:"name" validate:"required,min=1,max=100"`
}

// UpdateCollectionRequest represents the request to update a collection
type UpdateCollectionRequest struct {
	Name string `json:"name" validate:"required,min=1,max=100"`
}

// CollectionResponse represents the API response for a collection
type CollectionResponse struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}
