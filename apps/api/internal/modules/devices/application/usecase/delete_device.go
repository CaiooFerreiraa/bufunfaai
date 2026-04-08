package usecase

import (
	"context"

	"github.com/bufunfaai/bufunfaai/apps/api/internal/modules/devices/application/service"
	sharederrors "github.com/bufunfaai/bufunfaai/apps/api/internal/shared/errors"
)

type DeleteDeviceUseCase struct {
	deviceService *service.DeviceService
}

func NewDeleteDeviceUseCase(deviceService *service.DeviceService) *DeleteDeviceUseCase {
	return &DeleteDeviceUseCase{deviceService: deviceService}
}

func (useCase *DeleteDeviceUseCase) Execute(ctx context.Context, userID string, deviceID string) *sharederrors.AppError {
	return useCase.deviceService.Delete(ctx, userID, deviceID)
}
