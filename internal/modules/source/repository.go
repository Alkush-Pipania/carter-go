package source

import (
	"context"

	"github.com/Alkush-Pipania/carter-go/pkg/db"
	"github.com/jackc/pgx/v5/pgtype"
)

type repository struct {
	db *db.Queries
}

func NewRepository(db *db.Queries) *repository {
	return &repository{db: db}
}

func (r *repository) GetSourcesByCollectionID(ctx context.Context, collectionID pgtype.UUID) ([]db.Source, error) {
	return r.db.GetSourcesByCollectionID(ctx, collectionID)
}

func (r *repository) CreateSource(ctx context.Context, args db.CreateSourceParams) (db.Source, error) {
	return r.db.CreateSource(ctx, args)
}

func (r *repository) UpdateSourceStatus(ctx context.Context, args db.UpdateSourceStatusParams) (db.Source, error) {
	return r.db.UpdateSourceStatus(ctx, args)
}
