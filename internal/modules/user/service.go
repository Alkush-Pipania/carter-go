package user

import "context"

type Service interface {
	GetUserByID(context.Context, int) (User, error)
}

type service struct {
	repo Repository
}

func (s *service) GetUserByID(ctx context.Context, id int) (User, error) {
	return s.repo.GetUserByID(ctx, id)
}

// Exported constructor - the ONLY way to create a service
func NewService(repo Repository) Service {
	return &service{
		repo: repo,
	}
}
