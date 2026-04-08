package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/bufunfaai/bufunfaai/apps/api/internal/platform/config"
)

type Claims struct {
	Email     string `json:"email"`
	Role      string `json:"role"`
	SessionID string `json:"sid"`
	jwt.RegisteredClaims
}

type JWTService struct {
	secret   []byte
	issuer   string
	audience string
	ttl      time.Duration
}

type IssueTokenInput struct {
	UserID    string
	Email     string
	SessionID string
	Role      string
	Now       time.Time
}

func NewJWTService(cfg config.Config) *JWTService {
	return &JWTService{
		secret:   []byte(cfg.AccessTokenSecret),
		issuer:   cfg.AccessTokenIssuer,
		audience: cfg.AccessTokenAudience,
		ttl:      cfg.AccessTokenTTL,
	}
}

func (service *JWTService) Issue(input IssueTokenInput) (string, time.Time, error) {
	expiresAt := input.Now.Add(service.ttl)
	claims := Claims{
		Email:     input.Email,
		Role:      input.Role,
		SessionID: input.SessionID,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   input.UserID,
			Audience:  jwt.ClaimStrings{service.audience},
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(input.Now),
			Issuer:    service.issuer,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(service.secret)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("sign jwt: %w", err)
	}

	return signedToken, expiresAt, nil
}

func (service *JWTService) Parse(accessToken string) (Claims, error) {
	claims := Claims{}
	parsedToken, err := jwt.ParseWithClaims(accessToken, &claims, func(token *jwt.Token) (any, error) {
		return service.secret, nil
	}, jwt.WithAudience(service.audience), jwt.WithIssuer(service.issuer))
	if err != nil {
		return Claims{}, fmt.Errorf("parse jwt: %w", err)
	}

	if !parsedToken.Valid {
		return Claims{}, fmt.Errorf("invalid jwt")
	}

	return claims, nil
}
