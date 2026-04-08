CREATE TABLE IF NOT EXISTS transaction_classifications (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    transaction_id UUID NOT NULL,
    category TEXT NOT NULL,
    subcategory TEXT NOT NULL,
    confidence NUMERIC(5,4) NOT NULL,
    method TEXT NOT NULL,
    model_version TEXT NOT NULL DEFAULT '',
    prompt_version TEXT NOT NULL DEFAULT '',
    user_corrected BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_transaction_classifications_transaction_id
    ON transaction_classifications (transaction_id);

CREATE INDEX IF NOT EXISTS idx_transaction_classifications_user_id
    ON transaction_classifications (user_id);

CREATE TABLE IF NOT EXISTS category_budgets (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    category TEXT NOT NULL,
    monthly_budget_cents BIGINT NOT NULL,
    suggested_budget_cents BIGINT NOT NULL DEFAULT 0,
    source TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_category_budgets_user_category
    ON category_budgets (user_id, category);

CREATE TABLE IF NOT EXISTS budget_snapshots (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    month DATE NOT NULL,
    category TEXT NOT NULL,
    spent_amount_cents BIGINT NOT NULL,
    budget_amount_cents BIGINT NOT NULL,
    usage_percent NUMERIC(7,2) NOT NULL,
    status TEXT NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_budget_snapshots_user_month_category
    ON budget_snapshots (user_id, month, category);

CREATE TABLE IF NOT EXISTS financial_scores (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    period_start DATE NOT NULL,
    period_end DATE NOT NULL,
    total_score INTEGER NOT NULL,
    score_label TEXT NOT NULL,
    factors_json JSONB NOT NULL,
    created_at TIMESTAMPTZ NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_financial_scores_user_period_end
    ON financial_scores (user_id, period_end DESC);

CREATE TABLE IF NOT EXISTS anomalies (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    transaction_id UUID NOT NULL,
    anomaly_type TEXT NOT NULL,
    severity TEXT NOT NULL,
    reason TEXT NOT NULL,
    detected_at TIMESTAMPTZ NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_anomalies_user_detected_at
    ON anomalies (user_id, detected_at DESC);

CREATE TABLE IF NOT EXISTS cashflow_forecasts (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    reference_month DATE NOT NULL,
    projected_income_cents BIGINT NOT NULL,
    projected_expense_cents BIGINT NOT NULL,
    projected_balance_cents BIGINT NOT NULL,
    confidence_level TEXT NOT NULL,
    assumptions_json JSONB NOT NULL,
    created_at TIMESTAMPTZ NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_cashflow_forecasts_user_reference_month
    ON cashflow_forecasts (user_id, reference_month DESC);

CREATE TABLE IF NOT EXISTS financial_goals (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title TEXT NOT NULL,
    goal_type TEXT NOT NULL,
    target_amount_cents BIGINT NOT NULL,
    current_amount_cents BIGINT NOT NULL DEFAULT 0,
    due_date DATE NULL,
    status TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_financial_goals_user_created_at
    ON financial_goals (user_id, created_at DESC);

CREATE TABLE IF NOT EXISTS insights (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    insight_type TEXT NOT NULL,
    title TEXT NOT NULL,
    message TEXT NOT NULL,
    priority INTEGER NOT NULL DEFAULT 0,
    source TEXT NOT NULL,
    valid_until TIMESTAMPTZ NULL,
    created_at TIMESTAMPTZ NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_insights_user_priority
    ON insights (user_id, priority DESC, created_at DESC);

CREATE TABLE IF NOT EXISTS report_exports (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    format TEXT NOT NULL,
    period_start DATE NOT NULL,
    period_end DATE NOT NULL,
    file_url TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_report_exports_user_created_at
    ON report_exports (user_id, created_at DESC);
