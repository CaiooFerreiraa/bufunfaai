package usecase

import (
	"context"

	authdto "github.com/bufunfaai/bufunfaai/apps/api/internal/modules/auth/application/dto"
	authservice "github.com/bufunfaai/bufunfaai/apps/api/internal/modules/auth/application/service"
	sharederrors "github.com/bufunfaai/bufunfaai/apps/api/internal/shared/errors"
)

type LoginUseCase struct {
	authService *authservice.AuthService
}

func NewLoginUseCase(authService *authservice.AuthService) *LoginUseCase {
	return &LoginUseCase{authService: authService}
}

func (useCase *LoginUseCase) Execute(ctx context.Context, request authdto.LoginRequest, ipAddress string, userAgent string) (authdto.AuthResult, *sharederrors.AppError) {
	return useCase.authService.Login(ctx, authservice.LoginInput{
		Email:     request.Email,
		Password:  request.Password,
		Device:    request.Device,
		IPAddress: ipAddress,
		UserAgent: userAgent,
	})
}
