package user

import (
	"context"

	"github.com/Alkush-Pipania/carter-go/pkg/db"
)

type Repository interface {
	GetUserByID(context.Context, int) (User, error)
}

type repository struct {
	db *db.Queries
}

func NewRepository(db *db.Queries) Repository {
	return &repository{db: db}
}

func (r *repository) GetUserByID(ctx context.Context, id int) (User, error) {
	user, _ := r.db.GetUserByID(ctx, int32(id))
	return User{
		ID:        int(user.ID),
		Email:     user.Email,
		Username:  user.Username.String,
		Image:     user.Image.String,
		Password:  user.Password,
		SecretKey: user.Secretkey.String(),
		Verified:  user.Verified,
		CreatedAt: user.CreatedAt.Time,
		UpdatedAt: user.UpdatedAt.Time,
	}, nil
}
