package service

import (
	"context"
	"errors"
	"strings"

	"github.com/bufunfaai/bufunfaai/apps/api/internal/modules/users/domain/entity"
	userrepository "github.com/bufunfaai/bufunfaai/apps/api/internal/modules/users/domain/repository"
	sharederrors "github.com/bufunfaai/bufunfaai/apps/api/internal/shared/errors"
)

type UserService struct {
	userRepository userrepository.UserRepository
}

func NewUserService(userRepository userrepository.UserRepository) *UserService {
	return &UserService{
		userRepository: userRepository,
	}
}

func (service *UserService) GetCurrentUser(ctx context.Context, userID string) (entity.User, *sharederrors.AppError) {
	user, err := service.userRepository.GetByID(ctx, userID)
	if err != nil {
		return entity.User{}, mapUserRepositoryError(err)
	}

	return user, nil
}

func (service *UserService) UpdateCurrentUser(ctx context.Context, userID string, fullName string, phone string) (entity.User, *sharederrors.AppError) {
	user, err := service.userRepository.UpdateProfile(
		ctx,
		userID,
		strings.TrimSpace(fullName),
		strings.TrimSpace(phone),
	)
	if err != nil {
		return entity.User{}, mapUserRepositoryError(err)
	}

	return user, nil
}

func mapUserRepositoryError(err error) *sharederrors.AppError {
	if errors.Is(err, ErrUserNotFound) {
		return sharederrors.New("USER_NOT_FOUND", "usuario nao encontrado", 404)
	}

	return sharederrors.Wrap("INTERNAL_ERROR", "erro interno ao processar usuario", 500, err)
}

var ErrUserNotFound = errors.New("user not found")
