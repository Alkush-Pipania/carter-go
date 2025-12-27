package upload

import (
	"time"
)

// PresignUploadRequest represents a request for a presigned upload URL
type PresignUploadRequest struct {
	Filename     string `json:"filename" validate:"required"`
	ContentType  string `json:"content_type" validate:"required"`
	CollectionID string `json:"collection_id" validate:"required,uuid"`
	Title        string `json:"title" validate:"required"`
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
	SourceID string `json:"source_id" validate:"required,uuid"`
}

// SourceProcessingMessage is the message sent to RabbitMQ for processing
type SourceProcessingMessage struct {
	SourceID string `json:"source_id"`
	Type     string `json:"type"`
	UserID   string `json:"user_id"`
}
