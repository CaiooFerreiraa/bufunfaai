package repository

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	analyticsservice "github.com/bufunfaai/bufunfaai/apps/api/internal/modules/analytics/application/service"
	"github.com/bufunfaai/bufunfaai/apps/api/internal/modules/analytics/domain/entity"
)

type PostgresAnalyticsRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresAnalyticsRepository(pool *pgxpool.Pool) *PostgresAnalyticsRepository {
	return &PostgresAnalyticsRepository{pool: pool}
}

func (repository *PostgresAnalyticsRepository) ListBudgetSnapshotsByMonth(ctx context.Context, userID string, month time.Time) ([]entity.BudgetSnapshot, error) {
	rows, err := repository.pool.Query(ctx, `
		SELECT id, user_id, month, category, spent_amount_cents, budget_amount_cents, usage_percent, status, updated_at
		FROM budget_snapshots
		WHERE user_id = $1 AND month = $2
		ORDER BY usage_percent DESC, category ASC
	`, userID, month)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]entity.BudgetSnapshot, 0)
	for rows.Next() {
		item, scanErr := scanBudgetSnapshot(rows)
		if scanErr != nil {
			return nil, scanErr
		}

		items = append(items, item)
	}

	return items, rows.Err()
}

func (repository *PostgresAnalyticsRepository) GetLatestScore(ctx context.Context, userID string) (*entity.FinancialScore, error) {
	row := repository.pool.QueryRow(ctx, `
		SELECT id, user_id, period_start, period_end, total_score, score_label, factors_json, created_at
		FROM financial_scores
		WHERE user_id = $1
		ORDER BY period_end DESC, created_at DESC
		LIMIT 1
	`, userID)

	item, err := scanFinancialScore(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return &item, nil
}

func (repository *PostgresAnalyticsRepository) GetForecastByMonth(ctx context.Context, userID string, referenceMonth time.Time) (*entity.CashflowForecast, error) {
	row := repository.pool.QueryRow(ctx, `
		SELECT id, user_id, reference_month, projected_income_cents, projected_expense_cents, projected_balance_cents, confidence_level, assumptions_json, created_at
		FROM cashflow_forecasts
		WHERE user_id = $1 AND reference_month = $2
		ORDER BY created_at DESC
		LIMIT 1
	`, userID, referenceMonth)

	item, err := scanCashflowForecast(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return &item, nil
}

func (repository *PostgresAnalyticsRepository) ListInsights(ctx context.Context, userID string, limit int, now time.Time) ([]entity.Insight, error) {
	rows, err := repository.pool.Query(ctx, `
		SELECT id, user_id, insight_type, title, message, priority, source, valid_until, created_at
		FROM insights
		WHERE user_id = $1 AND (valid_until IS NULL OR valid_until >= $2)
		ORDER BY priority DESC, created_at DESC
		LIMIT $3
	`, userID, now, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]entity.Insight, 0)
	for rows.Next() {
		item, scanErr := scanInsight(rows)
		if scanErr != nil {
			return nil, scanErr
		}

		items = append(items, item)
	}

	return items, rows.Err()
}

func (repository *PostgresAnalyticsRepository) ListAnomalies(ctx context.Context, userID string, limit int) ([]entity.Anomaly, error) {
	rows, err := repository.pool.Query(ctx, `
		SELECT id, user_id, transaction_id, anomaly_type, severity, reason, detected_at
		FROM anomalies
		WHERE user_id = $1
		ORDER BY detected_at DESC
		LIMIT $2
	`, userID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]entity.Anomaly, 0)
	for rows.Next() {
		item, scanErr := scanAnomaly(rows)
		if scanErr != nil {
			return nil, scanErr
		}

		items = append(items, item)
	}

	return items, rows.Err()
}

func (repository *PostgresAnalyticsRepository) ListGoals(ctx context.Context, userID string) ([]entity.FinancialGoal, error) {
	rows, err := repository.pool.Query(ctx, `
		SELECT id, user_id, title, goal_type, target_amount_cents, current_amount_cents, due_date, status, created_at, updated_at
		FROM financial_goals
		WHERE user_id = $1
		ORDER BY created_at DESC
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]entity.FinancialGoal, 0)
	for rows.Next() {
		item, scanErr := scanFinancialGoal(rows)
		if scanErr != nil {
			return nil, scanErr
		}

		items = append(items, item)
	}

	return items, rows.Err()
}

func (repository *PostgresAnalyticsRepository) GetGoalByID(ctx context.Context, userID string, goalID string) (entity.FinancialGoal, error) {
	row := repository.pool.QueryRow(ctx, `
		SELECT id, user_id, title, goal_type, target_amount_cents, current_amount_cents, due_date, status, created_at, updated_at
		FROM financial_goals
		WHERE id = $1 AND user_id = $2
	`, goalID, userID)

	return scanFinancialGoal(row)
}

func (repository *PostgresAnalyticsRepository) CreateGoal(ctx context.Context, goal entity.FinancialGoal) (entity.FinancialGoal, error) {
	row := repository.pool.QueryRow(ctx, `
		INSERT INTO financial_goals (
			id, user_id, title, goal_type, target_amount_cents, current_amount_cents, due_date, status, created_at, updated_at
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)
		RETURNING id, user_id, title, goal_type, target_amount_cents, current_amount_cents, due_date, status, created_at, updated_at
	`, goal.ID, goal.UserID, goal.Title, goal.GoalType, goal.TargetAmountCents, goal.CurrentAmountCents, goal.DueDate, goal.Status, goal.CreatedAt, goal.UpdatedAt)

	return scanFinancialGoal(row)
}

func (repository *PostgresAnalyticsRepository) UpdateGoal(ctx context.Context, goal entity.FinancialGoal) (entity.FinancialGoal, error) {
	row := repository.pool.QueryRow(ctx, `
		UPDATE financial_goals
		SET title = $3, goal_type = $4, target_amount_cents = $5, current_amount_cents = $6, due_date = $7, status = $8, updated_at = $9
		WHERE id = $1 AND user_id = $2
		RETURNING id, user_id, title, goal_type, target_amount_cents, current_amount_cents, due_date, status, created_at, updated_at
	`, goal.ID, goal.UserID, goal.Title, goal.GoalType, goal.TargetAmountCents, goal.CurrentAmountCents, goal.DueDate, goal.Status, goal.UpdatedAt)

	return scanFinancialGoal(row)
}

type scanner interface {
	Scan(dest ...any) error
}

func scanBudgetSnapshot(row scanner) (entity.BudgetSnapshot, error) {
	var item entity.BudgetSnapshot
	err := row.Scan(&item.ID, &item.UserID, &item.Month, &item.Category, &item.SpentAmountCents, &item.BudgetAmountCents, &item.UsagePercent, &item.Status, &item.UpdatedAt)
	if err != nil {
		return entity.BudgetSnapshot{}, err
	}

	return item, nil
}

func scanFinancialScore(row scanner) (entity.FinancialScore, error) {
	var item entity.FinancialScore
	err := row.Scan(&item.ID, &item.UserID, &item.PeriodStart, &item.PeriodEnd, &item.TotalScore, &item.ScoreLabel, &item.FactorsJSON, &item.CreatedAt)
	if err != nil {
		return entity.FinancialScore{}, err
	}

	return item, nil
}

func scanCashflowForecast(row scanner) (entity.CashflowForecast, error) {
	var item entity.CashflowForecast
	err := row.Scan(&item.ID, &item.UserID, &item.ReferenceMonth, &item.ProjectedIncomeCents, &item.ProjectedExpenseCents, &item.ProjectedBalanceCents, &item.ConfidenceLevel, &item.AssumptionsJSON, &item.CreatedAt)
	if err != nil {
		return entity.CashflowForecast{}, err
	}

	return item, nil
}

func scanInsight(row scanner) (entity.Insight, error) {
	var item entity.Insight
	err := row.Scan(&item.ID, &item.UserID, &item.InsightType, &item.Title, &item.Message, &item.Priority, &item.Source, &item.ValidUntil, &item.CreatedAt)
	if err != nil {
		return entity.Insight{}, err
	}

	return item, nil
}

func scanAnomaly(row scanner) (entity.Anomaly, error) {
	var item entity.Anomaly
	err := row.Scan(&item.ID, &item.UserID, &item.TransactionID, &item.AnomalyType, &item.Severity, &item.Reason, &item.DetectedAt)
	if err != nil {
		return entity.Anomaly{}, err
	}

	return item, nil
}

func scanFinancialGoal(row scanner) (entity.FinancialGoal, error) {
	var item entity.FinancialGoal
	err := row.Scan(&item.ID, &item.UserID, &item.Title, &item.GoalType, &item.TargetAmountCents, &item.CurrentAmountCents, &item.DueDate, &item.Status, &item.CreatedAt, &item.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.FinancialGoal{}, analyticsservice.ErrGoalNotFound
		}

		return entity.FinancialGoal{}, err
	}

	return item, nil
}
