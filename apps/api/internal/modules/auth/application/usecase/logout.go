package usecase

import (
	"context"

	authservice "github.com/bufunfaai/bufunfaai/apps/api/internal/modules/auth/application/service"
	sharederrors "github.com/bufunfaai/bufunfaai/apps/api/internal/shared/errors"
)

type LogoutUseCase struct {
	authService *authservice.AuthService
}

func NewLogoutUseCase(authService *authservice.AuthService) *LogoutUseCase {
	return &LogoutUseCase{authService: authService}
}

func (useCase *LogoutUseCase) Execute(ctx context.Context, refreshToken string) *sharederrors.AppError {
	return useCase.authService.Logout(ctx, refreshToken)
}
