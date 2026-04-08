package dto

import (
	devicedto "github.com/bufunfaai/bufunfaai/apps/api/internal/modules/devices/application/dto"
)

type RegisterRequest struct {
	FullName string                    `json:"full_name" validate:"required,min=3,max=120"`
	Email    string                    `json:"email" validate:"required,email,max=255"`
	Password string                    `json:"password" validate:"required,min=8,max=72"`
	Phone    string                    `json:"phone" validate:"omitempty,min=8,max=20"`
	Device   *devicedto.DeviceMetadata `json:"device"`
}

type LoginRequest struct {
	Email    string                    `json:"email" validate:"required,email,max=255"`
	Password string                    `json:"password" validate:"required,min=8,max=72"`
	Device   *devicedto.DeviceMetadata `json:"device"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required,min=20"`
}

type LogoutRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required,min=20"`
}

type AuthTokensOutput struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresAt    string `json:"expires_at"`
}

type AuthUserOutput struct {
	ID       string `json:"id"`
	FullName string `json:"full_name"`
	Email    string `json:"email"`
	Status   string `json:"status"`
}

type AuthResult struct {
	User    AuthUserOutput   `json:"user"`
	Session AuthTokensOutput `json:"session"`
}
