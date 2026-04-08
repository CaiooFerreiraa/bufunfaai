package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/bufunfaai/bufunfaai/apps/api/internal/shared/response"
)

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

func (handler *Handler) List(context *gin.Context) {
	response.OK(context, gin.H{
		"consents": []any{},
	})
}

func (handler *Handler) Create(context *gin.Context) {
	response.Success(context, http.StatusNotImplemented, gin.H{
		"message": "modulo de consents sera implementado na proxima fase",
	})
}

func (handler *Handler) Delete(context *gin.Context) {
	response.Success(context, http.StatusNotImplemented, gin.H{
		"message": "modulo de consents sera implementado na proxima fase",
	})
}
