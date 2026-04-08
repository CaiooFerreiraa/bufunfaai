package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"math"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"

	ofdto "github.com/bufunfaai/bufunfaai/apps/api/internal/modules/openfinance/application/dto"
	"github.com/bufunfaai/bufunfaai/apps/api/internal/modules/openfinance/domain/entity"
	ofrepository "github.com/bufunfaai/bufunfaai/apps/api/internal/modules/openfinance/domain/repository"
	sharederrors "github.com/bufunfaai/bufunfaai/apps/api/internal/shared/errors"
)

var (
	ErrInstitutionNotFound = errors.New("institution not found")
	ErrConsentNotFound     = errors.New("consent not found")
	ErrConnectionNotFound  = errors.New("connection not found")
	ErrTokenNotFound       = errors.New("token not found")
)

type Provider interface {
	ListInstitutions(ctx context.Context) ([]entity.Institution, error)
	CreateConsent(ctx context.Context, institution entity.Institution, consent entity.Consent, permissions []string) (string, *time.Time, error)
	BuildAuthorizationURL(ctx context.Context, institution entity.Institution, consent entity.Consent) (string, error)
	CreateConnectToken(ctx context.Context, institution entity.Institution, consent entity.Consent) (ProviderConnectToken, error)
	GetItem(ctx context.Context, itemID string) (ProviderItem, error)
	ListAccounts(ctx context.Context, itemID string) ([]ProviderAccount, error)
	ListTransactions(ctx context.Context, accountID string, query ProviderTransactionQuery) ([]ProviderTransaction, error)
	ExchangeCode(ctx context.Context, institution entity.Institution, consent entity.Consent, code string) (ProviderTokenSet, error)
	RevokeConsent(ctx context.Context, institution entity.Institution, consent entity.Consent) error
	SyncResources(ctx context.Context, institution entity.Institution, consent entity.Consent, connection entity.Connection) ([]ProviderSyncResult, error)
}

type Cipher interface {
	Encrypt(plaintext string) (string, error)
	Decrypt(ciphertext string) (string, error)
}

type ProviderTokenSet struct {
	AccessToken      string
	RefreshToken     string
	TokenType        string
	Scope            string
	ExpiresAt        time.Time
	RefreshExpiresAt *time.Time
}

type ProviderConnectToken struct {
	ConnectToken        string
	SelectedConnectorID int64
}

type ProviderItem struct {
	ID              string
	ConnectorID     int64
	Status          string
	ExecutionStatus string
	LastUpdatedAt   *time.Time
	ErrorCode       string
}

type ProviderSyncResult struct {
	ResourceType string
	Status       string
	ErrorCode    string
	ErrorMessage string
}

type ProviderTransactionQuery struct {
	From     *time.Time
	To       *time.Time
	PageSize int
}

type ProviderAccount struct {
	ID                   string
	ItemID               string
	Type                 string
	Subtype              string
	Number               string
	Name                 string
	MarketingName        string
	Balance              float64
	CurrencyCode         string
	BankTransferNumber   string
	CreditBrand          string
	AvailableCreditLimit float64
}

type ProviderTransaction struct {
	ID           string
	AccountID    string
	Description  string
	Amount       float64
	Date         time.Time
	CurrencyCode string
	Category     string
	Type         string
	Status       string
	MerchantName string
}

type AccountSnapshot struct {
	ID                   string
	ConnectionID         string
	InstitutionID        string
	InstitutionName      string
	ItemID               string
	Type                 string
	Subtype              string
	Name                 string
	MarketingName        string
	Number               string
	CurrencyCode         string
	Balance              float64
	BankTransferNumber   string
	CreditBrand          string
	AvailableCreditLimit float64
}

type TransactionSnapshot struct {
	ID              string
	AccountID       string
	ConnectionID    string
	InstitutionName string
	AccountName     string
	Description     string
	Category        string
	Type            string
	Status          string
	CurrencyCode    string
	Amount          float64
	Date            time.Time
	MerchantName    string
}

type CategoryBreakdown struct {
	Category string
	Amount   float64
	Percent  float64
}

type Overview struct {
	Accounts            []AccountSnapshot
	RecentTransactions  []TransactionSnapshot
	ExpenseBreakdown    []CategoryBreakdown
	TotalAvailable      float64
	CreditCardBalance   float64
	MonthIncome         float64
	MonthExpenses       float64
	ConnectedAccounts   int
	ConnectionsWithData int
}

type TransactionFeed struct {
	Transactions              []TransactionSnapshot
	ExpenseBreakdown          []CategoryBreakdown
	MonthExpenseTotal         float64
	PreviousMonthExpenseTotal float64
	MonthIncomeTotal          float64
}

type ReconciliationResult struct {
	Processed   int
	Successful  int
	Failed      int
	JobsCreated int
}

type Service struct {
	institutionRepository   ofrepository.InstitutionRepository
	consentRepository       ofrepository.ConsentRepository
	authorizationRepository ofrepository.AuthorizationRepository
	tokenRepository         ofrepository.TokenRepository
	connectionRepository    ofrepository.ConnectionRepository
	syncJobRepository       ofrepository.SyncJobRepository
	provider                Provider
	cipher                  Cipher
	now                     func() time.Time
}

func NewService(
	institutionRepository ofrepository.InstitutionRepository,
	consentRepository ofrepository.ConsentRepository,
	authorizationRepository ofrepository.AuthorizationRepository,
	tokenRepository ofrepository.TokenRepository,
	connectionRepository ofrepository.ConnectionRepository,
	syncJobRepository ofrepository.SyncJobRepository,
	provider Provider,
	cipher Cipher,
) *Service {
	return &Service{
		institutionRepository:   institutionRepository,
		consentRepository:       consentRepository,
		authorizationRepository: authorizationRepository,
		tokenRepository:         tokenRepository,
		connectionRepository:    connectionRepository,
		syncJobRepository:       syncJobRepository,
		provider:                provider,
		cipher:                  cipher,
		now:                     time.Now,
	}
}

func (service *Service) EnsureInstitutions(ctx context.Context) *sharederrors.AppError {
	institutions, err := service.institutionRepository.List(ctx)
	if err == nil && len(institutions) > 0 && hasSandboxInstitution(institutions) {
		return nil
	}

	discoveredInstitutions, err := service.provider.ListInstitutions(ctx)
	if err != nil {
		if len(institutions) > 0 {
			return nil
		}
		return sharederrors.Wrap("OPEN_FINANCE_DISCOVERY_ERROR", "erro ao sincronizar instituicoes", 500, err)
	}

	if err := service.institutionRepository.SaveMany(ctx, discoveredInstitutions); err != nil {
		return sharederrors.Wrap("OPEN_FINANCE_DISCOVERY_ERROR", "erro ao persistir instituicoes", 500, err)
	}

	return nil
}

func (service *Service) ListInstitutions(ctx context.Context) ([]entity.Institution, *sharederrors.AppError) {
	if appError := service.EnsureInstitutions(ctx); appError != nil {
		return nil, appError
	}

	institutions, err := service.institutionRepository.List(ctx)
	if err != nil {
		return nil, sharederrors.Wrap("OPEN_FINANCE_INSTITUTIONS_ERROR", "erro ao listar instituicoes", 500, err)
	}

	return institutions, nil
}

func (service *Service) GetInstitution(ctx context.Context, institutionID string) (entity.Institution, *sharederrors.AppError) {
	if appError := service.EnsureInstitutions(ctx); appError != nil {
		return entity.Institution{}, appError
	}

	institution, err := service.institutionRepository.GetByID(ctx, institutionID)
	if err != nil {
		if errors.Is(err, ErrInstitutionNotFound) {
			return entity.Institution{}, sharederrors.New("INSTITUTION_NOT_FOUND", "instituicao nao encontrada", 404)
		}

		return entity.Institution{}, sharederrors.Wrap("OPEN_FINANCE_INSTITUTION_ERROR", "erro ao buscar instituicao", 500, err)
	}

	return institution, nil
}

func (service *Service) CreateConsent(
	ctx context.Context,
	userID string,
	request ofdto.CreateConsentRequest,
) (entity.Consent, *sharederrors.AppError) {
	institution, appError := service.GetInstitution(ctx, request.InstitutionID)
	if appError != nil {
		return entity.Consent{}, appError
	}

	permissionsJSON, err := json.Marshal(request.Permissions)
	if err != nil {
		return entity.Consent{}, sharederrors.Wrap("OPEN_FINANCE_CONSENT_ERROR", "erro ao serializar permissoes", 500, err)
	}

	codeVerifier := uuid.NewString() + uuid.NewString()
	encryptedCodeVerifier, err := service.cipher.Encrypt(codeVerifier)
	if err != nil {
		return entity.Consent{}, sharederrors.Wrap("OPEN_FINANCE_CONSENT_ERROR", "erro ao proteger pkce", 500, err)
	}

	now := service.now().UTC()
	consent := entity.Consent{
		ID:                    uuid.NewString(),
		UserID:                userID,
		InstitutionID:         institution.ID,
		Status:                entity.ConsentStatusPending,
		Purpose:               request.Purpose,
		PermissionsJSON:       string(permissionsJSON),
		RedirectURI:           request.RedirectURI,
		State:                 uuid.NewString(),
		Nonce:                 uuid.NewString(),
		CodeVerifierEncrypted: encryptedCodeVerifier,
		CreatedAt:             now,
		UpdatedAt:             now,
	}

	externalConsentID, expirationAt, err := service.provider.CreateConsent(ctx, institution, consent, request.Permissions)
	if err != nil {
		return entity.Consent{}, sharederrors.Wrap("OPEN_FINANCE_CONSENT_ERROR", "erro ao criar consentimento externo", 502, err)
	}

	consent.ExternalConsentID = externalConsentID
	consent.ExpirationAt = expirationAt
	consent.Status = entity.ConsentStatusAuthURLReady

	if err := service.consentRepository.Create(ctx, consent); err != nil {
		return entity.Consent{}, sharederrors.Wrap("OPEN_FINANCE_CONSENT_ERROR", "erro ao persistir consentimento", 500, err)
	}

	return consent, nil
}

func (service *Service) GetConsent(ctx context.Context, consentID string, userID string) (entity.Consent, *sharederrors.AppError) {
	consent, err := service.consentRepository.GetByID(ctx, consentID)
	if err != nil {
		if errors.Is(err, ErrConsentNotFound) {
			return entity.Consent{}, sharederrors.New("CONSENT_NOT_FOUND", "consentimento nao encontrado", 404)
		}

		return entity.Consent{}, sharederrors.Wrap("OPEN_FINANCE_CONSENT_ERROR", "erro ao buscar consentimento", 500, err)
	}

	if consent.UserID != userID {
		return entity.Consent{}, sharederrors.New("CONSENT_NOT_FOUND", "consentimento nao encontrado", 404)
	}

	return consent, nil
}

func (service *Service) AuthorizeConsent(ctx context.Context, consentID string, userID string) (string, *sharederrors.AppError) {
	consent, appError := service.GetConsent(ctx, consentID, userID)
	if appError != nil {
		return "", appError
	}

	institution, appError := service.GetInstitution(ctx, consent.InstitutionID)
	if appError != nil {
		return "", appError
	}

	authorizationURL, err := service.provider.BuildAuthorizationURL(ctx, institution, consent)
	if err != nil {
		return "", sharederrors.Wrap("OPEN_FINANCE_AUTHORIZE_ERROR", "erro ao montar url de autorizacao", 502, err)
	}

	consent.Status = entity.ConsentStatusAuthInProgress
	consent.UpdatedAt = service.now().UTC()
	if err := service.consentRepository.Update(ctx, consent); err != nil {
		return "", sharederrors.Wrap("OPEN_FINANCE_AUTHORIZE_ERROR", "erro ao atualizar consentimento", 500, err)
	}

	return authorizationURL, nil
}

func (service *Service) CreateConnectToken(ctx context.Context, consentID string, userID string) (ProviderConnectToken, *sharederrors.AppError) {
	consent, appError := service.GetConsent(ctx, consentID, userID)
	if appError != nil {
		return ProviderConnectToken{}, appError
	}

	institution, appError := service.GetInstitution(ctx, consent.InstitutionID)
	if appError != nil {
		return ProviderConnectToken{}, appError
	}

	connectToken, err := service.provider.CreateConnectToken(ctx, institution, consent)
	if err != nil {
		return ProviderConnectToken{}, sharederrors.Wrap("OPEN_FINANCE_CONNECT_TOKEN_ERROR", "erro ao preparar conexão com o banco", 502, err)
	}

	consent.Status = entity.ConsentStatusAuthInProgress
	consent.UpdatedAt = service.now().UTC()
	if err := service.consentRepository.Update(ctx, consent); err != nil {
		return ProviderConnectToken{}, sharederrors.Wrap("OPEN_FINANCE_CONNECT_TOKEN_ERROR", "erro ao atualizar consentimento", 500, err)
	}

	return connectToken, nil
}

func (service *Service) CompleteConsent(ctx context.Context, consentID string, userID string, itemID string) (entity.Consent, entity.Connection, *sharederrors.AppError) {
	consent, appError := service.GetConsent(ctx, consentID, userID)
	if appError != nil {
		return entity.Consent{}, entity.Connection{}, appError
	}

	institution, appError := service.GetInstitution(ctx, consent.InstitutionID)
	if appError != nil {
		return entity.Consent{}, entity.Connection{}, appError
	}

	item, err := service.provider.GetItem(ctx, itemID)
	if err != nil {
		return entity.Consent{}, entity.Connection{}, sharederrors.Wrap("OPEN_FINANCE_CONNECTION_ERROR", "erro ao validar conexão com o banco", 502, err)
	}

	if err := validateProviderItem(institution, item); err != nil {
		return entity.Consent{}, entity.Connection{}, sharederrors.Wrap("OPEN_FINANCE_CONNECTION_ERROR", "erro ao validar banco selecionado", 400, err)
	}

	now := service.now().UTC()
	consent.ExternalConsentID = item.ID
	consent.Status = entity.ConsentStatusAuthorised
	consent.AuthorisedAt = &now
	consent.UpdatedAt = now
	if err := service.consentRepository.Update(ctx, consent); err != nil {
		return entity.Consent{}, entity.Connection{}, sharederrors.Wrap("OPEN_FINANCE_CONNECTION_ERROR", "erro ao atualizar consentimento", 500, err)
	}

	connectionStatus := mapConnectionStatus(item)
	connection := entity.Connection{
		ID:            uuid.NewString(),
		UserID:        consent.UserID,
		InstitutionID: consent.InstitutionID,
		ConsentID:     consent.ID,
		Status:        connectionStatus,
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	if connectionStatus == entity.ConnectionStatusActive {
		connection.FirstSyncAt = item.LastUpdatedAt
		connection.LastSyncAt = item.LastUpdatedAt
		connection.LastSuccessfulSyncAt = item.LastUpdatedAt
	} else if connectionStatus == entity.ConnectionStatusPending {
		connection.LastSyncAt = item.LastUpdatedAt
	} else {
		connection.LastErrorCode = normalizeProviderErrorCode(item.ErrorCode, "PROVIDER_STATUS_ERROR")
		connection.LastErrorMessageRedacted = "Sua conexão com o banco ainda precisa de atenção."
		connection.LastSyncAt = item.LastUpdatedAt
	}

	connection, err = service.connectionRepository.CreateOrUpdate(ctx, connection)
	if err != nil {
		return entity.Consent{}, entity.Connection{}, sharederrors.Wrap("OPEN_FINANCE_CONNECTION_ERROR", "erro ao registrar conexão", 500, err)
	}

	jobs := buildSyncJobsFromConnection(connection, now)
	if err := service.syncJobRepository.ReplaceForConnection(ctx, connection.ID, jobs); err != nil {
		return entity.Consent{}, entity.Connection{}, sharederrors.Wrap("OPEN_FINANCE_SYNC_ERROR", "erro ao preparar status inicial da conexão", 500, err)
	}

	return consent, connection, nil
}

func (service *Service) HandleCallback(ctx context.Context, state string, code string) (entity.Consent, entity.Connection, *sharederrors.AppError) {
	consent, err := service.consentRepository.GetByState(ctx, state)
	if err != nil {
		if errors.Is(err, ErrConsentNotFound) {
			return entity.Consent{}, entity.Connection{}, sharederrors.New("CONSENT_NOT_FOUND", "consentimento nao encontrado", 404)
		}

		return entity.Consent{}, entity.Connection{}, sharederrors.Wrap("OPEN_FINANCE_CALLBACK_ERROR", "erro ao localizar consentimento", 500, err)
	}

	institution, appError := service.GetInstitution(ctx, consent.InstitutionID)
	if appError != nil {
		return entity.Consent{}, entity.Connection{}, appError
	}

	now := service.now().UTC()
	codeHash := hashString(code)
	authorizationExpiresAt := now.Add(5 * time.Minute)
	if err := service.authorizationRepository.Create(ctx, entity.Authorization{
		ID:                         uuid.NewString(),
		ConsentID:                  consent.ID,
		AuthorizationCodeHash:      codeHash,
		AuthorizationCodeExpiresAt: &authorizationExpiresAt,
		PKceMethod:                 "S256",
		RedirectReceivedAt:         &now,
		CreatedAt:                  now,
	}); err != nil {
		return entity.Consent{}, entity.Connection{}, sharederrors.Wrap("OPEN_FINANCE_CALLBACK_ERROR", "erro ao registrar callback", 500, err)
	}

	tokenSet, err := service.provider.ExchangeCode(ctx, institution, consent, code)
	if err != nil {
		consent.Status = entity.ConsentStatusError
		consent.UpdatedAt = now
		_ = service.consentRepository.Update(ctx, consent)
		return entity.Consent{}, entity.Connection{}, sharederrors.Wrap("OPEN_FINANCE_TOKEN_EXCHANGE_ERROR", "erro ao trocar authorization code", 502, err)
	}

	encryptedAccessToken, err := service.cipher.Encrypt(tokenSet.AccessToken)
	if err != nil {
		return entity.Consent{}, entity.Connection{}, sharederrors.Wrap("OPEN_FINANCE_TOKEN_VAULT_ERROR", "erro ao proteger access token", 500, err)
	}

	encryptedRefreshToken, err := service.cipher.Encrypt(tokenSet.RefreshToken)
	if err != nil {
		return entity.Consent{}, entity.Connection{}, sharederrors.Wrap("OPEN_FINANCE_TOKEN_VAULT_ERROR", "erro ao proteger refresh token", 500, err)
	}

	token := entity.Token{
		ID:                    uuid.NewString(),
		ConsentID:             consent.ID,
		InstitutionID:         consent.InstitutionID,
		AccessTokenEncrypted:  encryptedAccessToken,
		RefreshTokenEncrypted: encryptedRefreshToken,
		TokenType:             tokenSet.TokenType,
		Scope:                 tokenSet.Scope,
		ExpiresAt:             tokenSet.ExpiresAt,
		RefreshExpiresAt:      tokenSet.RefreshExpiresAt,
		CreatedAt:             now,
		UpdatedAt:             now,
	}
	if err := service.tokenRepository.UpsertByConsentID(ctx, token); err != nil {
		return entity.Consent{}, entity.Connection{}, sharederrors.Wrap("OPEN_FINANCE_TOKEN_VAULT_ERROR", "erro ao persistir tokens", 500, err)
	}

	consent.Status = entity.ConsentStatusAuthorised
	consent.AuthorisedAt = &now
	consent.UpdatedAt = now
	if err := service.consentRepository.Update(ctx, consent); err != nil {
		return entity.Consent{}, entity.Connection{}, sharederrors.Wrap("OPEN_FINANCE_CALLBACK_ERROR", "erro ao atualizar consentimento", 500, err)
	}

	connection, err := service.connectionRepository.CreateOrUpdate(ctx, entity.Connection{
		ID:                   uuid.NewString(),
		UserID:               consent.UserID,
		InstitutionID:        consent.InstitutionID,
		ConsentID:            consent.ID,
		Status:               entity.ConnectionStatusActive,
		FirstSyncAt:          &now,
		LastSyncAt:           &now,
		LastSuccessfulSyncAt: &now,
		CreatedAt:            now,
		UpdatedAt:            now,
	})
	if err != nil {
		return entity.Consent{}, entity.Connection{}, sharederrors.Wrap("OPEN_FINANCE_CONNECTION_ERROR", "erro ao registrar conexao", 500, err)
	}

	jobs := buildDefaultSyncJobs(connection.ID, now)
	if err := service.syncJobRepository.ReplaceForConnection(ctx, connection.ID, jobs); err != nil {
		return entity.Consent{}, entity.Connection{}, sharederrors.Wrap("OPEN_FINANCE_SYNC_ERROR", "erro ao preparar jobs iniciais", 500, err)
	}

	return consent, connection, nil
}

func (service *Service) RevokeConsent(ctx context.Context, consentID string, userID string) *sharederrors.AppError {
	consent, appError := service.GetConsent(ctx, consentID, userID)
	if appError != nil {
		return appError
	}

	institution, appError := service.GetInstitution(ctx, consent.InstitutionID)
	if appError != nil {
		return appError
	}

	if err := service.provider.RevokeConsent(ctx, institution, consent); err != nil {
		return sharederrors.Wrap("OPEN_FINANCE_REVOKE_ERROR", "erro ao revogar consentimento externo", 502, err)
	}

	now := service.now().UTC()
	consent.Status = entity.ConsentStatusRevoked
	consent.RevokedAt = &now
	consent.UpdatedAt = now
	if err := service.consentRepository.Update(ctx, consent); err != nil {
		return sharederrors.Wrap("OPEN_FINANCE_REVOKE_ERROR", "erro ao atualizar consentimento", 500, err)
	}

	if err := service.tokenRepository.RevokeByConsentID(ctx, consent.ID); err != nil {
		return sharederrors.Wrap("OPEN_FINANCE_REVOKE_ERROR", "erro ao revogar tokens", 500, err)
	}

	connection, err := service.connectionRepository.GetByConsentID(ctx, consent.ID)
	if err == nil {
		connection.Status = entity.ConnectionStatusRevoked
		connection.UpdatedAt = now
		if updateErr := service.connectionRepository.Update(ctx, connection); updateErr != nil {
			return sharederrors.Wrap("OPEN_FINANCE_REVOKE_ERROR", "erro ao atualizar conexao", 500, updateErr)
		}
	}

	return nil
}

func (service *Service) ListConnections(ctx context.Context, userID string) ([]entity.Connection, *sharederrors.AppError) {
	connections, err := service.connectionRepository.ListByUserID(ctx, userID)
	if err != nil {
		return nil, sharederrors.Wrap("OPEN_FINANCE_CONNECTION_ERROR", "erro ao listar conexoes", 500, err)
	}

	return connections, nil
}

func (service *Service) GetConnection(ctx context.Context, connectionID string, userID string) (entity.Connection, *sharederrors.AppError) {
	connection, err := service.connectionRepository.GetByID(ctx, connectionID)
	if err != nil {
		if errors.Is(err, ErrConnectionNotFound) {
			return entity.Connection{}, sharederrors.New("CONNECTION_NOT_FOUND", "conexao nao encontrada", 404)
		}

		return entity.Connection{}, sharederrors.Wrap("OPEN_FINANCE_CONNECTION_ERROR", "erro ao buscar conexao", 500, err)
	}

	if connection.UserID != userID {
		return entity.Connection{}, sharederrors.New("CONNECTION_NOT_FOUND", "conexao nao encontrada", 404)
	}

	return connection, nil
}

func (service *Service) ListAccountSnapshots(ctx context.Context, userID string) ([]AccountSnapshot, *sharederrors.AppError) {
	sources, appError := service.listConnectedSources(ctx, userID)
	if appError != nil {
		return nil, appError
	}

	accounts, _ := service.listAccountSnapshotsForSources(ctx, sources)
	return accounts, nil
}

func (service *Service) GetOverview(ctx context.Context, userID string) (Overview, *sharederrors.AppError) {
	sources, appError := service.listConnectedSources(ctx, userID)
	if appError != nil {
		return Overview{}, appError
	}

	accounts, activeSourceIDs := service.listAccountSnapshotsForSources(ctx, sources)
	monthStart := beginningOfMonth(service.now().UTC())
	transactionQuery := ProviderTransactionQuery{
		From:     &monthStart,
		PageSize: 500,
	}
	transactions := service.listTransactionSnapshots(ctx, accounts, transactionQuery)
	sortTransactionSnapshots(transactions)

	recentTransactions := transactions
	if len(recentTransactions) > 5 {
		recentTransactions = recentTransactions[:5]
	}

	overview := Overview{
		Accounts:            accounts,
		RecentTransactions:  recentTransactions,
		ExpenseBreakdown:    buildExpenseBreakdown(transactions),
		ConnectedAccounts:   len(accounts),
		ConnectionsWithData: len(activeSourceIDs),
	}

	for _, account := range accounts {
		if isCreditAccount(account.Type) {
			overview.CreditCardBalance += account.Balance
			continue
		}

		overview.TotalAvailable += account.Balance
	}

	for _, transaction := range transactions {
		normalizedAmount := normalizeTransactionAmount(transaction)
		if normalizedAmount < 0 {
			overview.MonthExpenses += math.Abs(normalizedAmount)
			continue
		}

		overview.MonthIncome += normalizedAmount
	}

	return overview, nil
}

func (service *Service) ListTransactions(ctx context.Context, userID string, query ProviderTransactionQuery) (TransactionFeed, *sharederrors.AppError) {
	sources, appError := service.listConnectedSources(ctx, userID)
	if appError != nil {
		return TransactionFeed{}, appError
	}

	accounts, _ := service.listAccountSnapshotsForSources(ctx, sources)
	transactions := service.listTransactionSnapshots(ctx, accounts, query)
	sortTransactionSnapshots(transactions)

	now := service.now().UTC()
	currentMonthStart := beginningOfMonth(now)
	previousMonthStart := currentMonthStart.AddDate(0, -1, 0)
	currentMonthEnd := currentMonthStart.AddDate(0, 1, 0)

	feed := TransactionFeed{
		Transactions:     transactions,
		ExpenseBreakdown: buildExpenseBreakdown(filterTransactionsByMonth(transactions, currentMonthStart, currentMonthEnd)),
	}

	for _, transaction := range transactions {
		normalizedAmount := normalizeTransactionAmount(transaction)

		switch {
		case !transaction.Date.Before(currentMonthStart) && transaction.Date.Before(currentMonthEnd):
			if normalizedAmount < 0 {
				feed.MonthExpenseTotal += math.Abs(normalizedAmount)
			} else {
				feed.MonthIncomeTotal += normalizedAmount
			}
		case !transaction.Date.Before(previousMonthStart) && transaction.Date.Before(currentMonthStart) && normalizedAmount < 0:
			feed.PreviousMonthExpenseTotal += math.Abs(normalizedAmount)
		}
	}

	return feed, nil
}

func (service *Service) SyncConnection(ctx context.Context, connectionID string, userID string) ([]entity.SyncJob, *sharederrors.AppError) {
	connection, appError := service.GetConnection(ctx, connectionID, userID)
	if appError != nil {
		return nil, appError
	}

	institution, appError := service.GetInstitution(ctx, connection.InstitutionID)
	if appError != nil {
		return nil, appError
	}

	now := service.now().UTC()
	jobs, appError := service.syncConnectionResources(ctx, institution, connection, now)
	if appError != nil {
		return nil, appError
	}

	return jobs, nil
}

func (service *Service) SyncStatus(ctx context.Context, connectionID string, userID string) (entity.Connection, []entity.SyncJob, *sharederrors.AppError) {
	connection, appError := service.GetConnection(ctx, connectionID, userID)
	if appError != nil {
		return entity.Connection{}, nil, appError
	}

	jobs, err := service.syncJobRepository.ListByConnectionID(ctx, connection.ID)
	if err != nil {
		return entity.Connection{}, nil, sharederrors.Wrap("OPEN_FINANCE_SYNC_ERROR", "erro ao listar jobs de sync", 500, err)
	}

	return connection, jobs, nil
}

func (service *Service) ReconcileConnections(ctx context.Context, limit int) (ReconciliationResult, *sharederrors.AppError) {
	connections, err := service.connectionRepository.ListActive(ctx, limit)
	if err != nil {
		return ReconciliationResult{}, sharederrors.Wrap("OPEN_FINANCE_RECONCILIATION_ERROR", "erro ao listar conexoes para reconciliacao", 500, err)
	}

	result := ReconciliationResult{}
	for _, connection := range connections {
		result.Processed++

		institution, appError := service.GetInstitution(ctx, connection.InstitutionID)
		if appError != nil {
			result.Failed++
			continue
		}

		jobs, appError := service.syncConnectionResources(ctx, institution, connection, service.now().UTC())
		if appError != nil {
			result.Failed++
			continue
		}

		result.Successful++
		result.JobsCreated += len(jobs)
	}

	return result, nil
}

func (service *Service) syncConnectionResources(
	ctx context.Context,
	institution entity.Institution,
	connection entity.Connection,
	now time.Time,
) ([]entity.SyncJob, *sharederrors.AppError) {
	consent, err := service.consentRepository.GetByID(ctx, connection.ConsentID)
	if err != nil {
		return nil, sharederrors.Wrap("OPEN_FINANCE_SYNC_ERROR", "erro ao localizar consentimento da conexão", 500, err)
	}

	results, err := service.provider.SyncResources(ctx, institution, consent, connection)
	if err != nil {
		connection.Status = entity.ConnectionStatusSyncError
		connection.LastErrorCode = "SYNC_PROVIDER_ERROR"
		connection.LastErrorMessageRedacted = "falha ao sincronizar recursos"
		connection.LastSyncAt = &now
		connection.UpdatedAt = now
		_ = service.connectionRepository.Update(ctx, connection)
		return nil, sharederrors.Wrap("OPEN_FINANCE_SYNC_ERROR", "erro ao sincronizar recursos", 502, err)
	}

	jobs := make([]entity.SyncJob, 0, len(results))
	for _, result := range results {
		startedAt := now
		finishedAt := now
		jobs = append(jobs, entity.SyncJob{
			ID:                   uuid.NewString(),
			ConnectionID:         connection.ID,
			ResourceType:         result.ResourceType,
			Status:               result.Status,
			AttemptCount:         1,
			ScheduledAt:          &now,
			StartedAt:            &startedAt,
			FinishedAt:           &finishedAt,
			ErrorCode:            result.ErrorCode,
			ErrorMessageRedacted: result.ErrorMessage,
			CreatedAt:            now,
		})
	}

	if err := service.syncJobRepository.ReplaceForConnection(ctx, connection.ID, jobs); err != nil {
		return nil, sharederrors.Wrap("OPEN_FINANCE_SYNC_ERROR", "erro ao persistir status de sync", 500, err)
	}

	connection.Status = entity.ConnectionStatusActive
	connection.LastSyncAt = &now
	connection.LastSuccessfulSyncAt = &now
	connection.LastErrorCode = ""
	connection.LastErrorMessageRedacted = ""
	connection.UpdatedAt = now
	if err := service.connectionRepository.Update(ctx, connection); err != nil {
		return nil, sharederrors.Wrap("OPEN_FINANCE_SYNC_ERROR", "erro ao atualizar conexao", 500, err)
	}

	return jobs, nil
}

func buildDefaultSyncJobs(connectionID string, now time.Time) []entity.SyncJob {
	resourceTypes := []string{
		entity.ResourceAccounts,
		entity.ResourceBalances,
		entity.ResourceTransactions,
	}

	jobs := make([]entity.SyncJob, 0, len(resourceTypes))
	for _, resourceType := range resourceTypes {
		scheduledAt := now
		jobs = append(jobs, entity.SyncJob{
			ID:           uuid.NewString(),
			ConnectionID: connectionID,
			ResourceType: resourceType,
			Status:       entity.SyncJobStatusPending,
			AttemptCount: 0,
			ScheduledAt:  &scheduledAt,
			CreatedAt:    now,
		})
	}

	return jobs
}

type connectedSource struct {
	connection  entity.Connection
	consent     entity.Consent
	institution entity.Institution
}

func (service *Service) listConnectedSources(ctx context.Context, userID string) ([]connectedSource, *sharederrors.AppError) {
	connections, err := service.connectionRepository.ListByUserID(ctx, userID)
	if err != nil {
		return nil, sharederrors.Wrap("OPEN_FINANCE_CONNECTION_ERROR", "erro ao listar conexoes", 500, err)
	}

	institutionCache := make(map[string]entity.Institution, len(connections))
	sources := make([]connectedSource, 0, len(connections))

	for _, connection := range connections {
		if connection.Status == entity.ConnectionStatusRevoked {
			continue
		}

		consent, err := service.consentRepository.GetByID(ctx, connection.ConsentID)
		if err != nil || strings.TrimSpace(consent.ExternalConsentID) == "" {
			continue
		}

		institution, exists := institutionCache[connection.InstitutionID]
		if !exists {
			institution, err = service.institutionRepository.GetByID(ctx, connection.InstitutionID)
			if err != nil {
				continue
			}

			institutionCache[connection.InstitutionID] = institution
		}

		sources = append(sources, connectedSource{
			connection:  connection,
			consent:     consent,
			institution: institution,
		})
	}

	return sources, nil
}

func (service *Service) listAccountSnapshotsForSources(ctx context.Context, sources []connectedSource) ([]AccountSnapshot, map[string]struct{}) {
	accounts := make([]AccountSnapshot, 0)
	activeSourceIDs := make(map[string]struct{}, len(sources))

	for _, source := range sources {
		providerAccounts, err := service.provider.ListAccounts(ctx, source.consent.ExternalConsentID)
		if err != nil {
			continue
		}

		for _, account := range providerAccounts {
			accounts = append(accounts, AccountSnapshot{
				ID:                   account.ID,
				ConnectionID:         source.connection.ID,
				InstitutionID:        source.institution.ID,
				InstitutionName:      source.institution.DisplayName,
				ItemID:               account.ItemID,
				Type:                 account.Type,
				Subtype:              account.Subtype,
				Name:                 account.Name,
				MarketingName:        account.MarketingName,
				Number:               account.Number,
				CurrencyCode:         account.CurrencyCode,
				Balance:              account.Balance,
				BankTransferNumber:   account.BankTransferNumber,
				CreditBrand:          account.CreditBrand,
				AvailableCreditLimit: account.AvailableCreditLimit,
			})
			activeSourceIDs[source.connection.ID] = struct{}{}
		}
	}

	return accounts, activeSourceIDs
}

func (service *Service) listTransactionSnapshots(ctx context.Context, accounts []AccountSnapshot, query ProviderTransactionQuery) []TransactionSnapshot {
	transactions := make([]TransactionSnapshot, 0)

	for _, account := range accounts {
		providerTransactions, err := service.provider.ListTransactions(ctx, account.ID, query)
		if err != nil {
			continue
		}

		for _, transaction := range providerTransactions {
			transactions = append(transactions, TransactionSnapshot{
				ID:              transaction.ID,
				AccountID:       transaction.AccountID,
				ConnectionID:    account.ConnectionID,
				InstitutionName: account.InstitutionName,
				AccountName:     resolveAccountName(account),
				Description:     strings.TrimSpace(transaction.Description),
				Category:        resolveTransactionCategory(transaction),
				Type:            transaction.Type,
				Status:          transaction.Status,
				CurrencyCode:    transaction.CurrencyCode,
				Amount:          transaction.Amount,
				Date:            transaction.Date,
				MerchantName:    strings.TrimSpace(transaction.MerchantName),
			})
		}
	}

	return transactions
}

func resolveAccountName(account AccountSnapshot) string {
	if value := strings.TrimSpace(account.MarketingName); value != "" {
		return value
	}
	if value := strings.TrimSpace(account.Name); value != "" {
		return value
	}
	if value := strings.TrimSpace(account.Number); value != "" {
		return value
	}
	return "Conta conectada"
}

func resolveTransactionCategory(transaction ProviderTransaction) string {
	if value := strings.TrimSpace(transaction.Category); value != "" {
		return value
	}
	if value := strings.TrimSpace(transaction.MerchantName); value != "" {
		return value
	}
	return "Sem categoria"
}

func sortTransactionSnapshots(transactions []TransactionSnapshot) {
	sort.SliceStable(transactions, func(left int, right int) bool {
		return transactions[left].Date.After(transactions[right].Date)
	})
}

func buildExpenseBreakdown(transactions []TransactionSnapshot) []CategoryBreakdown {
	totalsByCategory := make(map[string]float64)
	totalExpenses := 0.0

	for _, transaction := range transactions {
		normalizedAmount := normalizeTransactionAmount(transaction)
		if normalizedAmount >= 0 {
			continue
		}

		amount := math.Abs(normalizedAmount)
		totalsByCategory[transaction.Category] += amount
		totalExpenses += amount
	}

	if totalExpenses == 0 {
		return []CategoryBreakdown{}
	}

	breakdown := make([]CategoryBreakdown, 0, len(totalsByCategory))
	for category, amount := range totalsByCategory {
		breakdown = append(breakdown, CategoryBreakdown{
			Category: category,
			Amount:   amount,
			Percent:  (amount / totalExpenses) * 100,
		})
	}

	sort.SliceStable(breakdown, func(left int, right int) bool {
		return breakdown[left].Amount > breakdown[right].Amount
	})

	if len(breakdown) > 4 {
		breakdown = breakdown[:4]
	}

	return breakdown
}

func filterTransactionsByMonth(transactions []TransactionSnapshot, start time.Time, end time.Time) []TransactionSnapshot {
	filtered := make([]TransactionSnapshot, 0, len(transactions))
	for _, transaction := range transactions {
		if transaction.Date.Before(start) || !transaction.Date.Before(end) {
			continue
		}
		filtered = append(filtered, transaction)
	}
	return filtered
}

func beginningOfMonth(value time.Time) time.Time {
	return time.Date(value.Year(), value.Month(), 1, 0, 0, 0, 0, time.UTC)
}

func isCreditAccount(accountType string) bool {
	return strings.EqualFold(strings.TrimSpace(accountType), "CREDIT")
}

func normalizeTransactionAmount(transaction TransactionSnapshot) float64 {
	amount := math.Abs(transaction.Amount)
	if strings.EqualFold(transaction.Type, "DEBIT") {
		return -amount
	}
	return amount
}

func hasSandboxInstitution(institutions []entity.Institution) bool {
	for _, institution := range institutions {
		name := strings.ToLower(strings.TrimSpace(institution.DisplayName))
		if institution.Status == "sandbox" || strings.Contains(name, "pluggy") {
			return true
		}
	}

	return false
}

func buildSyncJobsFromConnection(connection entity.Connection, now time.Time) []entity.SyncJob {
	resourceTypes := []string{
		entity.ResourceAccounts,
		entity.ResourceBalances,
		entity.ResourceTransactions,
	}

	jobs := make([]entity.SyncJob, 0, len(resourceTypes))
	for _, resourceType := range resourceTypes {
		job := entity.SyncJob{
			ID:           uuid.NewString(),
			ConnectionID: connection.ID,
			ResourceType: resourceType,
			AttemptCount: 1,
			ScheduledAt:  &now,
			CreatedAt:    now,
		}

		switch connection.Status {
		case entity.ConnectionStatusActive:
			job.Status = entity.SyncJobStatusCompleted
			job.StartedAt = connection.LastSyncAt
			job.FinishedAt = connection.LastSuccessfulSyncAt
		case entity.ConnectionStatusPending:
			job.Status = entity.SyncJobStatusPending
		default:
			job.Status = entity.SyncJobStatusError
			job.StartedAt = connection.LastSyncAt
			job.FinishedAt = connection.LastSyncAt
			job.ErrorCode = normalizeProviderErrorCode(connection.LastErrorCode, "PROVIDER_STATUS_ERROR")
			job.ErrorMessageRedacted = connection.LastErrorMessageRedacted
		}

		jobs = append(jobs, job)
	}

	return jobs
}

func validateProviderItem(institution entity.Institution, item ProviderItem) error {
	expectedConnectorID := institution.DirectoryOrgID
	if expectedConnectorID == "" {
		return errors.New("connector id missing for institution")
	}

	if expectedConnectorID != stringifyConnectorID(item.ConnectorID) {
		return errors.New("item connector mismatch")
	}

	return nil
}

func mapConnectionStatus(item ProviderItem) string {
	switch item.Status {
	case "UPDATED":
		return entity.ConnectionStatusActive
	case "WAITING_USER_INPUT", "WAITING_USER_ACTION", "UPDATING", "CREATED":
		return entity.ConnectionStatusPending
	case "LOGIN_ERROR":
		return entity.ConnectionStatusReauthRequired
	default:
		if item.ExecutionStatus != "" {
			switch item.ExecutionStatus {
			case "WAITING_USER_INPUT", "LOGIN_IN_PROGRESS", "CREATING_IN_PROGRESS":
				return entity.ConnectionStatusPending
			case "LOGIN_ERROR", "USER_INPUT_TIMEOUT":
				return entity.ConnectionStatusReauthRequired
			}
		}

		return entity.ConnectionStatusSyncError
	}
}

func normalizeProviderErrorCode(value string, fallback string) string {
	if value == "" {
		return fallback
	}

	return value
}

func stringifyConnectorID(connectorID int64) string {
	return strconv.FormatInt(connectorID, 10)
}

func hashString(value string) string {
	sum := sha256.Sum256([]byte(value))
	return hex.EncodeToString(sum[:])
}
