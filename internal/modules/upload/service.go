package upload

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/Alkush-Pipania/carter-go/pkg/db"
	"github.com/Alkush-Pipania/carter-go/pkg/rabbitmq"
	"github.com/Alkush-Pipania/carter-go/pkg/s3"
	"github.com/Alkush-Pipania/carter-go/pkg/utils"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

const (
	// RoutingKeySourceProcess is the routing key for source processing queue
	RoutingKeySourceProcess = "source.process"
)

type Repository interface {
	CreateSourceWithS3Key(ctx context.Context, args db.CreateSourceWithS3KeyParams) (db.Source, error)
	UpdateSourceStatus(ctx context.Context, args db.UpdateSourceStatusParams) (db.Source, error)
	GetSourceByID(ctx context.Context, id pgtype.UUID) (db.Source, error)
}

type service struct {
	repo      Repository
	presigner *s3.Presigner
	producer  *rabbitmq.Producer
}

func NewService(repo Repository, presigner *s3.Presigner, producer *rabbitmq.Producer) Service {
	return &service{
		repo:      repo,
		presigner: presigner,
		producer:  producer,
	}
}

// RequestUploadURL creates a source record and returns a presigned S3 upload URL
func (s *service) RequestUploadURL(ctx context.Context, userID string, req PresignUploadRequest) (*PresignUploadResponse, error) {
	log.Printf("[RequestUploadURL] Starting request for userID=%s, filename=%s, contentType=%s, collectionID=%s, title=%s",
		userID, req.Filename, req.ContentType, req.CollectionID, req.Title)

	// Validate content type
	if !utils.IsAllowedContentType(req.ContentType) {
		log.Printf("[RequestUploadURL] ERROR: unsupported content type: %s", req.ContentType)
		return nil, errors.New("unsupported content type")
	}
	log.Printf("[RequestUploadURL] Content type validated: %s", req.ContentType)

	var userUUID pgtype.UUID
	if err := userUUID.Scan(userID); err != nil {
		log.Printf("[RequestUploadURL] ERROR: invalid user ID: %v", err)
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}
	log.Printf("[RequestUploadURL] User UUID parsed successfully")

	var collectionUUID pgtype.UUID
	if err := collectionUUID.Scan(req.CollectionID); err != nil {
		log.Printf("[RequestUploadURL] ERROR: invalid collection ID: %v", err)
		return nil, fmt.Errorf("invalid collection ID: %w", err)
	}
	log.Printf("[RequestUploadURL] Collection UUID parsed successfully")

	// Determine source type from content type
	sourceType := db.SourceType(utils.GetSourceTypeFromContentType(req.ContentType))
	log.Printf("[RequestUploadURL] Source type determined: %s", sourceType)

	// Generate a new source ID
	newSourceID := uuid.New()
	log.Printf("[RequestUploadURL] Generated new source ID: %s", newSourceID.String())

	// Generate S3 key
	s3Key := s3.GenerateS3Key(userID, newSourceID.String(), req.Filename)
	log.Printf("[RequestUploadURL] Generated S3 key: %s", s3Key)

	// Create source record with pending status
	var sourceUUID pgtype.UUID
	if err := sourceUUID.Scan(newSourceID.String()); err != nil {
		log.Printf("[RequestUploadURL] ERROR: failed to scan source UUID: %v", err)
		return nil, err
	}

	log.Printf("[RequestUploadURL] Creating source record in database...")
	source, err := s.repo.CreateSourceWithS3Key(ctx, db.CreateSourceWithS3KeyParams{
		UserID:       userUUID,
		CollectionID: collectionUUID,
		Type:         sourceType,
		Title:        req.Title,
		S3Key:        pgtype.Text{String: s3Key, Valid: true},
	})
	if err != nil {
		log.Printf("[RequestUploadURL] ERROR: failed to create source record: %v", err)
		return nil, fmt.Errorf("failed to create source record: %w", err)
	}
	log.Printf("[RequestUploadURL] Source record created successfully with ID: %s", source.ID.String())

	// Generate presigned upload URL
	log.Printf("[RequestUploadURL] Generating presigned URL for S3 key: %s", s3Key)
	result, err := s.presigner.GenerateUploadURL(ctx, s3Key, req.ContentType)
	if err != nil {
		log.Printf("[RequestUploadURL] ERROR: failed to generate presigned URL: %v", err)
		return nil, fmt.Errorf("failed to generate presigned URL: %w", err)
	}
	log.Printf("[RequestUploadURL] Presigned URL generated successfully, expires at: %v", result.ExpiresAt)

	log.Printf("[RequestUploadURL] Request completed successfully for source ID: %s", source.ID.String())
	return &PresignUploadResponse{
		SourceID:  source.ID.String(),
		UploadURL: result.URL,
		S3Key:     result.Key,
		ExpiresAt: result.ExpiresAt,
	}, nil
}

// ConfirmUpload marks a source as ready for processing after upload completes
func (s *service) ConfirmUpload(ctx context.Context, userID string, sourceID string) error {
	log.Printf("[ConfirmUpload] Starting confirmation for userID=%s, sourceID=%s", userID, sourceID)

	var sUUID pgtype.UUID
	if err := sUUID.Scan(sourceID); err != nil {
		log.Printf("[ConfirmUpload] ERROR: invalid source ID: %v", err)
		return fmt.Errorf("invalid source ID: %w", err)
	}
	log.Printf("[ConfirmUpload] Source UUID parsed successfully")

	// Get the source to retrieve its type
	log.Printf("[ConfirmUpload] Fetching source from database...")
	source, err := s.repo.GetSourceByID(ctx, sUUID)
	if err != nil {
		log.Printf("[ConfirmUpload] ERROR: source not found: %v", err)
		return fmt.Errorf("source not found: %w", err)
	}
	log.Printf("[ConfirmUpload] Source found with type: %s", source.Type)

	// Update status to processing
	log.Printf("[ConfirmUpload] Updating source status to processing...")
	_, err = s.repo.UpdateSourceStatus(ctx, db.UpdateSourceStatusParams{
		ID:     sUUID,
		Status: db.SourceStatusProcessing,
	})
	if err != nil {
		log.Printf("[ConfirmUpload] ERROR: failed to update source status: %v", err)
		return fmt.Errorf("failed to update source status: %w", err)
	}
	log.Printf("[ConfirmUpload] Source status updated to processing")

	// Publish to RabbitMQ for processing (parsing, embedding, etc.)
	if s.producer != nil {
		log.Printf("[ConfirmUpload] Publishing to RabbitMQ with routing key: %s", RoutingKeySourceProcess)
		msg := SourceProcessingMessage{
			SourceID: sourceID,
			Type:     string(source.Type),
			UserID:   userID,
		}
		if err := s.producer.Publish(ctx, RoutingKeySourceProcess, msg); err != nil {
			// Log error but don't fail - source is confirmed
			log.Printf("[ConfirmUpload] WARNING: failed to publish source for processing: %v", err)
		} else {
			log.Printf("[ConfirmUpload] Successfully published to RabbitMQ")
		}
	} else {
		log.Printf("[ConfirmUpload] WARNING: RabbitMQ producer is nil, skipping publish")
	}

	log.Printf("[ConfirmUpload] Confirmation completed successfully for sourceID=%s", sourceID)
	return nil
}
