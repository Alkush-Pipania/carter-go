package user

import (
	"context"

	"github.com/Alkush-Pipania/carter-go/pkg/db"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type repository struct {
	db *db.Queries
}

func NewRepository(db *db.Queries) Repository {
	return &repository{db: db}
}

func (r *repository) GetUserByID(ctx context.Context, id pgtype.UUID) (db.User, error) {
	user, err := r.db.GetUserById(ctx, id)
	if err != nil {
		return db.User{}, err
	}
	return user, nil
}

func (r *repository) CreateUser(ctx context.Context, params db.CreateUserParams) (*db.User, error) {
	// Generate UUID if not provided
	if !params.ID.Valid {
		newID := uuid.New()
		params.ID = pgtype.UUID{Bytes: newID, Valid: true}
	}

	user, err := r.db.CreateUser(ctx, params)
	if err != nil {
		return nil, err
	}
	return &user, nil
}


