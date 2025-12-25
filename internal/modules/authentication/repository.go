package authentication

import (
	"context"

	"github.com/Alkush-Pipania/carter-go/pkg/db"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type repository struct {
	db *db.Queries
}

func NewRepository(db *db.Queries) *repository {
	return &repository{
		db: db,
	}
}

// -------------------- User Methods --------------------

func (r *repository) GetUserByEmail(ctx context.Context, email string) (*db.User, error) {
	user, err := r.db.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	return &user, nil
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

// -------------------- Session Methods --------------------

func (r *repository) CreateSession(ctx context.Context, params db.CreateSessionParams) (*db.Session, error) {
	// Generate UUID if not provided
	if !params.ID.Valid {
		newID := uuid.New()
		params.ID = pgtype.UUID{Bytes: newID, Valid: true}
	}

	session, err := r.db.CreateSession(ctx, params)
	if err != nil {
		return nil, err
	}
	return &session, nil
}

func (r *repository) GetSession(ctx context.Context, sessionID pgtype.UUID) (*db.Session, error) {
	session, err := r.db.GetSessionByID(ctx, sessionID)
	if err != nil {
		return nil, err
	}
	return &session, nil
}

func (r *repository) RevokeSession(ctx context.Context, sessionID pgtype.UUID) error {
	return r.db.RevokeSession(ctx, sessionID)
}
