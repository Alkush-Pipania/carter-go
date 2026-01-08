package user_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Alkush-Pipania/carter-go/internal/modules/user"
	"github.com/Alkush-Pipania/carter-go/pkg/db"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/assert"
)

// MockRepository is a mock implementation of user.Repository
type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) GetUserByID(ctx context.Context, id pgtype.UUID) (db.User, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(db.User), args.Error(1)
}

func (m *MockRepository) CreateUser(ctx context.Context, params db.CreateUserParams) (*db.User, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*db.User), args.Error(1)
}

func TestGetUserByID(t *testing.T) {
	testUUID := uuid.New()
	userID := testUUID.String()
	pgUUID := pgtype.UUID{Bytes: testUUID, Valid: true}

	expectedUser := db.User{
		ID:        pgUUID,
		Email:     "test@example.com",
		Username:  pgtype.Text{String: "testuser", Valid: true},
		ImageUrl:  pgtype.Text{String: "http://example.com/image.png", Valid: true},
		Verified:  true,
		CreatedAt: pgtype.Timestamptz{Time: time.Now(), Valid: true},
	}

	t.Run("Success", func(t *testing.T) {
		mockRepo := new(MockRepository)
		service := user.NewService(mockRepo)
		ctx := context.Background()

		mockRepo.On("GetUserByID", ctx, pgUUID).Return(expectedUser, nil)

		result, err := service.GetUserByID(ctx, userID)

		assert.NoError(t, err)
		assert.Equal(t, userID, result.ID)
		assert.Equal(t, "test@example.com", result.Email)
		mockRepo.AssertExpectations(t)
	})

	t.Run("NotFound", func(t *testing.T) {
		mockRepo := new(MockRepository)
		service := user.NewService(mockRepo)
		ctx := context.Background()

		mockRepo.On("GetUserByID", ctx, pgUUID).Return(db.User{}, errors.New("user not found"))

		_, err := service.GetUserByID(ctx, userID)

		assert.Error(t, err)
		assert.Equal(t, "user not found", err.Error())
		mockRepo.AssertExpectations(t)
	})
}
