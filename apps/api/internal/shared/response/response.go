package response

import (
	"net/http"

	"github.com/gin-gonic/gin"

	sharederrors "github.com/bufunfaai/bufunfaai/apps/api/internal/shared/errors"
)

const requestIDKey string = "request_id"

type successMeta struct {
	RequestID string `json:"request_id"`
}

type errorBody struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func Success(context *gin.Context, status int, data any) {
	context.JSON(status, gin.H{
		"success": true,
		"data":    data,
		"meta": successMeta{
			RequestID: RequestID(context),
		},
	})
}

func OK(context *gin.Context, data any) {
	Success(context, http.StatusOK, data)
}

func Created(context *gin.Context, data any) {
	Success(context, http.StatusCreated, data)
}

func Error(context *gin.Context, appError *sharederrors.AppError) {
	context.JSON(appError.Status, gin.H{
		"success": false,
		"error": errorBody{
			Code:    appError.Code,
			Message: appError.Message,
		},
		"meta": successMeta{
			RequestID: RequestID(context),
		},
	})
}

func RequestID(context *gin.Context) string {
	value, ok := context.Get(requestIDKey)
	if !ok {
		return ""
	}

	requestID, ok := value.(string)
	if !ok {
		return ""
	}

	return requestID
}
