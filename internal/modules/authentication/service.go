package authentication

import (
	"context"
	"time"

	"github.com/Alkush-Pipania/carter-go/pkg/db"
	"github.com/Alkush-Pipania/carter-go/pkg/logger"
	"github.com/Alkush-Pipania/carter-go/pkg/redis"
	"github.com/Alkush-Pipania/carter-go/pkg/utils/toolchain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"go.uber.org/zap"
)

type Repository interface {
	GetUserByEmail(ctx context.Context, email string) (*db.User, error)
	CreateSession(ctx context.Context, params db.CreateSessionParams) (*db.Session, error)
	RevokeSession(ctx context.Context, sessionID pgtype.UUID) error
	GetSession(ctx context.Context, sessionID pgtype.UUID) (*db.Session, error)
}

type Service struct {
	repo        Repository
	redis       *redis.Client
	UserService UserService
}

type UserService interface {
	CreateUser(context.Context, string, string) error
}

func NewService(repository Repository, redis *redis.Client, userService UserService) *Service {
	return &Service{
		repo:        repository,
		redis:       redis,
		UserService: userService,
	}
}

func (s *Service) Login(ctx context.Context, email string, password string) (*LoginResponse, error) {
	logger.Debug("Login attempt", zap.String("email", email))

	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		logger.Warn("Login failed: user not found", zap.String("email", email))
		return nil, ErrInvalidCredentials
	}

	if !comparePassword(user.PasswordHash, password) {
		logger.Warn("Login failed: invalid password", zap.String("email", email))
		return nil, ErrInvalidCredentials
	}

	sessionID := uuid.New()
	expiresAt := time.Now().Add(30 * 24 * time.Hour) // 30 days

	session, err := s.repo.CreateSession(ctx, db.CreateSessionParams{
		ID:        pgtype.UUID{Bytes: sessionID, Valid: true},
		UserID:    user.ID,
		ExpiresAt: pgtype.Timestamptz{Time: expiresAt, Valid: true},
	})
	if err != nil {
		logger.Error("Login failed: could not create session", zap.String("email", email), zap.Error(err))
		return nil, err
	}

	if err := s.redis.SetSession(ctx, session.ID.String(), session.UserID.String(), time.Until(session.ExpiresAt.Time)); err != nil {
		logger.Warn("Failed to cache session in Redis", zap.String("session_id", session.ID.String()), zap.Error(err))
	}

	logger.Info("User logged in successfully",
		zap.String("user_id", user.ID.String()),
		zap.String("session_id", session.ID.String()),
	)

	return &LoginResponse{
		UserID:    user.ID.String(),
		SessionID: session.ID.String(),
	}, nil
}

func (s *Service) Logout(ctx context.Context, sessionID string) error {
	logger.Debug("Logout attempt", zap.String("session_id", sessionID))

	parsedID, err := toolchain.ParseUUID(sessionID)
	if err != nil {
		logger.Warn("Logout failed: invalid session ID format", zap.String("session_id", sessionID), zap.Error(err))
		return err
	}

	// 1. Revoke in DB (source of truth)
	if err := s.repo.RevokeSession(ctx, parsedID); err != nil {
		logger.Error("Logout failed: could not revoke session in DB", zap.String("session_id", sessionID), zap.Error(err))
		return err
	}

	// 2. Best-effort delete from Redis (cache)
	if err := s.redis.DeleteSession(ctx, sessionID); err != nil {
		logger.Warn("Failed to delete session from Redis cache", zap.String("session_id", sessionID), zap.Error(err))
	}

	logger.Info("User logged out successfully", zap.String("session_id", sessionID))
	return nil
}

func (s *Service) Register(ctx context.Context, password string, email string) error {
	logger.Debug("Registration attempt", zap.String("email", email))

	hash, err := hashPassword(password)
	if err != nil {
		logger.Error("Registration failed: could not hash password", zap.String("email", email), zap.Error(err))
		return err
	}

	err = s.UserService.CreateUser(ctx, email, hash)
	if err != nil {
		logger.Error("Registration failed: could not create user", zap.String("email", email), zap.Error(err))
		return err
	}

	logger.Info("User registered successfully", zap.String("email", email))
	return nil
}

func (s *Service) ValidateSession(ctx context.Context, sessionID string) (string, error) {
	logger.Debug("Validating session", zap.String("session_id", sessionID))

	userID, err := s.redis.GetSession(ctx, sessionID)
	if err == nil {
		logger.Debug("Session found in Redis cache", zap.String("session_id", sessionID), zap.String("user_id", userID))
		return userID, nil
	}

	logger.Debug("Session not in Redis, falling back to DB", zap.String("session_id", sessionID))

	// fallback to db
	parsedID, err := toolchain.ParseUUID(sessionID)
	if err != nil {
		logger.Warn("Session validation failed: invalid session ID format", zap.String("session_id", sessionID), zap.Error(err))
		return "", ErrUnauthorized
	}

	session, err := s.repo.GetSession(ctx, parsedID)
	if err != nil {
		logger.Warn("Session validation failed: session not found in DB", zap.String("session_id", sessionID), zap.Error(err))
		return "", ErrUnauthorized
	}

	if err := s.redis.SetSession(
		ctx,
		session.ID.String(),
		session.UserID.String(),
		time.Until(session.ExpiresAt.Time),
	); err != nil {
		logger.Warn("Failed to cache session in Redis after DB lookup", zap.String("session_id", sessionID), zap.Error(err))
	}

	logger.Debug("Session validated from DB and cached",
		zap.String("session_id", sessionID),
		zap.String("user_id", session.UserID.String()),
	)

	return session.UserID.String(), nil
}
