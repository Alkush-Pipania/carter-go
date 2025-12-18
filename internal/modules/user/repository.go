package user

import (
	"context"

	"github.com/Alkush-Pipania/carter-go/pkg/db"
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

func (r *repository) CreateUser(ctx context.Context, user InputCreateUser) error {
	_, err := r.db.CreateUser(ctx, db.CreateUserParams{
		Email:        user.Email,
		Username:     user.Username,
		PasswordHash: user.PasswordHash,
	})
	if err != nil {
		return err
	}
	return nil
}
