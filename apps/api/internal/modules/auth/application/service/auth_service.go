package service

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"

	authdto "github.com/bufunfaai/bufunfaai/apps/api/internal/modules/auth/application/dto"
	authentity "github.com/bufunfaai/bufunfaai/apps/api/internal/modules/auth/domain/entity"
	authrepository "github.com/bufunfaai/bufunfaai/apps/api/internal/modules/auth/domain/repository"
	devicedto "github.com/bufunfaai/bufunfaai/apps/api/internal/modules/devices/application/dto"
	deviceentity "github.com/bufunfaai/bufunfaai/apps/api/internal/modules/devices/domain/entity"
	devicerepository "github.com/bufunfaai/bufunfaai/apps/api/internal/modules/devices/domain/repository"
	userservice "github.com/bufunfaai/bufunfaai/apps/api/internal/modules/users/application/service"
	userentity "github.com/bufunfaai/bufunfaai/apps/api/internal/modules/users/domain/entity"
	userrepository "github.com/bufunfaai/bufunfaai/apps/api/internal/modules/users/domain/repository"
	sharederrors "github.com/bufunfaai/bufunfaai/apps/api/internal/shared/errors"
)

type PasswordHasher interface {
	HashPassword(password string) (string, error)
	ComparePassword(encodedHash string, password string) error
}

type RefreshTokenManager interface {
	GenerateToken() (string, error)
	HashToken(rawToken string) string
}

type AccessTokenIssuer interface {
	IssueAccessToken(user userentity.User, sessionID string, now time.Time) (string, time.Time, error)
}

type AuthService struct {
	userRepository         userrepository.UserRepository
	refreshTokenRepository authrepository.RefreshTokenRepository
	deviceRepository       devicerepository.DeviceRepository
	passwordHasher         PasswordHasher
	refreshTokenManager    RefreshTokenManager
	accessTokenIssuer      AccessTokenIssuer
	refreshTokenTTL        time.Duration
	now                    func() time.Time
}

type RegisterInput struct {
	FullName  string
	Email     string
	Password  string
	Phone     string
	Device    *devicedto.DeviceMetadata
	IPAddress string
	UserAgent string
}

type LoginInput struct {
	Email     string
	Password  string
	Device    *devicedto.DeviceMetadata
	IPAddress string
	UserAgent string
}

type RefreshInput struct {
	RefreshToken string
	IPAddress    string
	UserAgent    string
}

func NewAuthService(
	userRepository userrepository.UserRepository,
	refreshTokenRepository authrepository.RefreshTokenRepository,
	deviceRepository devicerepository.DeviceRepository,
	passwordHasher PasswordHasher,
	refreshTokenManager RefreshTokenManager,
	accessTokenIssuer AccessTokenIssuer,
	refreshTokenTTL time.Duration,
) *AuthService {
	return &AuthService{
		userRepository:         userRepository,
		refreshTokenRepository: refreshTokenRepository,
		deviceRepository:       deviceRepository,
		passwordHasher:         passwordHasher,
		refreshTokenManager:    refreshTokenManager,
		accessTokenIssuer:      accessTokenIssuer,
		refreshTokenTTL:        refreshTokenTTL,
		now:                    time.Now,
	}
}

func (service *AuthService) Register(ctx context.Context, input RegisterInput) (authdto.AuthResult, *sharederrors.AppError) {
	email := strings.ToLower(strings.TrimSpace(input.Email))
	if _, err := service.userRepository.GetByEmail(ctx, email); err == nil {
		return authdto.AuthResult{}, sharederrors.New("EMAIL_ALREADY_IN_USE", "email ja cadastrado", 409)
	} else if !errors.Is(err, userservice.ErrUserNotFound) {
		return authdto.AuthResult{}, sharederrors.Wrap("INTERNAL_ERROR", "erro interno ao consultar usuario", 500, err)
	}

	passwordHash, err := service.passwordHasher.HashPassword(input.Password)
	if err != nil {
		return authdto.AuthResult{}, sharederrors.Wrap("INTERNAL_ERROR", "erro interno ao proteger senha", 500, err)
	}

	now := service.now().UTC()
	user := userentity.User{
		ID:           uuid.NewString(),
		FullName:     strings.TrimSpace(input.FullName),
		Email:        email,
		Phone:        strings.TrimSpace(input.Phone),
		PasswordHash: passwordHash,
		Status:       "active",
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	createdUser, err := service.userRepository.Create(ctx, user)
	if err != nil {
		return authdto.AuthResult{}, sharederrors.Wrap("INTERNAL_ERROR", "erro interno ao criar usuario", 500, err)
	}

	return service.issueSession(ctx, createdUser, input.Device, input.IPAddress, input.UserAgent, now)
}

func (service *AuthService) Login(ctx context.Context, input LoginInput) (authdto.AuthResult, *sharederrors.AppError) {
	user, err := service.userRepository.GetByEmail(ctx, input.Email)
	if err != nil {
		if errors.Is(err, userservice.ErrUserNotFound) {
			return authdto.AuthResult{}, sharederrors.New("INVALID_CREDENTIALS", "credenciais invalidas", 401)
		}

		return authdto.AuthResult{}, sharederrors.Wrap("INTERNAL_ERROR", "erro interno ao consultar usuario", 500, err)
	}

	if err := service.passwordHasher.ComparePassword(user.PasswordHash, input.Password); err != nil {
		return authdto.AuthResult{}, sharederrors.New("INVALID_CREDENTIALS", "credenciais invalidas", 401)
	}

	if err := service.userRepository.UpdateLastLogin(ctx, user.ID); err != nil {
		return authdto.AuthResult{}, sharederrors.Wrap("INTERNAL_ERROR", "erro interno ao atualizar ultimo login", 500, err)
	}

	return service.issueSession(ctx, user, input.Device, input.IPAddress, input.UserAgent, service.now().UTC())
}

func (service *AuthService) Refresh(ctx context.Context, input RefreshInput) (authdto.AuthResult, *sharederrors.AppError) {
	tokenHash := service.refreshTokenManager.HashToken(input.RefreshToken)
	storedToken, err := service.refreshTokenRepository.GetByTokenHash(ctx, tokenHash)
	if err != nil {
		if errors.Is(err, ErrRefreshTokenNotFound) {
			return authdto.AuthResult{}, sharederrors.New("INVALID_REFRESH_TOKEN", "refresh token invalido ou expirado", 401)
		}

		return authdto.AuthResult{}, sharederrors.Wrap("INTERNAL_ERROR", "erro interno ao consultar refresh token", 500, err)
	}

	now := service.now().UTC()
	if !storedToken.IsActive(now) {
		return authdto.AuthResult{}, sharederrors.New("INVALID_REFRESH_TOKEN", "refresh token invalido ou expirado", 401)
	}

	user, err := service.userRepository.GetByID(ctx, storedToken.UserID)
	if err != nil {
		if errors.Is(err, userservice.ErrUserNotFound) {
			return authdto.AuthResult{}, sharederrors.New("USER_NOT_FOUND", "usuario nao encontrado", 404)
		}

		return authdto.AuthResult{}, sharederrors.Wrap("INTERNAL_ERROR", "erro interno ao consultar usuario", 500, err)
	}

	if err := service.refreshTokenRepository.Revoke(ctx, storedToken.ID, now); err != nil {
		return authdto.AuthResult{}, sharederrors.Wrap("INTERNAL_ERROR", "erro interno ao revogar refresh token", 500, err)
	}

	var device *devicedto.DeviceMetadata
	if storedToken.DeviceID != nil {
		device = &devicedto.DeviceMetadata{}
	}

	return service.issueSession(ctx, user, device, input.IPAddress, input.UserAgent, now)
}

func (service *AuthService) Logout(ctx context.Context, refreshToken string) *sharederrors.AppError {
	tokenHash := service.refreshTokenManager.HashToken(refreshToken)
	storedToken, err := service.refreshTokenRepository.GetByTokenHash(ctx, tokenHash)
	if err != nil {
		if errors.Is(err, ErrRefreshTokenNotFound) {
			return sharederrors.New("INVALID_REFRESH_TOKEN", "refresh token invalido ou expirado", 401)
		}

		return sharederrors.Wrap("INTERNAL_ERROR", "erro interno ao consultar refresh token", 500, err)
	}

	if err := service.refreshTokenRepository.Revoke(ctx, storedToken.ID, service.now().UTC()); err != nil {
		return sharederrors.Wrap("INTERNAL_ERROR", "erro interno ao revogar refresh token", 500, err)
	}

	return nil
}

func (service *AuthService) LogoutAll(ctx context.Context, userID string) *sharederrors.AppError {
	if err := service.refreshTokenRepository.RevokeByUserID(ctx, userID, service.now().UTC()); err != nil {
		return sharederrors.Wrap("INTERNAL_ERROR", "erro interno ao revogar sessoes", 500, err)
	}

	return nil
}

func (service *AuthService) issueSession(
	ctx context.Context,
	user userentity.User,
	device *devicedto.DeviceMetadata,
	ipAddress string,
	userAgent string,
	now time.Time,
) (authdto.AuthResult, *sharederrors.AppError) {
	var deviceID *string
	if device != nil {
		upsertedDevice, err := service.deviceRepository.Upsert(ctx, deviceentity.Device{
			ID:              uuid.NewString(),
			UserID:          user.ID,
			DeviceName:      strings.TrimSpace(device.DeviceName),
			Platform:        strings.TrimSpace(device.Platform),
			AppVersion:      strings.TrimSpace(device.AppVersion),
			FingerprintHash: strings.TrimSpace(device.Fingerprint),
			LastSeenAt:      now,
			CreatedAt:       now,
		})
		if err != nil {
			return authdto.AuthResult{}, sharederrors.Wrap("INTERNAL_ERROR", "erro interno ao registrar dispositivo", 500, err)
		}

		deviceID = &upsertedDevice.ID
	}

	rawRefreshToken, err := service.refreshTokenManager.GenerateToken()
	if err != nil {
		return authdto.AuthResult{}, sharederrors.Wrap("INTERNAL_ERROR", "erro interno ao gerar refresh token", 500, err)
	}

	refreshToken := authentity.RefreshToken{
		ID:        uuid.NewString(),
		UserID:    user.ID,
		TokenHash: service.refreshTokenManager.HashToken(rawRefreshToken),
		DeviceID:  deviceID,
		ExpiresAt: now.Add(service.refreshTokenTTL),
		CreatedAt: now,
		IPAddress: ipAddress,
		UserAgent: userAgent,
	}

	if err := service.refreshTokenRepository.Create(ctx, refreshToken); err != nil {
		return authdto.AuthResult{}, sharederrors.Wrap("INTERNAL_ERROR", "erro interno ao persistir refresh token", 500, err)
	}

	accessToken, expiresAt, err := service.accessTokenIssuer.IssueAccessToken(user, refreshToken.ID, now)
	if err != nil {
		return authdto.AuthResult{}, sharederrors.Wrap("INTERNAL_ERROR", "erro interno ao gerar access token", 500, err)
	}

	return authdto.AuthResult{
		User: authdto.AuthUserOutput{
			ID:       user.ID,
			FullName: user.FullName,
			Email:    user.Email,
			Status:   user.Status,
		},
		Session: authdto.AuthTokensOutput{
			AccessToken:  accessToken,
			RefreshToken: rawRefreshToken,
			ExpiresAt:    expiresAt.Format(time.RFC3339),
		},
	}, nil
}

var ErrRefreshTokenNotFound = errors.New("refresh token not found")
