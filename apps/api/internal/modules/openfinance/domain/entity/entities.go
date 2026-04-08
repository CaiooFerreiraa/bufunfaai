package entity

import "time"

const (
	ConsentStatusPending        string = "PENDING"
	ConsentStatusAuthURLReady   string = "AUTH_URL_READY"
	ConsentStatusAuthInProgress string = "AUTH_IN_PROGRESS"
	ConsentStatusAuthorised     string = "AUTHORISED"
	ConsentStatusDenied         string = "DENIED"
	ConsentStatusRevoked        string = "REVOKED"
	ConsentStatusExpired        string = "EXPIRED"
	ConsentStatusError          string = "ERROR"

	ConnectionStatusPending        string = "PENDING"
	ConnectionStatusActive         string = "ACTIVE"
	ConnectionStatusTokenExpired   string = "TOKEN_EXPIRED"
	ConnectionStatusConsentExpired string = "CONSENT_EXPIRED"
	ConnectionStatusRevoked        string = "REVOKED"
	ConnectionStatusSyncError      string = "SYNC_ERROR"
	ConnectionStatusReauthRequired string = "REAUTH_REQUIRED"

	SyncJobStatusPending   string = "PENDING"
	SyncJobStatusRunning   string = "RUNNING"
	SyncJobStatusCompleted string = "COMPLETED"
	SyncJobStatusError     string = "ERROR"

	ResourceAccounts     string = "accounts"
	ResourceBalances     string = "balances"
	ResourceTransactions string = "transactions"
)

type Institution struct {
	ID                     string
	DirectoryOrgID         string
	BrandName              string
	DisplayName            string
	AuthorisationServerURL string
	ResourcesBaseURL       string
	LogoURL                string
	Status                 string
	SupportsDataSharing    bool
	SupportsPayments       bool
	LastDirectorySyncAt    *time.Time
	CreatedAt              time.Time
	UpdatedAt              time.Time
}

type Consent struct {
	ID                    string
	UserID                string
	InstitutionID         string
	ExternalConsentID     string
	Status                string
	Purpose               string
	PermissionsJSON       string
	ExpirationAt          *time.Time
	RevokedAt             *time.Time
	AuthorisedAt          *time.Time
	DeniedAt              *time.Time
	RedirectURI           string
	State                 string
	Nonce                 string
	CodeVerifierEncrypted string
	CreatedAt             time.Time
	UpdatedAt             time.Time
}

type Authorization struct {
	ID                         string
	ConsentID                  string
	AuthorizationCodeHash      string
	AuthorizationCodeExpiresAt *time.Time
	PKceMethod                 string
	RedirectReceivedAt         *time.Time
	CreatedAt                  time.Time
}

type Token struct {
	ID                    string
	ConsentID             string
	InstitutionID         string
	AccessTokenEncrypted  string
	RefreshTokenEncrypted string
	TokenType             string
	Scope                 string
	ExpiresAt             time.Time
	RefreshExpiresAt      *time.Time
	LastRefreshAt         *time.Time
	RevokedAt             *time.Time
	CreatedAt             time.Time
	UpdatedAt             time.Time
}

type Connection struct {
	ID                       string
	UserID                   string
	InstitutionID            string
	ConsentID                string
	Status                   string
	FirstSyncAt              *time.Time
	LastSyncAt               *time.Time
	LastSuccessfulSyncAt     *time.Time
	LastErrorCode            string
	LastErrorMessageRedacted string
	CreatedAt                time.Time
	UpdatedAt                time.Time
}

type SyncJob struct {
	ID                   string
	ConnectionID         string
	ResourceType         string
	Status               string
	Cursor               string
	WindowStart          *time.Time
	WindowEnd            *time.Time
	AttemptCount         int
	ScheduledAt          *time.Time
	StartedAt            *time.Time
	FinishedAt           *time.Time
	ErrorCode            string
	ErrorMessageRedacted string
	CreatedAt            time.Time
}

type SyncCheckpoint struct {
	ID                    string
	ConnectionID          string
	ResourceType          string
	Cursor                string
	LastReferenceDatetime *time.Time
	ETag                  string
	UpdatedAt             time.Time
}
