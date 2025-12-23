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

// PresignUploadRequest represents a request for a presigned upload URL
type PresignUploadRequest struct {
	Filename     string `json:"filename"`
	ContentType  string `json:"content_type"`
	CollectionID string `json:"collection_id"`
	Title        string `json:"title"`
}

// PresignUploadResponse contains the presigned URL and source info
type PresignUploadResponse struct {
	SourceID  string    `json:"source_id"`
	UploadURL string    `json:"upload_url"`
	S3Key     string    `json:"s3_key"`
	ExpiresAt time.Time `json:"expires_at"`
}

// ConfirmUploadRequest for confirming an upload is complete
type ConfirmUploadRequest struct {
	SourceID string `json:"source_id"`
}

// AllowedContentTypes maps file extensions to MIME types
var AllowedContentTypes = map[string]string{
	"application/pdf":               "pdf",
	"application/vnd.ms-powerpoint": "ppt",
	"application/vnd.openxmlformats-officedocument.presentationml.presentation": "ppt",
	"application/msword": "doc",
	"application/vnd.openxmlformats-officedocument.wordprocessingml.document": "doc",
}

// IsAllowedContentType checks if the content type is allowed for upload
func IsAllowedContentType(contentType string) bool {
	_, ok := AllowedContentTypes[contentType]
	return ok
}

// GetSourceTypeFromContentType returns the source type for a content type
func GetSourceTypeFromContentType(contentType string) string {
	return AllowedContentTypes[contentType]
}
