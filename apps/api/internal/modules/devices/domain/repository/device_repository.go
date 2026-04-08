package repository

import (
	"context"

	"github.com/bufunfaai/bufunfaai/apps/api/internal/modules/devices/domain/entity"
)

type DeviceRepository interface {
	Upsert(ctx context.Context, device entity.Device) (entity.Device, error)
	ListByUserID(ctx context.Context, userID string) ([]entity.Device, error)
	Delete(ctx context.Context, userID string, deviceID string) error
}
