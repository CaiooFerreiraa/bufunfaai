package usecase

import (
	"context"

	"github.com/bufunfaai/bufunfaai/apps/api/internal/modules/users/application/service"
	"github.com/bufunfaai/bufunfaai/apps/api/internal/modules/users/domain/entity"
	sharederrors "github.com/bufunfaai/bufunfaai/apps/api/internal/shared/errors"
)

type GetCurrentUserUseCase struct {
	userService *service.UserService
}

func NewGetCurrentUserUseCase(userService *service.UserService) *GetCurrentUserUseCase {
	return &GetCurrentUserUseCase{userService: userService}
}

func (useCase *GetCurrentUserUseCase) Execute(ctx context.Context, userID string) (entity.User, *sharederrors.AppError) {
	return useCase.userService.GetCurrentUser(ctx, userID)
}
