package authentication

import (
	"context"
	"time"

	"github.com/Alkush-Pipania/carter-go/pkg/db"
	"github.com/Alkush-Pipania/carter-go/pkg/redis"
	"github.com/Alkush-Pipania/carter-go/pkg/utils/toolchain"
	"github.com/jackc/pgx/v5/pgtype"
)

type Repository interface {
	GetUserByEmail(ctx context.Context, email string) (*db.User, error)
	CreateUser(ctx context.Context, params db.CreateUserParams) (*db.User, error)

	CreateSession(ctx context.Context, params db.CreateSessionParams) (*db.Session, error)
	RevokeSession(ctx context.Context, sessionID pgtype.UUID) error
	GetSession(ctx context.Context, sessionID pgtype.UUID) (*db.Session, error)
}

type Service struct {
	repo  Repository
	redis *redis.Client
}

func NewService(repository Repository, redis *redis.Client) *Service {
	return &Service{
		repo:  repository,
		redis: redis,
	}
}

func (s *Service) Login(ctx context.Context, email string, password string) (*LoginResponse, error) {
	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	if !comparePassword(user.PasswordHash, password) {
		return nil, ErrInvalidCredentials
	}
	session, err := s.repo.CreateSession(ctx, db.CreateSessionParams{
		UserID: user.ID,
	})
	if err != nil {
		return nil, err
	}

	_ = s.redis.SetSession(ctx, session.ID.String(), session.UserID.String(), time.Until(session.ExpiresAt.Time))

	return &LoginResponse{
		UserID:    user.ID.String(),
		SessionID: session.ID.String(),
	}, nil
}

func (s *Service) Logout(ctx context.Context, sessionID string) error {
	parsedID, err := toolchain.ParseUUID(sessionID)
	if err != nil {
		return err
	}

	// 1. Revoke in DB (source of truth)
	if err := s.repo.RevokeSession(ctx, parsedID); err != nil {
		return err
	}

	// 2. Best-effort delete from Redis (cache)
	_ = s.redis.DeleteSession(ctx, sessionID)

	return nil
}

func (s *Service) Register(ctx context.Context, password string, email string) error {
	hash, err := hashPassword(password)
	if err != nil {
		return err
	}

	_, err = s.repo.CreateUser(ctx, db.CreateUserParams{
		Email:        email,
		PasswordHash: hash,
	})
	return err
}

func (s *Service) ValidateSession(ctx context.Context, sessionID string) (string, error) {
	userID, err := s.redis.GetSession(ctx, sessionID)
	if err == nil {
		return userID, nil
	}

	// fallback to db
	parsedID, err := toolchain.ParseUUID(sessionID)
	if err != nil {
		return "", ErrUnauthorized
	}

	session, err := s.repo.GetSession(ctx, parsedID)
	if err != nil {
		return "", ErrUnauthorized
	}

	_ = s.redis.SetSession(
		ctx,
		session.ID.String(),
		session.UserID.String(),
		time.Until(session.ExpiresAt.Time),
	)

	return session.UserID.String(), nil
}
