package usecase

import (
	"context"

	authservice "github.com/bufunfaai/bufunfaai/apps/api/internal/modules/auth/application/service"
	sharederrors "github.com/bufunfaai/bufunfaai/apps/api/internal/shared/errors"
)

type LogoutAllUseCase struct {
	authService *authservice.AuthService
}

func NewLogoutAllUseCase(authService *authservice.AuthService) *LogoutAllUseCase {
	return &LogoutAllUseCase{authService: authService}
}

func (useCase *LogoutAllUseCase) Execute(ctx context.Context, userID string) *sharederrors.AppError {
	return useCase.authService.LogoutAll(ctx, userID)
}
