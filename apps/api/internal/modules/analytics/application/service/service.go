package service

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"

	analyticsdto "github.com/bufunfaai/bufunfaai/apps/api/internal/modules/analytics/application/dto"
	"github.com/bufunfaai/bufunfaai/apps/api/internal/modules/analytics/domain/entity"
	analyticsrepository "github.com/bufunfaai/bufunfaai/apps/api/internal/modules/analytics/domain/repository"
	sharederrors "github.com/bufunfaai/bufunfaai/apps/api/internal/shared/errors"
)

var ErrGoalNotFound = errors.New("goal not found")

type Service struct {
	repository analyticsrepository.AnalyticsRepository
	now        func() time.Time
}

func NewService(repository analyticsrepository.AnalyticsRepository) *Service {
	return &Service{
		repository: repository,
		now:        time.Now,
	}
}

func (service *Service) GetOverview(ctx context.Context, userID string, referenceMonth string) (entity.Overview, *sharederrors.AppError) {
	month, appError := resolveReferenceMonth(referenceMonth, service.now().UTC())
	if appError != nil {
		return entity.Overview{}, appError
	}

	budgets, appError := service.ListBudgetSnapshots(ctx, userID, month)
	if appError != nil {
		return entity.Overview{}, appError
	}

	score, appError := service.GetLatestScore(ctx, userID)
	if appError != nil {
		return entity.Overview{}, appError
	}

	forecast, appError := service.GetForecastByMonth(ctx, userID, month)
	if appError != nil {
		return entity.Overview{}, appError
	}

	insights, appError := service.ListInsights(ctx, userID, 3)
	if appError != nil {
		return entity.Overview{}, appError
	}

	anomalies, appError := service.ListAnomalies(ctx, userID, 5)
	if appError != nil {
		return entity.Overview{}, appError
	}

	goals, appError := service.ListGoals(ctx, userID)
	if appError != nil {
		return entity.Overview{}, appError
	}

	return entity.Overview{
		ReferenceMonth: month.Format("2006-01"),
		Budgets:        budgets,
		Score:          score,
		Forecast:       forecast,
		Insights:       insights,
		Anomalies:      anomalies,
		Goals:          goals,
	}, nil
}

func (service *Service) ListBudgetSnapshots(ctx context.Context, userID string, month time.Time) ([]entity.BudgetSnapshot, *sharederrors.AppError) {
	items, err := service.repository.ListBudgetSnapshotsByMonth(ctx, userID, month)
	if err != nil {
		return nil, sharederrors.Wrap("ANALYTICS_BUDGETS_ERROR", "erro ao listar snapshots de orcamento", 500, err)
	}

	return items, nil
}

func (service *Service) GetLatestScore(ctx context.Context, userID string) (*entity.FinancialScore, *sharederrors.AppError) {
	item, err := service.repository.GetLatestScore(ctx, userID)
	if err != nil {
		return nil, sharederrors.Wrap("ANALYTICS_SCORE_ERROR", "erro ao buscar score financeiro", 500, err)
	}

	return item, nil
}

func (service *Service) GetForecastByMonth(ctx context.Context, userID string, month time.Time) (*entity.CashflowForecast, *sharederrors.AppError) {
	item, err := service.repository.GetForecastByMonth(ctx, userID, month)
	if err != nil {
		return nil, sharederrors.Wrap("ANALYTICS_FORECAST_ERROR", "erro ao buscar previsao de caixa", 500, err)
	}

	return item, nil
}

func (service *Service) ListInsights(ctx context.Context, userID string, limit int) ([]entity.Insight, *sharederrors.AppError) {
	items, err := service.repository.ListInsights(ctx, userID, limit, service.now().UTC())
	if err != nil {
		return nil, sharederrors.Wrap("ANALYTICS_INSIGHTS_ERROR", "erro ao listar insights", 500, err)
	}

	return items, nil
}

func (service *Service) ListAnomalies(ctx context.Context, userID string, limit int) ([]entity.Anomaly, *sharederrors.AppError) {
	items, err := service.repository.ListAnomalies(ctx, userID, limit)
	if err != nil {
		return nil, sharederrors.Wrap("ANALYTICS_ANOMALIES_ERROR", "erro ao listar anomalias", 500, err)
	}

	return items, nil
}

func (service *Service) ListGoals(ctx context.Context, userID string) ([]entity.FinancialGoal, *sharederrors.AppError) {
	items, err := service.repository.ListGoals(ctx, userID)
	if err != nil {
		return nil, sharederrors.Wrap("ANALYTICS_GOALS_ERROR", "erro ao listar metas financeiras", 500, err)
	}

	return items, nil
}

func (service *Service) CreateGoal(ctx context.Context, userID string, request analyticsdto.CreateGoalRequest) (entity.FinancialGoal, *sharederrors.AppError) {
	now := service.now().UTC()
	dueDate, appError := parseDueDate(request.DueDate)
	if appError != nil {
		return entity.FinancialGoal{}, appError
	}

	status := request.Status
	if status == "" {
		status = entity.GoalStatusInProgress
	}

	goal := entity.FinancialGoal{
		ID:                 uuid.NewString(),
		UserID:             userID,
		Title:              request.Title,
		GoalType:           request.GoalType,
		TargetAmountCents:  request.TargetAmountCents,
		CurrentAmountCents: request.CurrentAmountCents,
		DueDate:            dueDate,
		Status:             status,
		CreatedAt:          now,
		UpdatedAt:          now,
	}

	createdGoal, err := service.repository.CreateGoal(ctx, goal)
	if err != nil {
		return entity.FinancialGoal{}, sharederrors.Wrap("ANALYTICS_GOALS_ERROR", "erro ao criar meta financeira", 500, err)
	}

	return createdGoal, nil
}

func (service *Service) UpdateGoal(ctx context.Context, userID string, goalID string, request analyticsdto.UpdateGoalRequest) (entity.FinancialGoal, *sharederrors.AppError) {
	goal, err := service.repository.GetGoalByID(ctx, userID, goalID)
	if err != nil {
		if errors.Is(err, ErrGoalNotFound) {
			return entity.FinancialGoal{}, sharederrors.New("GOAL_NOT_FOUND", "meta financeira nao encontrada", 404)
		}

		return entity.FinancialGoal{}, sharederrors.Wrap("ANALYTICS_GOALS_ERROR", "erro ao buscar meta financeira", 500, err)
	}

	if request.Title != nil {
		goal.Title = *request.Title
	}

	if request.GoalType != nil {
		goal.GoalType = *request.GoalType
	}

	if request.TargetAmountCents != nil {
		goal.TargetAmountCents = *request.TargetAmountCents
	}

	if request.CurrentAmountCents != nil {
		goal.CurrentAmountCents = *request.CurrentAmountCents
	}

	if request.DueDate != nil {
		dueDate, appError := parseDueDate(*request.DueDate)
		if appError != nil {
			return entity.FinancialGoal{}, appError
		}

		goal.DueDate = dueDate
	}

	if request.Status != nil {
		goal.Status = *request.Status
	}

	goal.UpdatedAt = service.now().UTC()
	updatedGoal, err := service.repository.UpdateGoal(ctx, goal)
	if err != nil {
		if errors.Is(err, ErrGoalNotFound) {
			return entity.FinancialGoal{}, sharederrors.New("GOAL_NOT_FOUND", "meta financeira nao encontrada", 404)
		}

		return entity.FinancialGoal{}, sharederrors.Wrap("ANALYTICS_GOALS_ERROR", "erro ao atualizar meta financeira", 500, err)
	}

	return updatedGoal, nil
}

func resolveReferenceMonth(rawMonth string, now time.Time) (time.Time, *sharederrors.AppError) {
	if rawMonth == "" {
		return time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC), nil
	}

	month, err := time.Parse("2006-01", rawMonth)
	if err != nil {
		return time.Time{}, sharederrors.New("INVALID_REFERENCE_MONTH", "mes de referencia invalido", 400)
	}

	return time.Date(month.Year(), month.Month(), 1, 0, 0, 0, 0, time.UTC), nil
}

func parseDueDate(rawDate string) (*time.Time, *sharederrors.AppError) {
	if rawDate == "" {
		return nil, nil
	}

	dueDate, err := time.Parse("2006-01-02", rawDate)
	if err != nil {
		return nil, sharederrors.New("INVALID_DUE_DATE", "data limite invalida", 400)
	}

	value := dueDate.UTC()
	return &value, nil
}
