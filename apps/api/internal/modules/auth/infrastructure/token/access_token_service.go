package token

import (
	"time"

	userentity "github.com/bufunfaai/bufunfaai/apps/api/internal/modules/users/domain/entity"
	platformauth "github.com/bufunfaai/bufunfaai/apps/api/internal/platform/auth"
)

type AccessTokenService struct {
	jwtService *platformauth.JWTService
}

func NewAccessTokenService(jwtService *platformauth.JWTService) *AccessTokenService {
	return &AccessTokenService{jwtService: jwtService}
}

func (service *AccessTokenService) IssueAccessToken(user userentity.User, sessionID string, now time.Time) (string, time.Time, error) {
	return service.jwtService.Issue(platformauth.IssueTokenInput{
		UserID:    user.ID,
		Email:     user.Email,
		SessionID: sessionID,
		Role:      "user",
		Now:       now,
	})
}
