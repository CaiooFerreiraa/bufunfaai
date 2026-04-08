package presenter

import (
	"github.com/bufunfaai/bufunfaai/apps/api/internal/modules/devices/application/dto"
	"github.com/bufunfaai/bufunfaai/apps/api/internal/modules/devices/domain/entity"
)

func DeviceOutput(device entity.Device) dto.DeviceOutput {
	return dto.DeviceOutput{
		ID:         device.ID,
		DeviceName: device.DeviceName,
		Platform:   device.Platform,
		AppVersion: device.AppVersion,
		LastSeenAt: device.LastSeenAt.Format("2006-01-02T15:04:05Z07:00"),
		CreatedAt:  device.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}
