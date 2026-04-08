package dto

type CreateConsentRequest struct {
	InstitutionID string   `json:"institution_id" validate:"required,uuid"`
	Purpose       string   `json:"purpose" validate:"required,min=5,max=255"`
	Permissions   []string `json:"permissions" validate:"required,min=1,dive,min=3,max=80"`
	RedirectURI   string   `json:"redirect_uri" validate:"required,url,max=500"`
}

type CallbackRequest struct {
	State string `json:"state" validate:"required,min=8,max=255"`
	Code  string `json:"code" validate:"required,min=4,max=255"`
}

type CompleteConsentRequest struct {
	ItemID string `json:"item_id" validate:"required,uuid"`
}

type RevokeConsentRequest struct {
	Reason string `json:"reason" validate:"omitempty,max=255"`
}

type InstitutionOutput struct {
	ID                     string `json:"id"`
	DirectoryOrgID         string `json:"directory_org_id"`
	BrandName              string `json:"brand_name"`
	DisplayName            string `json:"display_name"`
	AuthorisationServerURL string `json:"authorisation_server_url"`
	ResourcesBaseURL       string `json:"resources_base_url"`
	LogoURL                string `json:"logo_url"`
	Status                 string `json:"status"`
	SupportsDataSharing    bool   `json:"supports_data_sharing"`
	SupportsPayments       bool   `json:"supports_payments"`
}

type ConsentOutput struct {
	ID                string   `json:"id"`
	UserID            string   `json:"user_id"`
	InstitutionID     string   `json:"institution_id"`
	ExternalConsentID string   `json:"external_consent_id"`
	Status            string   `json:"status"`
	Purpose           string   `json:"purpose"`
	Permissions       []string `json:"permissions"`
	ExpirationAt      string   `json:"expiration_at,omitempty"`
	RevokedAt         string   `json:"revoked_at,omitempty"`
	AuthorisedAt      string   `json:"authorised_at,omitempty"`
	RedirectURI       string   `json:"redirect_uri"`
	CreatedAt         string   `json:"created_at"`
	UpdatedAt         string   `json:"updated_at"`
}

type AuthorizationURLOutput struct {
	ConsentID        string `json:"consent_id"`
	AuthorizationURL string `json:"authorization_url"`
}

type ConnectTokenOutput struct {
	ConsentID           string `json:"consent_id"`
	ConnectToken        string `json:"connect_token"`
	SelectedConnectorID int64  `json:"selected_connector_id"`
}

type ConnectionOutput struct {
	ID                   string `json:"id"`
	UserID               string `json:"user_id"`
	InstitutionID        string `json:"institution_id"`
	ConsentID            string `json:"consent_id"`
	Status               string `json:"status"`
	FirstSyncAt          string `json:"first_sync_at,omitempty"`
	LastSyncAt           string `json:"last_sync_at,omitempty"`
	LastSuccessfulSyncAt string `json:"last_successful_sync_at,omitempty"`
	LastErrorCode        string `json:"last_error_code,omitempty"`
	LastErrorMessage     string `json:"last_error_message_redacted,omitempty"`
	CreatedAt            string `json:"created_at"`
	UpdatedAt            string `json:"updated_at"`
}

type SyncJobOutput struct {
	ID                   string `json:"id"`
	ConnectionID         string `json:"connection_id"`
	ResourceType         string `json:"resource_type"`
	Status               string `json:"status"`
	AttemptCount         int    `json:"attempt_count"`
	ScheduledAt          string `json:"scheduled_at,omitempty"`
	StartedAt            string `json:"started_at,omitempty"`
	FinishedAt           string `json:"finished_at,omitempty"`
	ErrorCode            string `json:"error_code,omitempty"`
	ErrorMessageRedacted string `json:"error_message_redacted,omitempty"`
}

type SyncStatusOutput struct {
	Connection ConnectionOutput `json:"connection"`
	Jobs       []SyncJobOutput  `json:"jobs"`
}

type CallbackResultOutput struct {
	Consent    ConsentOutput    `json:"consent"`
	Connection ConnectionOutput `json:"connection"`
}

type ReconciliationResultOutput struct {
	Processed   int `json:"processed"`
	Successful  int `json:"successful"`
	Failed      int `json:"failed"`
	JobsCreated int `json:"jobs_created"`
}

type TransactionsQuery struct {
	From  string `form:"from" validate:"omitempty,datetime=2006-01-02"`
	To    string `form:"to" validate:"omitempty,datetime=2006-01-02"`
	Limit int    `form:"limit" validate:"omitempty,min=1,max=500"`
}

type AccountSnapshotOutput struct {
	ID                   string  `json:"id"`
	ConnectionID         string  `json:"connection_id"`
	InstitutionID        string  `json:"institution_id"`
	InstitutionName      string  `json:"institution_name"`
	ItemID               string  `json:"item_id"`
	Type                 string  `json:"type"`
	Subtype              string  `json:"subtype,omitempty"`
	Name                 string  `json:"name"`
	MarketingName        string  `json:"marketing_name,omitempty"`
	Number               string  `json:"number,omitempty"`
	CurrencyCode         string  `json:"currency_code"`
	Balance              float64 `json:"balance"`
	BankTransferNumber   string  `json:"bank_transfer_number,omitempty"`
	CreditBrand          string  `json:"credit_brand,omitempty"`
	AvailableCreditLimit float64 `json:"available_credit_limit,omitempty"`
}

type TransactionSnapshotOutput struct {
	ID              string  `json:"id"`
	AccountID       string  `json:"account_id"`
	ConnectionID    string  `json:"connection_id"`
	InstitutionName string  `json:"institution_name"`
	AccountName     string  `json:"account_name"`
	Description     string  `json:"description"`
	Category        string  `json:"category"`
	Type            string  `json:"type"`
	Status          string  `json:"status"`
	CurrencyCode    string  `json:"currency_code"`
	Amount          float64 `json:"amount"`
	SignedAmount    float64 `json:"signed_amount"`
	Date            string  `json:"date"`
	MerchantName    string  `json:"merchant_name,omitempty"`
}

type CategoryBreakdownOutput struct {
	Category string  `json:"category"`
	Amount   float64 `json:"amount"`
	Percent  float64 `json:"percent"`
}

type OverviewOutput struct {
	Accounts            []AccountSnapshotOutput     `json:"accounts"`
	RecentTransactions  []TransactionSnapshotOutput `json:"recent_transactions"`
	ExpenseBreakdown    []CategoryBreakdownOutput   `json:"expense_breakdown"`
	TotalAvailable      float64                     `json:"total_available"`
	CreditCardBalance   float64                     `json:"credit_card_balance"`
	MonthIncome         float64                     `json:"month_income"`
	MonthExpenses       float64                     `json:"month_expenses"`
	ConnectedAccounts   int                         `json:"connected_accounts"`
	ConnectionsWithData int                         `json:"connections_with_data"`
}

type TransactionFeedOutput struct {
	Transactions              []TransactionSnapshotOutput `json:"transactions"`
	ExpenseBreakdown          []CategoryBreakdownOutput   `json:"expense_breakdown"`
	MonthExpenseTotal         float64                     `json:"month_expense_total"`
	PreviousMonthExpenseTotal float64                     `json:"previous_month_expense_total"`
	MonthIncomeTotal          float64                     `json:"month_income_total"`
}
