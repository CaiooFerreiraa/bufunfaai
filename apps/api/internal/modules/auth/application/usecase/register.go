package usecase

import (
	"context"

	authdto "github.com/bufunfaai/bufunfaai/apps/api/internal/modules/auth/application/dto"
	authservice "github.com/bufunfaai/bufunfaai/apps/api/internal/modules/auth/application/service"
	sharederrors "github.com/bufunfaai/bufunfaai/apps/api/internal/shared/errors"
)

type RegisterUseCase struct {
	authService *authservice.AuthService
}

func NewRegisterUseCase(authService *authservice.AuthService) *RegisterUseCase {
	return &RegisterUseCase{authService: authService}
}

func (useCase *RegisterUseCase) Execute(ctx context.Context, request authdto.RegisterRequest, ipAddress string, userAgent string) (authdto.AuthResult, *sharederrors.AppError) {
	return useCase.authService.Register(ctx, authservice.RegisterInput{
		FullName:  request.FullName,
		Email:     request.Email,
		Password:  request.Password,
		Phone:     request.Phone,
		Device:    request.Device,
		IPAddress: ipAddress,
		UserAgent: userAgent,
	})
}
