package presenter

import (
	"time"

	analyticsdto "github.com/bufunfaai/bufunfaai/apps/api/internal/modules/analytics/application/dto"
	"github.com/bufunfaai/bufunfaai/apps/api/internal/modules/analytics/domain/entity"
)

func OverviewOutput(overview entity.Overview) analyticsdto.OverviewOutput {
	budgets := make([]analyticsdto.BudgetSnapshotOutput, 0, len(overview.Budgets))
	for _, item := range overview.Budgets {
		budgets = append(budgets, BudgetSnapshotOutput(item))
	}

	insights := make([]analyticsdto.InsightOutput, 0, len(overview.Insights))
	for _, item := range overview.Insights {
		insights = append(insights, InsightOutput(item))
	}

	anomalies := make([]analyticsdto.AnomalyOutput, 0, len(overview.Anomalies))
	for _, item := range overview.Anomalies {
		anomalies = append(anomalies, AnomalyOutput(item))
	}

	goals := make([]analyticsdto.GoalOutput, 0, len(overview.Goals))
	for _, item := range overview.Goals {
		goals = append(goals, GoalOutput(item))
	}

	var scoreOutput *analyticsdto.ScoreOutput
	if overview.Score != nil {
		score := ScoreOutput(*overview.Score)
		scoreOutput = &score
	}

	var forecastOutput *analyticsdto.ForecastOutput
	if overview.Forecast != nil {
		forecast := ForecastOutput(*overview.Forecast)
		forecastOutput = &forecast
	}

	return analyticsdto.OverviewOutput{
		ReferenceMonth: overview.ReferenceMonth,
		Budgets:        budgets,
		Score:          scoreOutput,
		Forecast:       forecastOutput,
		Insights:       insights,
		Anomalies:      anomalies,
		Goals:          goals,
	}
}

func BudgetSnapshotOutput(item entity.BudgetSnapshot) analyticsdto.BudgetSnapshotOutput {
	return analyticsdto.BudgetSnapshotOutput{
		ID:                item.ID,
		Month:             item.Month.Format("2006-01"),
		Category:          item.Category,
		SpentAmountCents:  item.SpentAmountCents,
		BudgetAmountCents: item.BudgetAmountCents,
		UsagePercent:      item.UsagePercent,
		Status:            item.Status,
		UpdatedAt:         item.UpdatedAt.Format(time.RFC3339),
	}
}

func ScoreOutput(item entity.FinancialScore) analyticsdto.ScoreOutput {
	return analyticsdto.ScoreOutput{
		ID:          item.ID,
		PeriodStart: item.PeriodStart.Format("2006-01-02"),
		PeriodEnd:   item.PeriodEnd.Format("2006-01-02"),
		TotalScore:  item.TotalScore,
		ScoreLabel:  item.ScoreLabel,
		FactorsJSON: item.FactorsJSON,
		CreatedAt:   item.CreatedAt.Format(time.RFC3339),
	}
}

func ForecastOutput(item entity.CashflowForecast) analyticsdto.ForecastOutput {
	return analyticsdto.ForecastOutput{
		ID:                    item.ID,
		ReferenceMonth:        item.ReferenceMonth.Format("2006-01"),
		ProjectedIncomeCents:  item.ProjectedIncomeCents,
		ProjectedExpenseCents: item.ProjectedExpenseCents,
		ProjectedBalanceCents: item.ProjectedBalanceCents,
		ConfidenceLevel:       item.ConfidenceLevel,
		AssumptionsJSON:       item.AssumptionsJSON,
		CreatedAt:             item.CreatedAt.Format(time.RFC3339),
	}
}

func InsightOutput(item entity.Insight) analyticsdto.InsightOutput {
	return analyticsdto.InsightOutput{
		ID:          item.ID,
		InsightType: item.InsightType,
		Title:       item.Title,
		Message:     item.Message,
		Priority:    item.Priority,
		Source:      item.Source,
		ValidUntil:  formatTimePointer(item.ValidUntil),
		CreatedAt:   item.CreatedAt.Format(time.RFC3339),
	}
}

func AnomalyOutput(item entity.Anomaly) analyticsdto.AnomalyOutput {
	return analyticsdto.AnomalyOutput{
		ID:            item.ID,
		TransactionID: item.TransactionID,
		AnomalyType:   item.AnomalyType,
		Severity:      item.Severity,
		Reason:        item.Reason,
		DetectedAt:    item.DetectedAt.Format(time.RFC3339),
	}
}

func GoalOutput(item entity.FinancialGoal) analyticsdto.GoalOutput {
	return analyticsdto.GoalOutput{
		ID:                 item.ID,
		Title:              item.Title,
		GoalType:           item.GoalType,
		TargetAmountCents:  item.TargetAmountCents,
		CurrentAmountCents: item.CurrentAmountCents,
		DueDate:            formatTimePointer(item.DueDate),
		Status:             item.Status,
		CreatedAt:          item.CreatedAt.Format(time.RFC3339),
		UpdatedAt:          item.UpdatedAt.Format(time.RFC3339),
	}
}

func formatTimePointer(value *time.Time) string {
	if value == nil {
		return ""
	}

	return value.Format(time.RFC3339)
}
