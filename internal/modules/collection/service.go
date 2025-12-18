package collection

import (
	"context"

	"github.com/Alkush-Pipania/carter-go/pkg/db"
	"github.com/jackc/pgx/v5/pgtype"
)

type Repository interface {
	GetCollectionsByUserID(ctx context.Context, userID pgtype.UUID) ([]db.Collection, error)
	CreateCollection(ctx context.Context, args db.CreateCollectionParams) (db.Collection, error)
	UpdateCollection(ctx context.Context, args db.UpdateCollectionParams) (db.Collection, error)
	DeleteCollection(ctx context.Context, id pgtype.UUID) error
}

type service struct {
	repo Repository
}

func NewService(repository Repository) *service {
	return &service{
		repo: repository,
	}
}

func (s *service) GetCollectionsByUserID(ctx context.Context, userID string) ([]db.Collection, error) {
	var uuid pgtype.UUID
	if err := uuid.Scan(userID); err != nil {
		return nil, err
	}
	return s.repo.GetCollectionsByUserID(ctx, uuid)
}

func (s *service) CreateCollection(ctx context.Context, userID string, name string) (db.Collection, error) {
	var uuid pgtype.UUID
	if err := uuid.Scan(userID); err != nil {
		return db.Collection{}, err
	}

	return s.repo.CreateCollection(ctx, db.CreateCollectionParams{
		UserID: uuid,
		Name:   name,
	})
}

func (s *service) UpdateCollection(ctx context.Context, id string, name string) (db.Collection, error) {
	var uuid pgtype.UUID
	if err := uuid.Scan(id); err != nil {
		return db.Collection{}, err
	}
	return s.repo.UpdateCollection(ctx, db.UpdateCollectionParams{
		ID:   uuid,
		Name: name,
	})
}

func (s *service) DeleteCollection(ctx context.Context, id string) error {
	var uuid pgtype.UUID
	if err := uuid.Scan(id); err != nil {
		return err
	}
	return s.repo.DeleteCollection(ctx, uuid)
}