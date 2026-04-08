package repository

import (
	"context"

	"github.com/bufunfaai/bufunfaai/apps/api/internal/modules/users/domain/entity"
)

type UserRepository interface {
	Create(ctx context.Context, user entity.User) (entity.User, error)
	GetByEmail(ctx context.Context, email string) (entity.User, error)
	GetByID(ctx context.Context, userID string) (entity.User, error)
	UpdateProfile(ctx context.Context, userID string, fullName string, phone string) (entity.User, error)
	UpdateLastLogin(ctx context.Context, userID string) error
}
