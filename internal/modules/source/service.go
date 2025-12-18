package source

import (
	"context"

	"github.com/Alkush-Pipania/carter-go/pkg/db"
	"github.com/jackc/pgx/v5/pgtype"
)

type Repository interface {
	GetSourcesByCollectionID(ctx context.Context, collectionID pgtype.UUID) ([]db.Source, error)
	CreateSource(ctx context.Context, args db.CreateSourceParams) (db.Source, error)
	UpdateSourceStatus(ctx context.Context, args db.UpdateSourceStatusParams) (db.Source, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{
		repo: repo,
	}
}

func (s *service) GetSourcesByCollectionID(ctx context.Context, collectionID string) ([]db.Source, error) {
	var uuid pgtype.UUID
	if err := uuid.Scan(collectionID); err != nil {
		return nil, err
	}
	return s.repo.GetSourcesByCollectionID(ctx, uuid) // fix it up have to check what it returns
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
