package collection

import (
	"context"

	"github.com/Alkush-Pipania/carter-go/pkg/db"
	"github.com/jackc/pgx/v5/pgtype"
)

type Repository interface {
	GetCollectionsByUserID(ctx context.Context, userID pgtype.UUID) ([]db.Collection, error)
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
