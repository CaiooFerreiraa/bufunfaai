package auth

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	clerk "github.com/clerk/clerk-sdk-go/v2"
	clerkjwt "github.com/clerk/clerk-sdk-go/v2/jwt"
	clerkuser "github.com/clerk/clerk-sdk-go/v2/user"
	"github.com/google/uuid"

	userservice "github.com/bufunfaai/bufunfaai/apps/api/internal/modules/users/application/service"
	userentity "github.com/bufunfaai/bufunfaai/apps/api/internal/modules/users/domain/entity"
	userrepository "github.com/bufunfaai/bufunfaai/apps/api/internal/modules/users/domain/repository"
	"github.com/bufunfaai/bufunfaai/apps/api/internal/platform/config"
)

type AuthenticatedIdentity struct {
	LocalUserID string
	ClerkUserID string
	Email       string
	SessionID   string
}

type TokenAuthenticator interface {
	Authenticate(ctx context.Context, bearerToken string) (AuthenticatedIdentity, error)
}

type ClerkService struct {
	userRepository userrepository.UserRepository
}

func NewClerkService(cfg config.Config, userRepository userrepository.UserRepository) (*ClerkService, error) {
	if strings.TrimSpace(cfg.ClerkSecretKey) == "" {
		return nil, fmt.Errorf("missing CLERK_SECRET_KEY")
	}

	clerk.SetKey(cfg.ClerkSecretKey)

	return &ClerkService{
		userRepository: userRepository,
	}, nil
}

func (service *ClerkService) Authenticate(ctx context.Context, bearerToken string) (AuthenticatedIdentity, error) {
	claims, err := clerkjwt.Verify(ctx, &clerkjwt.VerifyParams{
		Token: bearerToken,
	})
	if err != nil {
		return AuthenticatedIdentity{}, fmt.Errorf("verify clerk token: %w", err)
	}

	clerkUser, err := clerkuser.Get(ctx, claims.Subject)
	if err != nil {
		return AuthenticatedIdentity{}, fmt.Errorf("load clerk user: %w", err)
	}

	email := primaryEmailAddress(clerkUser)
	if email == "" {
		return AuthenticatedIdentity{}, fmt.Errorf("clerk user without primary email")
	}

	localUser, err := service.resolveLocalUser(ctx, clerkUser, email)
	if err != nil {
		return AuthenticatedIdentity{}, err
	}

	if updateErr := service.userRepository.UpdateLastLogin(ctx, localUser.ID); updateErr != nil && !errors.Is(updateErr, userservice.ErrUserNotFound) {
		return AuthenticatedIdentity{}, fmt.Errorf("update local user login: %w", updateErr)
	}

	return AuthenticatedIdentity{
		LocalUserID: localUser.ID,
		ClerkUserID: clerkUser.ID,
		Email:       email,
		SessionID:   claims.SessionID,
	}, nil
}

func (service *ClerkService) resolveLocalUser(
	ctx context.Context,
	clerkUser *clerk.User,
	email string,
) (userentity.User, error) {
	localUser, err := service.userRepository.GetByEmail(ctx, email)
	if err == nil {
		return localUser, nil
	}

	if !errors.Is(err, userservice.ErrUserNotFound) {
		return userentity.User{}, fmt.Errorf("find local user by email: %w", err)
	}

	now := time.Now().UTC()
	fullName := strings.TrimSpace(strings.Join([]string{
		stringValue(clerkUser.FirstName),
		stringValue(clerkUser.LastName),
	}, " "))
	if fullName == "" {
		fullName = strings.Split(email, "@")[0]
	}

	var emailVerifiedAt *time.Time
	if primaryEmailVerified(clerkUser) {
		emailVerifiedAt = &now
	}

	createdUser, createErr := service.userRepository.Create(ctx, userentity.User{
		ID:              uuid.NewString(),
		FullName:        fullName,
		Email:           email,
		PasswordHash:    "clerk_managed",
		Status:          "active",
		EmailVerifiedAt: emailVerifiedAt,
		LastLoginAt:     &now,
		CreatedAt:       now,
		UpdatedAt:       now,
	})
	if createErr != nil {
		return userentity.User{}, fmt.Errorf("create local user from clerk identity: %w", createErr)
	}

	return createdUser, nil
}

func primaryEmailAddress(user *clerk.User) string {
	if user == nil {
		return ""
	}

	if user.PrimaryEmailAddressID != nil {
		for _, emailAddress := range user.EmailAddresses {
			if emailAddress != nil && emailAddress.ID == *user.PrimaryEmailAddressID {
				return strings.ToLower(strings.TrimSpace(emailAddress.EmailAddress))
			}
		}
	}

	for _, emailAddress := range user.EmailAddresses {
		if emailAddress != nil && strings.TrimSpace(emailAddress.EmailAddress) != "" {
			return strings.ToLower(strings.TrimSpace(emailAddress.EmailAddress))
		}
	}

	return ""
}

func primaryEmailVerified(user *clerk.User) bool {
	if user == nil || user.PrimaryEmailAddressID == nil {
		return false
	}

	for _, emailAddress := range user.EmailAddresses {
		if emailAddress == nil || emailAddress.ID != *user.PrimaryEmailAddressID || emailAddress.Verification == nil {
			continue
		}

		return strings.EqualFold(emailAddress.Verification.Status, "verified")
	}

	return false
}

func stringValue(value *string) string {
	if value == nil {
		return ""
	}

	return strings.TrimSpace(*value)
}
