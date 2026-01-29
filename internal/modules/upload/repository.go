package upload

import (
	"context"

	"github.com/Alkush-Pipania/carter-go/pkg/db"
	"github.com/jackc/pgx/v5/pgtype"
)

type repository struct {
	dbQuerier *db.Queries
}

func NewRepository(db *db.Queries) *repository {
	return &repository{
		dbQuerier: db,
	}
}

func (r *repository) CreateSourceWithS3Key(ctx context.Context, args db.CreateSourceWithS3KeyParams) (db.Source, error) {
	return r.dbQuerier.CreateSourceWithS3Key(ctx, args)
}

func (r *repository) UpdateSourceStatus(ctx context.Context, args db.UpdateSourceStatusParams) (db.Source, error) {
	return r.dbQuerier.UpdateSourceStatus(ctx, args)
}

func (r *repository) GetSourceByID(ctx context.Context, id pgtype.UUID) (db.Source, error) {
	return r.dbQuerier.GetSourceByID(ctx, id)
}
