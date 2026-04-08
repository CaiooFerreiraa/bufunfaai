package dto

type CreateConsentRequest struct {
	InstitutionID string   `json:"institution_id" validate:"required,uuid4"`
	Purpose       string   `json:"purpose" validate:"required,min=5,max=255"`
	Permissions   []string `json:"permissions" validate:"required,min=1,dive,min=3,max=80"`
	RedirectURI   string   `json:"redirect_uri" validate:"required,url,max=500"`
}

type CallbackRequest struct {
	State string `json:"state" validate:"required,min=8,max=255"`
	Code  string `json:"code" validate:"required,min=4,max=255"`
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
