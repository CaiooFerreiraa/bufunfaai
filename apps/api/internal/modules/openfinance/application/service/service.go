package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
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
	ExchangeCode(ctx context.Context, institution entity.Institution, consent entity.Consent, code string) (ProviderTokenSet, error)
	RevokeConsent(ctx context.Context, institution entity.Institution, consent entity.Consent) error
	SyncResources(ctx context.Context, institution entity.Institution, connection entity.Connection) ([]ProviderSyncResult, error)
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

type ProviderSyncResult struct {
	ResourceType string
	Status       string
	ErrorCode    string
	ErrorMessage string
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
	if err == nil && len(institutions) > 0 {
		return nil
	}

	discoveredInstitutions, err := service.provider.ListInstitutions(ctx)
	if err != nil {
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
	results, err := service.provider.SyncResources(ctx, institution, connection)
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

func hashString(value string) string {
	sum := sha256.Sum256([]byte(value))
	return hex.EncodeToString(sum[:])
}
