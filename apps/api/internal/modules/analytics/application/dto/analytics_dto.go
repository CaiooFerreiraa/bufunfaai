package dto

type BudgetSnapshotOutput struct {
	ID                string  `json:"id"`
	Month             string  `json:"month"`
	Category          string  `json:"category"`
	SpentAmountCents  int64   `json:"spent_amount_cents"`
	BudgetAmountCents int64   `json:"budget_amount_cents"`
	UsagePercent      float64 `json:"usage_percent"`
	Status            string  `json:"status"`
	UpdatedAt         string  `json:"updated_at"`
}

type ScoreOutput struct {
	ID          string `json:"id"`
	PeriodStart string `json:"period_start"`
	PeriodEnd   string `json:"period_end"`
	TotalScore  int    `json:"total_score"`
	ScoreLabel  string `json:"score_label"`
	FactorsJSON string `json:"factors_json"`
	CreatedAt   string `json:"created_at"`
}

type ForecastOutput struct {
	ID                    string `json:"id"`
	ReferenceMonth        string `json:"reference_month"`
	ProjectedIncomeCents  int64  `json:"projected_income_cents"`
	ProjectedExpenseCents int64  `json:"projected_expense_cents"`
	ProjectedBalanceCents int64  `json:"projected_balance_cents"`
	ConfidenceLevel       string `json:"confidence_level"`
	AssumptionsJSON       string `json:"assumptions_json"`
	CreatedAt             string `json:"created_at"`
}

type InsightOutput struct {
	ID          string `json:"id"`
	InsightType string `json:"insight_type"`
	Title       string `json:"title"`
	Message     string `json:"message"`
	Priority    int    `json:"priority"`
	Source      string `json:"source"`
	ValidUntil  string `json:"valid_until"`
	CreatedAt   string `json:"created_at"`
}

type AnomalyOutput struct {
	ID            string `json:"id"`
	TransactionID string `json:"transaction_id"`
	AnomalyType   string `json:"anomaly_type"`
	Severity      string `json:"severity"`
	Reason        string `json:"reason"`
	DetectedAt    string `json:"detected_at"`
}

type GoalOutput struct {
	ID                 string `json:"id"`
	Title              string `json:"title"`
	GoalType           string `json:"goal_type"`
	TargetAmountCents  int64  `json:"target_amount_cents"`
	CurrentAmountCents int64  `json:"current_amount_cents"`
	DueDate            string `json:"due_date"`
	Status             string `json:"status"`
	CreatedAt          string `json:"created_at"`
	UpdatedAt          string `json:"updated_at"`
}

type OverviewOutput struct {
	ReferenceMonth string                 `json:"reference_month"`
	Budgets        []BudgetSnapshotOutput `json:"budgets"`
	Score          *ScoreOutput           `json:"score,omitempty"`
	Forecast       *ForecastOutput        `json:"forecast,omitempty"`
	Insights       []InsightOutput        `json:"insights"`
	Anomalies      []AnomalyOutput        `json:"anomalies"`
	Goals          []GoalOutput           `json:"goals"`
}

type CreateGoalRequest struct {
	Title              string `json:"title" validate:"required,min=3,max=120"`
	GoalType           string `json:"goal_type" validate:"required,min=3,max=60"`
	TargetAmountCents  int64  `json:"target_amount_cents" validate:"required,gt=0"`
	CurrentAmountCents int64  `json:"current_amount_cents" validate:"gte=0"`
	DueDate            string `json:"due_date" validate:"omitempty,datetime=2006-01-02"`
	Status             string `json:"status" validate:"omitempty,oneof=in_progress completed at_risk"`
}

type UpdateGoalRequest struct {
	Title              *string `json:"title" validate:"omitempty,min=3,max=120"`
	GoalType           *string `json:"goal_type" validate:"omitempty,min=3,max=60"`
	TargetAmountCents  *int64  `json:"target_amount_cents" validate:"omitempty,gt=0"`
	CurrentAmountCents *int64  `json:"current_amount_cents" validate:"omitempty,gte=0"`
	DueDate            *string `json:"due_date" validate:"omitempty,datetime=2006-01-02"`
	Status             *string `json:"status" validate:"omitempty,oneof=in_progress completed at_risk"`
}
