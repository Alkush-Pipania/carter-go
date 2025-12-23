package source

import (
	"context"
	"errors"
	"fmt"

	"github.com/Alkush-Pipania/carter-go/pkg/db"
	"github.com/Alkush-Pipania/carter-go/pkg/rabbitmq"
	"github.com/Alkush-Pipania/carter-go/pkg/s3"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type Repository interface {
	GetSourcesByCollectionID(ctx context.Context, collectionID pgtype.UUID) ([]db.Source, error)
	CreateSource(ctx context.Context, args db.CreateSourceParams) (db.Source, error)
	CreateSourceWithS3Key(ctx context.Context, args db.CreateSourceWithS3KeyParams) (db.Source, error)
	UpdateSourceStatus(ctx context.Context, args db.UpdateSourceStatusParams) (db.Source, error)
}

type service struct {
	repo      Repository
	producer  *rabbitmq.Producer
	presigner *s3.Presigner
}

func NewService(repo Repository, producer *rabbitmq.Producer, presigner *s3.Presigner) Service {
	return &service{
		repo:      repo,
		producer:  producer,
		presigner: presigner,
	}
}

func (s *service) GetSourcesByCollectionID(ctx context.Context, collectionID string) ([]db.Source, error) {
	var id pgtype.UUID
	if err := id.Scan(collectionID); err != nil {
		return nil, err
	}
	return s.repo.GetSourcesByCollectionID(ctx, id)
}

func (s *service) CreateSource(ctx context.Context, userID string, req CreateSourceRequest) (db.Source, error) {
	var userUUID pgtype.UUID
	if err := userUUID.Scan(userID); err != nil {
		return db.Source{}, err
	}

	var collectionUUID pgtype.UUID
	if err := collectionUUID.Scan(req.CollectionID); err != nil {
		return db.Source{}, err
	}

	sourceType := db.SourceType(req.Type)
	if sourceType != "link" {
		return db.Source{}, errors.New("invalid source type")
	}

	var originalUrl pgtype.Text
	if req.OriginalUrl != "" {
		originalUrl = pgtype.Text{String: req.OriginalUrl, Valid: true}
	}

	// Create the source in the database (status defaults to 'pending' in DB)
	source, err := s.repo.CreateSource(ctx, db.CreateSourceParams{
		UserID:       userUUID,
		CollectionID: collectionUUID,
		Type:         sourceType,
		Title:        req.Title,
		OriginalUrl:  originalUrl,
	})
	if err != nil {
		return db.Source{}, err
	}

	// TODO: Add to queue for embedding/processing based on source type

	return source, nil
}

// RequestUploadURL creates a source record and returns a presigned S3 upload URL
func (s *service) RequestUploadURL(ctx context.Context, userID string, req PresignUploadRequest) (*PresignUploadResponse, error) {
	// Validate content type
	if !IsAllowedContentType(req.ContentType) {
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
	sourceType := db.SourceType(GetSourceTypeFromContentType(req.ContentType))

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

	// Update status to processing
	_, err := s.repo.UpdateSourceStatus(ctx, db.UpdateSourceStatusParams{
		ID:     sUUID,
		Status: db.SourceStatusProcessing,
	})
	if err != nil {
		return fmt.Errorf("failed to update source status: %w", err)
	}

	// TODO: Publish to RabbitMQ for embedding processing
	// This is where you would add the queue message for the consumer

	return nil
}
