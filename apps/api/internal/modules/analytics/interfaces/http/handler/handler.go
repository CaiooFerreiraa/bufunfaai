package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	analyticsdto "github.com/bufunfaai/bufunfaai/apps/api/internal/modules/analytics/application/dto"
	analyticsusecase "github.com/bufunfaai/bufunfaai/apps/api/internal/modules/analytics/application/usecase"
	analyticspresenter "github.com/bufunfaai/bufunfaai/apps/api/internal/modules/analytics/interfaces/http/presenter"
	platformvalidator "github.com/bufunfaai/bufunfaai/apps/api/internal/platform/validator"
	sharederrors "github.com/bufunfaai/bufunfaai/apps/api/internal/shared/errors"
	"github.com/bufunfaai/bufunfaai/apps/api/internal/shared/middleware"
	"github.com/bufunfaai/bufunfaai/apps/api/internal/shared/response"
)

type Handler struct {
	useCases  *analyticsusecase.UseCases
	validator *platformvalidator.Validator
}

func NewHandler(useCases *analyticsusecase.UseCases, validator *platformvalidator.Validator) *Handler {
	return &Handler{useCases: useCases, validator: validator}
}

func (handler *Handler) Overview(context *gin.Context) {
	userID := middleware.CurrentUserID(context)
	overview, appError := handler.useCases.GetOverview(context.Request.Context(), userID, context.Query("month"))
	if appError != nil {
		response.Error(context, appError)
		return
	}

	response.OK(context, gin.H{"overview": analyticspresenter.OverviewOutput(overview)})
}

func (handler *Handler) Budgets(context *gin.Context) {
	userID := middleware.CurrentUserID(context)
	month, appError := parseMonth(context.Query("month"))
	if appError != nil {
		response.Error(context, appError)
		return
	}

	items, appError := handler.useCases.ListBudgetSnapshots(context.Request.Context(), userID, month)
	if appError != nil {
		response.Error(context, appError)
		return
	}

	payload := make([]analyticsdto.BudgetSnapshotOutput, 0, len(items))
	for _, item := range items {
		payload = append(payload, analyticspresenter.BudgetSnapshotOutput(item))
	}

	response.OK(context, gin.H{"budgets": payload})
}

func (handler *Handler) LatestScore(context *gin.Context) {
	userID := middleware.CurrentUserID(context)
	item, appError := handler.useCases.GetLatestScore(context.Request.Context(), userID)
	if appError != nil {
		response.Error(context, appError)
		return
	}

	if item == nil {
		response.OK(context, gin.H{"score": nil})
		return
	}

	response.OK(context, gin.H{"score": analyticspresenter.ScoreOutput(*item)})
}

func (handler *Handler) Forecast(context *gin.Context) {
	userID := middleware.CurrentUserID(context)
	month, appError := parseMonth(context.Query("month"))
	if appError != nil {
		response.Error(context, appError)
		return
	}

	item, appError := handler.useCases.GetForecastByMonth(context.Request.Context(), userID, month)
	if appError != nil {
		response.Error(context, appError)
		return
	}

	if item == nil {
		response.OK(context, gin.H{"forecast": nil})
		return
	}

	response.OK(context, gin.H{"forecast": analyticspresenter.ForecastOutput(*item)})
}

func (handler *Handler) Insights(context *gin.Context) {
	userID := middleware.CurrentUserID(context)
	items, appError := handler.useCases.ListInsights(context.Request.Context(), userID, 10)
	if appError != nil {
		response.Error(context, appError)
		return
	}

	payload := make([]analyticsdto.InsightOutput, 0, len(items))
	for _, item := range items {
		payload = append(payload, analyticspresenter.InsightOutput(item))
	}

	response.OK(context, gin.H{"insights": payload})
}

func (handler *Handler) Anomalies(context *gin.Context) {
	userID := middleware.CurrentUserID(context)
	items, appError := handler.useCases.ListAnomalies(context.Request.Context(), userID, 10)
	if appError != nil {
		response.Error(context, appError)
		return
	}

	payload := make([]analyticsdto.AnomalyOutput, 0, len(items))
	for _, item := range items {
		payload = append(payload, analyticspresenter.AnomalyOutput(item))
	}

	response.OK(context, gin.H{"anomalies": payload})
}

func (handler *Handler) Goals(context *gin.Context) {
	userID := middleware.CurrentUserID(context)
	items, appError := handler.useCases.ListGoals(context.Request.Context(), userID)
	if appError != nil {
		response.Error(context, appError)
		return
	}

	payload := make([]analyticsdto.GoalOutput, 0, len(items))
	for _, item := range items {
		payload = append(payload, analyticspresenter.GoalOutput(item))
	}

	response.OK(context, gin.H{"goals": payload})
}

func (handler *Handler) CreateGoal(context *gin.Context) {
	userID := middleware.CurrentUserID(context)

	var request analyticsdto.CreateGoalRequest
	if err := context.ShouldBindJSON(&request); err != nil {
		response.Error(context, sharederrors.New("INVALID_PAYLOAD", "payload invalido", http.StatusBadRequest))
		return
	}

	if appError := handler.validator.Validate(request); appError != nil {
		response.Error(context, appError)
		return
	}

	item, appError := handler.useCases.CreateGoal(context.Request.Context(), userID, request)
	if appError != nil {
		response.Error(context, appError)
		return
	}

	response.Created(context, gin.H{"goal": analyticspresenter.GoalOutput(item)})
}

func (handler *Handler) UpdateGoal(context *gin.Context) {
	userID := middleware.CurrentUserID(context)

	var request analyticsdto.UpdateGoalRequest
	if err := context.ShouldBindJSON(&request); err != nil {
		response.Error(context, sharederrors.New("INVALID_PAYLOAD", "payload invalido", http.StatusBadRequest))
		return
	}

	if appError := handler.validator.Validate(request); appError != nil {
		response.Error(context, appError)
		return
	}

	item, appError := handler.useCases.UpdateGoal(context.Request.Context(), userID, context.Param("id"), request)
	if appError != nil {
		response.Error(context, appError)
		return
	}

	response.OK(context, gin.H{"goal": analyticspresenter.GoalOutput(item)})
}

func parseMonth(rawMonth string) (time.Time, *sharederrors.AppError) {
	if rawMonth == "" {
		now := time.Now().UTC()
		return time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC), nil
	}

	month, err := time.Parse("2006-01", rawMonth)
	if err != nil {
		return time.Time{}, sharederrors.New("INVALID_REFERENCE_MONTH", "mes de referencia invalido", 400)
	}

	return time.Date(month.Year(), month.Month(), 1, 0, 0, 0, 0, time.UTC), nil
}
