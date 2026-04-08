package presenter

import (
	"github.com/bufunfaai/bufunfaai/apps/api/internal/modules/users/application/dto"
	"github.com/bufunfaai/bufunfaai/apps/api/internal/modules/users/domain/entity"
)

func UserOutput(user entity.User) dto.UserOutput {
	return dto.UserOutput{
		ID:        user.ID,
		FullName:  user.FullName,
		Email:     user.Email,
		Phone:     user.Phone,
		Status:    user.Status,
		CreatedAt: user.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: user.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}
