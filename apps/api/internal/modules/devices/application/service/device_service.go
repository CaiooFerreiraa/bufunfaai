package service

import (
	"context"

	"github.com/bufunfaai/bufunfaai/apps/api/internal/modules/devices/domain/entity"
	devicerepository "github.com/bufunfaai/bufunfaai/apps/api/internal/modules/devices/domain/repository"
	sharederrors "github.com/bufunfaai/bufunfaai/apps/api/internal/shared/errors"
)

type DeviceService struct {
	deviceRepository devicerepository.DeviceRepository
}

func NewDeviceService(deviceRepository devicerepository.DeviceRepository) *DeviceService {
	return &DeviceService{deviceRepository: deviceRepository}
}

func (service *DeviceService) List(ctx context.Context, userID string) ([]entity.Device, *sharederrors.AppError) {
	devices, err := service.deviceRepository.ListByUserID(ctx, userID)
	if err != nil {
		return nil, sharederrors.Wrap("INTERNAL_ERROR", "erro interno ao listar dispositivos", 500, err)
	}

	return devices, nil
}

func (service *DeviceService) Delete(ctx context.Context, userID string, deviceID string) *sharederrors.AppError {
	if err := service.deviceRepository.Delete(ctx, userID, deviceID); err != nil {
		if err == ErrDeviceNotFound {
			return sharederrors.New("DEVICE_NOT_FOUND", "dispositivo nao encontrado", 404)
		}

		return sharederrors.Wrap("INTERNAL_ERROR", "erro interno ao remover dispositivo", 500, err)
	}

	return nil
}
