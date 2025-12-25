package upload

import (
	"context"
	"errors"
	"fmt"

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
	// Validate content type
	if !utils.IsAllowedContentType(req.ContentType) {
		return nil, errors.New("unsupported content type")
	}

	var userUUID pgtype.UUID
	if err := userUUID.Scan(userID); err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	var collectionUUID pgtype.UUID
	if err := collectionUUID.Scan(req.CollectionID); err != nil {
		return nil, fmt.Errorf("invalid collection ID: %w", err)
	}

	// Determine source type from content type
	sourceType := db.SourceType(utils.GetSourceTypeFromContentType(req.ContentType))

	// Generate a new source ID
	newSourceID := uuid.New()

	// Generate S3 key
	s3Key := s3.GenerateS3Key(userID, newSourceID.String(), req.Filename)

	// Create source record with pending status
	var sourceUUID pgtype.UUID
	if err := sourceUUID.Scan(newSourceID.String()); err != nil {
		return nil, err
	}

	source, err := s.repo.CreateSourceWithS3Key(ctx, db.CreateSourceWithS3KeyParams{
		UserID:       userUUID,
		CollectionID: collectionUUID,
		Type:         sourceType,
		Title:        req.Title,
		S3Key:        pgtype.Text{String: s3Key, Valid: true},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create source record: %w", err)
	}

	// Generate presigned upload URL
	result, err := s.presigner.GenerateUploadURL(ctx, s3Key, req.ContentType)
	if err != nil {
		return nil, fmt.Errorf("failed to generate presigned URL: %w", err)
	}

	return &PresignUploadResponse{
		SourceID:  source.ID.String(),
		UploadURL: result.URL,
		S3Key:     result.Key,
		ExpiresAt: result.ExpiresAt,
	}, nil
}

// ConfirmUpload marks a source as ready for processing after upload completes
func (s *service) ConfirmUpload(ctx context.Context, userID string, sourceID string) error {
	var sUUID pgtype.UUID
	if err := sUUID.Scan(sourceID); err != nil {
		return fmt.Errorf("invalid source ID: %w", err)
	}

	// Get the source to retrieve its type
	source, err := s.repo.GetSourceByID(ctx, sUUID)
	if err != nil {
		return fmt.Errorf("source not found: %w", err)
	}

	// Update status to processing
	_, err = s.repo.UpdateSourceStatus(ctx, db.UpdateSourceStatusParams{
		ID:     sUUID,
		Status: db.SourceStatusProcessing,
	})
	if err != nil {
		return fmt.Errorf("failed to update source status: %w", err)
	}

	// Publish to RabbitMQ for processing (parsing, embedding, etc.)
	if s.producer != nil {
		msg := SourceProcessingMessage{
			SourceID: sourceID,
			Type:     string(source.Type),
			UserID:   userID,
		}
		if err := s.producer.Publish(ctx, RoutingKeySourceProcess, msg); err != nil {
			// Log error but don't fail - source is confirmed
			fmt.Printf("failed to publish source for processing: %v\n", err)
		}
	}

	return nil
}
