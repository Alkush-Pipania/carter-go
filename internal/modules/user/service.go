package user

import (
	"context"

	"github.com/Alkush-Pipania/carter-go/pkg/db"
	"github.com/jackc/pgx/v5/pgtype"
)

type Repository interface {
	GetUserByID(context.Context, pgtype.UUID) (db.User, error)
	CreateUser(context.Context, InputCreateUser) error
}

type service struct {
	repo Repository
}

func (s *service) GetUserByID(ctx context.Context, id string) (User, error) {
	var uuid pgtype.UUID
	if err := uuid.Scan(id); err != nil {
		return User{}, err
	}
	user, err := s.repo.GetUserByID(ctx, uuid)
	if err != nil {
		return User{}, err
	}
	return User{
		ID:        user.ID.String(),
		Email:     user.Email,
		Username:  user.Username.String,
		Image:     user.ImageUrl.String,
		Verified:  user.Verified,
		CreatedAt: user.CreatedAt.Time,
	}, nil
}

func (s *service) CreateUser(ctx context.Context, user InputCreateUser) error {
	return s.repo.CreateUser(ctx, user)
}

// Exported constructor - the ONLY way to create a service
func NewService(repo Repository) Service {
	return &service{
		repo: repo,
	}
}
