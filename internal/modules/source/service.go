package source

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/Alkush-Pipania/carter-go/internal/middleware"
	"github.com/Alkush-Pipania/carter-go/pkg/db"
	"github.com/Alkush-Pipania/carter-go/pkg/rabbitmq"
	"github.com/jackc/pgx/v5/pgtype"
)

type Repository interface {
	GetSourcesByCollectionID(ctx context.Context, collectionID pgtype.UUID) ([]db.Source, error)
	CreateSource(ctx context.Context, args db.CreateSourceParams) (db.Source, error)
	CreateSourceContent(ctx context.Context, args db.CreateSourceContentParams) (db.SourceContent, error)
	GetSourceByID(ctx context.Context, id pgtype.UUID) (db.Source, error)
}

type service struct {
	repo     Repository
	producer *rabbitmq.Publisher
	dltpbh   *rabbitmq.Publisher
}

func NewService(repo Repository, producer *rabbitmq.Publisher, dltpbh *rabbitmq.Publisher) Service {
	return &service{
		repo:     repo,
		producer: producer,
		dltpbh:   dltpbh,
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
	switch req.Type {
	case "link":
		return s.createLinkSource(ctx, userID, req)
	case "note":
		return s.createNoteSource(ctx, userID, req)
	default:
		return db.Source{}, errors.New("invalid source type: use /upload for documents")
	}
}

// createLinkSource handles creation of link-type sources
func (s *service) createLinkSource(ctx context.Context, userID string, req CreateSourceRequest) (db.Source, error) {
	if req.OriginalUrl == "" {
		return db.Source{}, errors.New("original_url is required for link type")
	}

	var userUUID pgtype.UUID
	if err := userUUID.Scan(userID); err != nil {
		return db.Source{}, fmt.Errorf("invalid user ID: %w", err)
	}

	var collectionUUID pgtype.UUID
	if err := collectionUUID.Scan(req.CollectionID); err != nil {
		return db.Source{}, fmt.Errorf("invalid collection ID: %w", err)
	}

	// Create the source in the database (status defaults to 'pending' in DB)
	source, err := s.repo.CreateSource(ctx, db.CreateSourceParams{
		UserID:       userUUID,
		CollectionID: collectionUUID,
		Type:         db.SourceTypeLink,
		Title:        "",
		OriginalUrl:  pgtype.Text{String: req.OriginalUrl, Valid: true},
	})
	if err != nil {
		return db.Source{}, fmt.Errorf("failed to create source: %w", err)
	}

	// Publish to queue for processing (scraping, embedding, etc.)
	if err := s.publishSourceProcessing(ctx, source, userID); err != nil {
		// Log error but don't fail the request - source is created
		// TODO: Add proper logging
		fmt.Printf("failed to publish source for processing: %v\n", err)
	}

	return source, nil
}

// createNoteSource handles creation of note-type sources
func (s *service) createNoteSource(ctx context.Context, userID string, req CreateSourceRequest) (db.Source, error) {
	if req.Content == "" {
		return db.Source{}, errors.New("content is required for note type")
	}

	var userUUID pgtype.UUID
	if err := userUUID.Scan(userID); err != nil {
		return db.Source{}, fmt.Errorf("invalid user ID: %w", err)
	}

	var collectionUUID pgtype.UUID
	if err := collectionUUID.Scan(req.CollectionID); err != nil {
		return db.Source{}, fmt.Errorf("invalid collection ID: %w", err)
	}

	// Create the source in the database
	source, err := s.repo.CreateSource(ctx, db.CreateSourceParams{
		UserID:       userUUID,
		CollectionID: collectionUUID,
		Type:         db.SourceTypeNote,
		Title:        "",
		OriginalUrl:  pgtype.Text{Valid: false}, // No URL for notes
	})
	if err != nil {
		return db.Source{}, fmt.Errorf("failed to create source: %w", err)
	}

	// Save the content to source_contents table
	_, err = s.repo.CreateSourceContent(ctx, db.CreateSourceContentParams{
		SourceID:    source.ID,
		ContentText: req.Content,
	})
	if err != nil {
		return db.Source{}, fmt.Errorf("failed to save source content: %w", err)
	}

	// Publish to queue for processing (embedding)
	if err := s.publishSourceProcessing(ctx, source, userID); err != nil {
		// Log error but don't fail the request
		fmt.Printf("failed to publish source for processing: %v\n", err)
	}

	return source, nil
}

// publishSourceProcessing sends a message to RabbitMQ for processing
func (s *service) publishSourceProcessing(ctx context.Context, source db.Source, userID string) error {
	if s.producer == nil {
		return nil // No producer configured, skip
	}

	msgBody, _ := json.Marshal(&SourceProcessingMessage{
		SourceID: source.ID.String(),
		Type:     string(source.Type),
		UserID:   userID,
	})

	err := s.producer.Publish(ctx, msgBody)
	if err != nil {
		return err
	}
	return nil
}

// GetSourceByID returns a single source by ID
func (s *service) GetSourceByID(ctx context.Context, sourceID string) (db.Source, error) {
	var id pgtype.UUID
	if err := id.Scan(sourceID); err != nil {
		return db.Source{}, fmt.Errorf("invalid source ID: %w", err)
	}
	return s.repo.GetSourceByID(ctx, id)
}

// DeleteSource removes a source by ID
func (s *service) DeleteSource(ctx context.Context, sourceID string, sourceType string) error {
	userID, _ := middleware.GetUserIDFromContext(ctx)

	// Get the source to know its type for cleanup
	var id pgtype.UUID
	if err := id.Scan(sourceID); err != nil {
		return fmt.Errorf("invalid source ID: %w", err)
	}

	msgBody, _ := json.Marshal(&SourceDeleteMessage{
		SourceID: sourceID,
		UserID:   userID,
		Type:     sourceType,
	})
	err := s.dltpbh.Publish(ctx, msgBody)
	if err != nil {
		return err
	}
	return nil
}
