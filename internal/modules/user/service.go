package user

import "context"

type Repository interface {
	GetUserByID(context.Context, string) (User, error)
}

type service struct {
	repo Repository
}

func (s *service) GetUserByID(ctx context.Context, id string) (User, error) {
	return s.repo.GetUserByID(ctx, id)
}

// Exported constructor - the ONLY way to create a service
func NewService(repo Repository) Service {
	return &service{
		repo: repo,
	}
}
