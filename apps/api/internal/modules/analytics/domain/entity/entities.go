package entity

import "time"

const (
	GoalStatusInProgress = "in_progress"
	GoalStatusCompleted  = "completed"
	GoalStatusAtRisk     = "at_risk"
)

type TransactionClassification struct {
	ID            string
	UserID        string
	TransactionID string
	Category      string
	Subcategory   string
	Confidence    float64
	Method        string
	ModelVersion  string
	PromptVersion string
	UserCorrected bool
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type CategoryBudget struct {
	ID                   string
	UserID               string
	Category             string
	MonthlyBudgetCents   int64
	SuggestedBudgetCents int64
	Source               string
	CreatedAt            time.Time
	UpdatedAt            time.Time
}

type BudgetSnapshot struct {
	ID                string
	UserID            string
	Month             time.Time
	Category          string
	SpentAmountCents  int64
	BudgetAmountCents int64
	UsagePercent      float64
	Status            string
	UpdatedAt         time.Time
}

type FinancialScore struct {
	ID          string
	UserID      string
	PeriodStart time.Time
	PeriodEnd   time.Time
	TotalScore  int
	ScoreLabel  string
	FactorsJSON string
	CreatedAt   time.Time
}

type Anomaly struct {
	ID            string
	UserID        string
	TransactionID string
	AnomalyType   string
	Severity      string
	Reason        string
	DetectedAt    time.Time
}

type CashflowForecast struct {
	ID                    string
	UserID                string
	ReferenceMonth        time.Time
	ProjectedIncomeCents  int64
	ProjectedExpenseCents int64
	ProjectedBalanceCents int64
	ConfidenceLevel       string
	AssumptionsJSON       string
	CreatedAt             time.Time
}

type FinancialGoal struct {
	ID                 string
	UserID             string
	Title              string
	GoalType           string
	TargetAmountCents  int64
	CurrentAmountCents int64
	DueDate            *time.Time
	Status             string
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

type Insight struct {
	ID          string
	UserID      string
	InsightType string
	Title       string
	Message     string
	Priority    int
	Source      string
	ValidUntil  *time.Time
	CreatedAt   time.Time
}

type ReportExport struct {
	ID          string
	UserID      string
	Format      string
	PeriodStart time.Time
	PeriodEnd   time.Time
	FileURL     string
	CreatedAt   time.Time
}

type Overview struct {
	ReferenceMonth string
	Budgets        []BudgetSnapshot
	Score          *FinancialScore
	Forecast       *CashflowForecast
	Insights       []Insight
	Anomalies      []Anomaly
	Goals          []FinancialGoal
}
