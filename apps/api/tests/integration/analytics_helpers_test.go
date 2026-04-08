package integration

import (
	"context"
	"sync"
	"time"

	analyticsservice "github.com/bufunfaai/bufunfaai/apps/api/internal/modules/analytics/application/service"
	analyticsentity "github.com/bufunfaai/bufunfaai/apps/api/internal/modules/analytics/domain/entity"
	analyticsrepository "github.com/bufunfaai/bufunfaai/apps/api/internal/modules/analytics/domain/repository"
)

type inMemoryAnalyticsRepository struct {
	mutex     sync.RWMutex
	budgets   []analyticsentity.BudgetSnapshot
	insights  []analyticsentity.Insight
	anomalies []analyticsentity.Anomaly
	goals     map[string]analyticsentity.FinancialGoal
	score     *analyticsentity.FinancialScore
	forecast  *analyticsentity.CashflowForecast
}

func newInMemoryAnalyticsRepository() *inMemoryAnalyticsRepository {
	now := time.Now().UTC()
	month := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)

	return &inMemoryAnalyticsRepository{
		budgets: []analyticsentity.BudgetSnapshot{
			{
				ID:                "budget-1",
				UserID:            "",
				Month:             month,
				Category:          "alimentacao",
				SpentAmountCents:  92000,
				BudgetAmountCents: 110000,
				UsagePercent:      83.64,
				Status:            "warning",
				UpdatedAt:         now,
			},
		},
		insights: []analyticsentity.Insight{
			{
				ID:          "insight-1",
				UserID:      "",
				InsightType: "attention",
				Title:       "Alimentacao perto do limite",
				Message:     "Voce ja consumiu mais de 80% do orcamento de alimentacao.",
				Priority:    90,
				Source:      "rule_engine",
				CreatedAt:   now,
			},
		},
		anomalies: []analyticsentity.Anomaly{
			{
				ID:            "anomaly-1",
				UserID:        "",
				TransactionID: "transaction-1",
				AnomalyType:   "value_outlier",
				Severity:      "medium",
				Reason:        "Compra acima do padrao historico para a categoria.",
				DetectedAt:    now,
			},
		},
		goals: make(map[string]analyticsentity.FinancialGoal),
		score: &analyticsentity.FinancialScore{
			ID:          "score-1",
			UserID:      "",
			PeriodStart: month,
			PeriodEnd:   month.AddDate(0, 1, -1),
			TotalScore:  66,
			ScoreLabel:  "regular",
			FactorsJSON: `{"positive":["renda estavel"],"negative":["alimentacao acima do planejado"]}`,
			CreatedAt:   now,
		},
		forecast: &analyticsentity.CashflowForecast{
			ID:                    "forecast-1",
			UserID:                "",
			ReferenceMonth:        month,
			ProjectedIncomeCents:  420000,
			ProjectedExpenseCents: 389000,
			ProjectedBalanceCents: 31000,
			ConfidenceLevel:       "expected",
			AssumptionsJSON:       `{"income":"stable","expense":"recent_average"}`,
			CreatedAt:             now,
		},
	}
}

func (repository *inMemoryAnalyticsRepository) ListBudgetSnapshotsByMonth(_ context.Context, userID string, month time.Time) ([]analyticsentity.BudgetSnapshot, error) {
	repository.mutex.RLock()
	defer repository.mutex.RUnlock()

	items := make([]analyticsentity.BudgetSnapshot, 0)
	for _, item := range repository.budgets {
		if item.Month.Equal(month) {
			item.UserID = userID
			items = append(items, item)
		}
	}

	return items, nil
}

func (repository *inMemoryAnalyticsRepository) GetLatestScore(_ context.Context, userID string) (*analyticsentity.FinancialScore, error) {
	repository.mutex.RLock()
	defer repository.mutex.RUnlock()

	if repository.score == nil {
		return nil, nil
	}

	value := *repository.score
	value.UserID = userID
	return &value, nil
}

func (repository *inMemoryAnalyticsRepository) GetForecastByMonth(_ context.Context, userID string, referenceMonth time.Time) (*analyticsentity.CashflowForecast, error) {
	repository.mutex.RLock()
	defer repository.mutex.RUnlock()

	if repository.forecast == nil || !repository.forecast.ReferenceMonth.Equal(referenceMonth) {
		return nil, nil
	}

	value := *repository.forecast
	value.UserID = userID
	return &value, nil
}

func (repository *inMemoryAnalyticsRepository) ListInsights(_ context.Context, userID string, limit int, now time.Time) ([]analyticsentity.Insight, error) {
	repository.mutex.RLock()
	defer repository.mutex.RUnlock()

	items := make([]analyticsentity.Insight, 0)
	for _, item := range repository.insights {
		if item.ValidUntil != nil && item.ValidUntil.Before(now) {
			continue
		}

		item.UserID = userID
		items = append(items, item)
		if len(items) == limit {
			break
		}
	}

	return items, nil
}

func (repository *inMemoryAnalyticsRepository) ListAnomalies(_ context.Context, userID string, limit int) ([]analyticsentity.Anomaly, error) {
	repository.mutex.RLock()
	defer repository.mutex.RUnlock()

	items := make([]analyticsentity.Anomaly, 0)
	for _, item := range repository.anomalies {
		item.UserID = userID
		items = append(items, item)
		if len(items) == limit {
			break
		}
	}

	return items, nil
}

func (repository *inMemoryAnalyticsRepository) ListGoals(_ context.Context, userID string) ([]analyticsentity.FinancialGoal, error) {
	repository.mutex.RLock()
	defer repository.mutex.RUnlock()

	items := make([]analyticsentity.FinancialGoal, 0)
	for _, item := range repository.goals {
		if item.UserID == userID {
			items = append(items, item)
		}
	}

	return items, nil
}

func (repository *inMemoryAnalyticsRepository) GetGoalByID(_ context.Context, userID string, goalID string) (analyticsentity.FinancialGoal, error) {
	repository.mutex.RLock()
	defer repository.mutex.RUnlock()

	item, ok := repository.goals[goalID]
	if !ok || item.UserID != userID {
		return analyticsentity.FinancialGoal{}, analyticsservice.ErrGoalNotFound
	}

	return item, nil
}

func (repository *inMemoryAnalyticsRepository) CreateGoal(_ context.Context, goal analyticsentity.FinancialGoal) (analyticsentity.FinancialGoal, error) {
	repository.mutex.Lock()
	defer repository.mutex.Unlock()

	repository.goals[goal.ID] = goal
	return goal, nil
}

func (repository *inMemoryAnalyticsRepository) UpdateGoal(_ context.Context, goal analyticsentity.FinancialGoal) (analyticsentity.FinancialGoal, error) {
	repository.mutex.Lock()
	defer repository.mutex.Unlock()

	if _, ok := repository.goals[goal.ID]; !ok {
		return analyticsentity.FinancialGoal{}, analyticsservice.ErrGoalNotFound
	}

	repository.goals[goal.ID] = goal
	return goal, nil
}

var _ analyticsrepository.AnalyticsRepository = (*inMemoryAnalyticsRepository)(nil)
