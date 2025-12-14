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

func (r *repository) GetUserByID(ctx context.Context, id string) (User, error) {
	var uuid pgtype.UUID
	if err := uuid.Scan(id); err != nil {
		return User{}, err
	}
	user, err := r.db.GetUserById(ctx, uuid)
	if err != nil {
		return User{}, err
	}
	return User{
		ID:        user.ID.String(),
		Email:     user.Email,
		Username:  user.Username.String,
		Image:     user.ImageUrl.String,
		Password:  user.PasswordHash,
		Verified:  user.Verified,
		CreatedAt: user.CreatedAt.Time,
	}, nil
}
