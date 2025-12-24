package authentication

import (
	"context"
	"time"

	"github.com/Alkush-Pipania/carter-go/pkg/db"
	"github.com/Alkush-Pipania/carter-go/pkg/redis"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type Repository interface {
	GetUserByEmail(ctx context.Context, email string) (*db.User, error)
	CreateUser(ctx context.Context, params db.CreateUserParams) (*db.User, error)

	CreateSession(ctx context.Context, params db.CreateSessionParams) (*db.Session, error)
	RevokeSession(ctx context.Context, sessionID pgtype.UUID) error
	GetSession(ctx context.Context, sessionID pgtype.UUID) (*db.Session, error)
}

type service struct {
	repo  Repository
	redis *redis.Client
}

func NewService(repository Repository, redis *redis.Client) *service {
	return &service{
		repo:  repository,
		redis: redis,
	}
}

func (s *service) Login(ctx context.Context, email string, password string) (*LoginResponse, error) {
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

func (s *service) Logout(ctx context.Context, sessionID string) error {
	parsedID, err := parseUUID(sessionID)
	if err != nil {
		return err
	}
	return s.repo.RevokeSession(ctx, parsedID)
}

func (s *service) Register(ctx context.Context, password string, email string) error {
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

func (s *service) ValidateSession(ctx context.Context, sessionID string) (string, error) {
	userID, err := s.redis.GetSession(ctx, sessionID)
	if err == nil {
		return userID, nil
	}

	// fallback to db
	parsedID, err := parseUUID(sessionID)
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

// parseUUID converts a string to pgtype.UUID
func parseUUID(id string) (pgtype.UUID, error) {
	parsed, err := uuid.Parse(id)
	if err != nil {
		return pgtype.UUID{}, err
	}
	return pgtype.UUID{Bytes: parsed, Valid: true}, nil
}
