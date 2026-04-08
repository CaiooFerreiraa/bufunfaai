package usecase

import (
	"context"
	"time"

	analyticsdto "github.com/bufunfaai/bufunfaai/apps/api/internal/modules/analytics/application/dto"
	analyticsservice "github.com/bufunfaai/bufunfaai/apps/api/internal/modules/analytics/application/service"
	"github.com/bufunfaai/bufunfaai/apps/api/internal/modules/analytics/domain/entity"
	sharederrors "github.com/bufunfaai/bufunfaai/apps/api/internal/shared/errors"
)

type UseCases struct {
	service *analyticsservice.Service
}

func New(service *analyticsservice.Service) *UseCases {
	return &UseCases{service: service}
}

func (useCases *UseCases) GetOverview(ctx context.Context, userID string, referenceMonth string) (entity.Overview, *sharederrors.AppError) {
	return useCases.service.GetOverview(ctx, userID, referenceMonth)
}

func (useCases *UseCases) ListBudgetSnapshots(ctx context.Context, userID string, month time.Time) ([]entity.BudgetSnapshot, *sharederrors.AppError) {
	return useCases.service.ListBudgetSnapshots(ctx, userID, month)
}

func (useCases *UseCases) GetLatestScore(ctx context.Context, userID string) (*entity.FinancialScore, *sharederrors.AppError) {
	return useCases.service.GetLatestScore(ctx, userID)
}

func (useCases *UseCases) GetForecastByMonth(ctx context.Context, userID string, month time.Time) (*entity.CashflowForecast, *sharederrors.AppError) {
	return useCases.service.GetForecastByMonth(ctx, userID, month)
}

func (useCases *UseCases) ListInsights(ctx context.Context, userID string, limit int) ([]entity.Insight, *sharederrors.AppError) {
	return useCases.service.ListInsights(ctx, userID, limit)
}

func (useCases *UseCases) ListAnomalies(ctx context.Context, userID string, limit int) ([]entity.Anomaly, *sharederrors.AppError) {
	return useCases.service.ListAnomalies(ctx, userID, limit)
}

func (useCases *UseCases) ListGoals(ctx context.Context, userID string) ([]entity.FinancialGoal, *sharederrors.AppError) {
	return useCases.service.ListGoals(ctx, userID)
}

func (useCases *UseCases) CreateGoal(ctx context.Context, userID string, request analyticsdto.CreateGoalRequest) (entity.FinancialGoal, *sharederrors.AppError) {
	return useCases.service.CreateGoal(ctx, userID, request)
}

func (useCases *UseCases) UpdateGoal(ctx context.Context, userID string, goalID string, request analyticsdto.UpdateGoalRequest) (entity.FinancialGoal, *sharederrors.AppError) {
	return useCases.service.UpdateGoal(ctx, userID, goalID, request)
}
