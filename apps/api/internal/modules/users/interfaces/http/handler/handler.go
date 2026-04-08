package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/bufunfaai/bufunfaai/apps/api/internal/modules/users/application/dto"
	getcurrentuser "github.com/bufunfaai/bufunfaai/apps/api/internal/modules/users/application/usecase"
	userpresenter "github.com/bufunfaai/bufunfaai/apps/api/internal/modules/users/interfaces/http/presenter"
	platformvalidator "github.com/bufunfaai/bufunfaai/apps/api/internal/platform/validator"
	sharederrors "github.com/bufunfaai/bufunfaai/apps/api/internal/shared/errors"
	"github.com/bufunfaai/bufunfaai/apps/api/internal/shared/middleware"
	"github.com/bufunfaai/bufunfaai/apps/api/internal/shared/response"
)

type Handler struct {
	getCurrentUserUseCase    *getcurrentuser.GetCurrentUserUseCase
	updateCurrentUserUseCase *getcurrentuser.UpdateCurrentUserUseCase
	validator                *platformvalidator.Validator
}

func NewHandler(
	getCurrentUserUseCase *getcurrentuser.GetCurrentUserUseCase,
	updateCurrentUserUseCase *getcurrentuser.UpdateCurrentUserUseCase,
	validator *platformvalidator.Validator,
) *Handler {
	return &Handler{
		getCurrentUserUseCase:    getCurrentUserUseCase,
		updateCurrentUserUseCase: updateCurrentUserUseCase,
		validator:                validator,
	}
}

func (handler *Handler) Me(context *gin.Context) {
	userID := middleware.CurrentUserID(context)
	if userID == "" {
		response.Error(context, sharederrors.New("UNAUTHORIZED", "autenticacao obrigatoria", http.StatusUnauthorized))
		return
	}

	user, appError := handler.getCurrentUserUseCase.Execute(context.Request.Context(), userID)
	if appError != nil {
		response.Error(context, appError)
		return
	}

	response.OK(context, gin.H{
		"user": userpresenter.UserOutput(user),
	})
}

func (handler *Handler) UpdateMe(context *gin.Context) {
	userID := middleware.CurrentUserID(context)
	if userID == "" {
		response.Error(context, sharederrors.New("UNAUTHORIZED", "autenticacao obrigatoria", http.StatusUnauthorized))
		return
	}

	var request dto.UpdateCurrentUserRequest
	if err := context.ShouldBindJSON(&request); err != nil {
		response.Error(context, sharederrors.New("INVALID_PAYLOAD", "payload invalido", http.StatusBadRequest))
		return
	}

	if appError := handler.validator.Validate(request); appError != nil {
		response.Error(context, appError)
		return
	}

	user, appError := handler.updateCurrentUserUseCase.Execute(
		context.Request.Context(),
		userID,
		request.FullName,
		request.Phone,
	)
	if appError != nil {
		response.Error(context, appError)
		return
	}

	response.OK(context, gin.H{
		"user": userpresenter.UserOutput(user),
	})
}
