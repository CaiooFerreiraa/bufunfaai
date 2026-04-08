package usecase

import (
	"context"

	authdto "github.com/bufunfaai/bufunfaai/apps/api/internal/modules/auth/application/dto"
	authservice "github.com/bufunfaai/bufunfaai/apps/api/internal/modules/auth/application/service"
	sharederrors "github.com/bufunfaai/bufunfaai/apps/api/internal/shared/errors"
)

type RefreshUseCase struct {
	authService *authservice.AuthService
}

func NewRefreshUseCase(authService *authservice.AuthService) *RefreshUseCase {
	return &RefreshUseCase{authService: authService}
}

func (useCase *RefreshUseCase) Execute(ctx context.Context, request authdto.RefreshRequest, ipAddress string, userAgent string) (authdto.AuthResult, *sharederrors.AppError) {
	return useCase.authService.Refresh(ctx, authservice.RefreshInput{
		RefreshToken: request.RefreshToken,
		IPAddress:    ipAddress,
		UserAgent:    userAgent,
	})
}
