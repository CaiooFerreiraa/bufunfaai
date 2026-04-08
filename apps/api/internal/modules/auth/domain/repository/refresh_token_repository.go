package repository

import (
	"context"
	"time"

	"github.com/bufunfaai/bufunfaai/apps/api/internal/modules/auth/domain/entity"
)

type RefreshTokenRepository interface {
	Create(ctx context.Context, token entity.RefreshToken) error
	GetByTokenHash(ctx context.Context, tokenHash string) (entity.RefreshToken, error)
	Revoke(ctx context.Context, tokenID string, revokedAt time.Time) error
	RevokeByUserID(ctx context.Context, userID string, revokedAt time.Time) error
	RevokeByUserIDAndDeviceID(ctx context.Context, userID string, deviceID string, revokedAt time.Time) error
}
