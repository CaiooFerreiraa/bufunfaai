package repository

import (
	"context"
	"time"

	"github.com/bufunfaai/bufunfaai/apps/api/internal/modules/analytics/domain/entity"
)

type AnalyticsRepository interface {
	ListBudgetSnapshotsByMonth(ctx context.Context, userID string, month time.Time) ([]entity.BudgetSnapshot, error)
	GetLatestScore(ctx context.Context, userID string) (*entity.FinancialScore, error)
	GetForecastByMonth(ctx context.Context, userID string, referenceMonth time.Time) (*entity.CashflowForecast, error)
	ListInsights(ctx context.Context, userID string, limit int, now time.Time) ([]entity.Insight, error)
	ListAnomalies(ctx context.Context, userID string, limit int) ([]entity.Anomaly, error)
	ListGoals(ctx context.Context, userID string) ([]entity.FinancialGoal, error)
	GetGoalByID(ctx context.Context, userID string, goalID string) (entity.FinancialGoal, error)
	CreateGoal(ctx context.Context, goal entity.FinancialGoal) (entity.FinancialGoal, error)
	UpdateGoal(ctx context.Context, goal entity.FinancialGoal) (entity.FinancialGoal, error)
}
