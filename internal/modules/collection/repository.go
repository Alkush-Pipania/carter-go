package collection

import (
	"context"

	"github.com/Alkush-Pipania/carter-go/pkg/db"
	"github.com/jackc/pgx/v5/pgtype"
)

type repository struct {
	db *db.Queries
}

func NewRepository(db *db.Queries) *repository {
	return &repository{
		db: db,
	}
}

func (r *repository) GetCollectionsByUserID(ctx context.Context, userID pgtype.UUID) ([]db.Collection, error) {
	collections, err := r.db.GetCollectionsByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	return collections, nil
}
