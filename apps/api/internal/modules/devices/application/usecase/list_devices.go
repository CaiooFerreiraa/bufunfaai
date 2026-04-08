package usecase

import (
	"context"

	"github.com/bufunfaai/bufunfaai/apps/api/internal/modules/devices/application/service"
	"github.com/bufunfaai/bufunfaai/apps/api/internal/modules/devices/domain/entity"
	sharederrors "github.com/bufunfaai/bufunfaai/apps/api/internal/shared/errors"
)

type ListDevicesUseCase struct {
	deviceService *service.DeviceService
}

func NewListDevicesUseCase(deviceService *service.DeviceService) *ListDevicesUseCase {
	return &ListDevicesUseCase{deviceService: deviceService}
}

func (useCase *ListDevicesUseCase) Execute(ctx context.Context, userID string) ([]entity.Device, *sharederrors.AppError) {
	return useCase.deviceService.List(ctx, userID)
}
