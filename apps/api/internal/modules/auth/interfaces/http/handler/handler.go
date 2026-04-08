package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	authdto "github.com/bufunfaai/bufunfaai/apps/api/internal/modules/auth/application/dto"
	authusecase "github.com/bufunfaai/bufunfaai/apps/api/internal/modules/auth/application/usecase"
	authpresenter "github.com/bufunfaai/bufunfaai/apps/api/internal/modules/auth/interfaces/http/presenter"
	platformvalidator "github.com/bufunfaai/bufunfaai/apps/api/internal/platform/validator"
	sharederrors "github.com/bufunfaai/bufunfaai/apps/api/internal/shared/errors"
	"github.com/bufunfaai/bufunfaai/apps/api/internal/shared/middleware"
	"github.com/bufunfaai/bufunfaai/apps/api/internal/shared/response"
)

type Handler struct {
	registerUseCase  *authusecase.RegisterUseCase
	loginUseCase     *authusecase.LoginUseCase
	refreshUseCase   *authusecase.RefreshUseCase
	logoutUseCase    *authusecase.LogoutUseCase
	logoutAllUseCase *authusecase.LogoutAllUseCase
	validator        *platformvalidator.Validator
}

func NewHandler(
	registerUseCase *authusecase.RegisterUseCase,
	loginUseCase *authusecase.LoginUseCase,
	refreshUseCase *authusecase.RefreshUseCase,
	logoutUseCase *authusecase.LogoutUseCase,
	logoutAllUseCase *authusecase.LogoutAllUseCase,
	validator *platformvalidator.Validator,
) *Handler {
	return &Handler{
		registerUseCase:  registerUseCase,
		loginUseCase:     loginUseCase,
		refreshUseCase:   refreshUseCase,
		logoutUseCase:    logoutUseCase,
		logoutAllUseCase: logoutAllUseCase,
		validator:        validator,
	}
}

func (handler *Handler) Register(context *gin.Context) {
	var request authdto.RegisterRequest
	if err := context.ShouldBindJSON(&request); err != nil {
		response.Error(context, sharederrors.New("INVALID_PAYLOAD", "payload invalido", http.StatusBadRequest))
		return
	}

	if appError := handler.validator.Validate(request); appError != nil {
		response.Error(context, appError)
		return
	}

	result, appError := handler.registerUseCase.Execute(
		context.Request.Context(),
		request,
		context.ClientIP(),
		context.Request.UserAgent(),
	)
	if appError != nil {
		response.Error(context, appError)
		return
	}

	response.Created(context, authpresenter.AuthResult(result))
}

func (handler *Handler) Login(context *gin.Context) {
	var request authdto.LoginRequest
	if err := context.ShouldBindJSON(&request); err != nil {
		response.Error(context, sharederrors.New("INVALID_PAYLOAD", "payload invalido", http.StatusBadRequest))
		return
	}

	if appError := handler.validator.Validate(request); appError != nil {
		response.Error(context, appError)
		return
	}

	result, appError := handler.loginUseCase.Execute(
		context.Request.Context(),
		request,
		context.ClientIP(),
		context.Request.UserAgent(),
	)
	if appError != nil {
		response.Error(context, appError)
		return
	}

	response.OK(context, authpresenter.AuthResult(result))
}

func (handler *Handler) Refresh(context *gin.Context) {
	var request authdto.RefreshRequest
	if err := context.ShouldBindJSON(&request); err != nil {
		response.Error(context, sharederrors.New("INVALID_PAYLOAD", "payload invalido", http.StatusBadRequest))
		return
	}

	if appError := handler.validator.Validate(request); appError != nil {
		response.Error(context, appError)
		return
	}

	result, appError := handler.refreshUseCase.Execute(
		context.Request.Context(),
		request,
		context.ClientIP(),
		context.Request.UserAgent(),
	)
	if appError != nil {
		response.Error(context, appError)
		return
	}

	response.OK(context, authpresenter.AuthResult(result))
}

func (handler *Handler) Logout(context *gin.Context) {
	var request authdto.LogoutRequest
	if err := context.ShouldBindJSON(&request); err != nil {
		response.Error(context, sharederrors.New("INVALID_PAYLOAD", "payload invalido", http.StatusBadRequest))
		return
	}

	if appError := handler.validator.Validate(request); appError != nil {
		response.Error(context, appError)
		return
	}

	appError := handler.logoutUseCase.Execute(context.Request.Context(), request.RefreshToken)
	if appError != nil {
		response.Error(context, appError)
		return
	}

	response.OK(context, gin.H{
		"logged_out": true,
	})
}

func (handler *Handler) LogoutAll(context *gin.Context) {
	userID := middleware.CurrentUserID(context)
	if userID == "" {
		response.Error(context, sharederrors.New("UNAUTHORIZED", "autenticacao obrigatoria", http.StatusUnauthorized))
		return
	}

	appError := handler.logoutAllUseCase.Execute(context.Request.Context(), userID)
	if appError != nil {
		response.Error(context, appError)
		return
	}

	response.OK(context, gin.H{
		"logged_out_all": true,
	})
}
