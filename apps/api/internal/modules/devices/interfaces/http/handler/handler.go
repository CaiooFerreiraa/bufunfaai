package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	deleteusecase "github.com/bufunfaai/bufunfaai/apps/api/internal/modules/devices/application/usecase"
	devicepresenter "github.com/bufunfaai/bufunfaai/apps/api/internal/modules/devices/interfaces/http/presenter"
	sharederrors "github.com/bufunfaai/bufunfaai/apps/api/internal/shared/errors"
	"github.com/bufunfaai/bufunfaai/apps/api/internal/shared/middleware"
	"github.com/bufunfaai/bufunfaai/apps/api/internal/shared/response"
)

type Handler struct {
	listDevicesUseCase  *deleteusecase.ListDevicesUseCase
	deleteDeviceUseCase *deleteusecase.DeleteDeviceUseCase
}

func NewHandler(
	listDevicesUseCase *deleteusecase.ListDevicesUseCase,
	deleteDeviceUseCase *deleteusecase.DeleteDeviceUseCase,
) *Handler {
	return &Handler{
		listDevicesUseCase:  listDevicesUseCase,
		deleteDeviceUseCase: deleteDeviceUseCase,
	}
}

func (handler *Handler) List(context *gin.Context) {
	userID := middleware.CurrentUserID(context)
	if userID == "" {
		response.Error(context, sharederrors.New("UNAUTHORIZED", "autenticacao obrigatoria", http.StatusUnauthorized))
		return
	}

	devices, appError := handler.listDevicesUseCase.Execute(context.Request.Context(), userID)
	if appError != nil {
		response.Error(context, appError)
		return
	}

	payload := make([]any, 0, len(devices))
	for _, device := range devices {
		payload = append(payload, devicepresenter.DeviceOutput(device))
	}

	response.OK(context, gin.H{
		"devices": payload,
	})
}

func (handler *Handler) Delete(context *gin.Context) {
	userID := middleware.CurrentUserID(context)
	if userID == "" {
		response.Error(context, sharederrors.New("UNAUTHORIZED", "autenticacao obrigatoria", http.StatusUnauthorized))
		return
	}

	deviceID := context.Param("id")
	if deviceID == "" {
		response.Error(context, sharederrors.New("VALIDATION_ERROR", "id do dispositivo obrigatorio", http.StatusBadRequest))
		return
	}

	appError := handler.deleteDeviceUseCase.Execute(context.Request.Context(), userID, deviceID)
	if appError != nil {
		response.Error(context, appError)
		return
	}

	response.OK(context, gin.H{
		"deleted": true,
	})
}
