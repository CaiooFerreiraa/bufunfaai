package usecase

import (
	"context"

	userservice "github.com/bufunfaai/bufunfaai/apps/api/internal/modules/users/application/service"
	"github.com/bufunfaai/bufunfaai/apps/api/internal/modules/users/domain/entity"
	sharederrors "github.com/bufunfaai/bufunfaai/apps/api/internal/shared/errors"
)

type UpdateCurrentUserUseCase struct {
	userService *userservice.UserService
}

func NewUpdateCurrentUserUseCase(userService *userservice.UserService) *UpdateCurrentUserUseCase {
	return &UpdateCurrentUserUseCase{userService: userService}
}

func (useCase *UpdateCurrentUserUseCase) Execute(ctx context.Context, userID string, fullName string, phone string) (entity.User, *sharederrors.AppError) {
	return useCase.userService.UpdateCurrentUser(ctx, userID, fullName, phone)
}
