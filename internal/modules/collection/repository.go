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

func (r *repository) UpdateCollection(ctx context.Context, args db.UpdateCollectionParams) (db.Collection, error) {
	collection, err := r.db.UpdateCollection(ctx, args)
	if err != nil {
		return db.Collection{}, err
	}
	return collection, nil
}

func (r *repository) DeleteCollection(ctx context.Context, id pgtype.UUID) error {
	if err := r.db.DeleteCollection(ctx, id); err != nil {
		return err
	}
	return nil
}


func (r *repository) CreateCollection(ctx context.Context, args db.CreateCollectionParams) (db.Collection, error) {
	collection, err := r.db.CreateCollection(ctx, args)
	if err != nil {
		return db.Collection{}, err
	}
	return collection, nil
} 