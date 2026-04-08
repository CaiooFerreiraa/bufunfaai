package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	ofservice "github.com/bufunfaai/bufunfaai/apps/api/internal/modules/openfinance/application/service"
	"github.com/bufunfaai/bufunfaai/apps/api/internal/modules/openfinance/domain/entity"
)

type PostgresInstitutionRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresInstitutionRepository(pool *pgxpool.Pool) *PostgresInstitutionRepository {
	return &PostgresInstitutionRepository{pool: pool}
}

func (repository *PostgresInstitutionRepository) SaveMany(ctx context.Context, institutions []entity.Institution) error {
	for _, institution := range institutions {
		_, err := repository.pool.Exec(ctx, `
			INSERT INTO of_institutions (
				id, directory_org_id, brand_name, display_name, authorisation_server_url, resources_base_url, logo_url,
				status, supports_data_sharing, supports_payments, last_directory_sync_at, created_at, updated_at
			) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13)
			ON CONFLICT (id) DO UPDATE SET
				directory_org_id = EXCLUDED.directory_org_id,
				brand_name = EXCLUDED.brand_name,
				display_name = EXCLUDED.display_name,
				authorisation_server_url = EXCLUDED.authorisation_server_url,
				resources_base_url = EXCLUDED.resources_base_url,
				logo_url = EXCLUDED.logo_url,
				status = EXCLUDED.status,
				supports_data_sharing = EXCLUDED.supports_data_sharing,
				supports_payments = EXCLUDED.supports_payments,
				last_directory_sync_at = EXCLUDED.last_directory_sync_at,
				updated_at = EXCLUDED.updated_at
		`,
			institution.ID,
			institution.DirectoryOrgID,
			institution.BrandName,
			institution.DisplayName,
			institution.AuthorisationServerURL,
			institution.ResourcesBaseURL,
			institution.LogoURL,
			institution.Status,
			institution.SupportsDataSharing,
			institution.SupportsPayments,
			institution.LastDirectorySyncAt,
			institution.CreatedAt,
			institution.UpdatedAt,
		)
		if err != nil {
			return err
		}
	}

	return nil
}

func (repository *PostgresInstitutionRepository) List(ctx context.Context) ([]entity.Institution, error) {
	rows, err := repository.pool.Query(ctx, `
		SELECT id, directory_org_id, brand_name, display_name, authorisation_server_url, resources_base_url, logo_url,
		       status, supports_data_sharing, supports_payments, last_directory_sync_at, created_at, updated_at
		FROM of_institutions
		ORDER BY display_name
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	institutions := make([]entity.Institution, 0)
	for rows.Next() {
		institution, scanErr := scanInstitution(rows)
		if scanErr != nil {
			return nil, scanErr
		}

		institutions = append(institutions, institution)
	}

	return institutions, rows.Err()
}

func (repository *PostgresInstitutionRepository) GetByID(ctx context.Context, institutionID string) (entity.Institution, error) {
	row := repository.pool.QueryRow(ctx, `
		SELECT id, directory_org_id, brand_name, display_name, authorisation_server_url, resources_base_url, logo_url,
		       status, supports_data_sharing, supports_payments, last_directory_sync_at, created_at, updated_at
		FROM of_institutions
		WHERE id = $1
	`, institutionID)

	return scanInstitution(row)
}

type PostgresConsentRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresConsentRepository(pool *pgxpool.Pool) *PostgresConsentRepository {
	return &PostgresConsentRepository{pool: pool}
}

func (repository *PostgresConsentRepository) Create(ctx context.Context, consent entity.Consent) error {
	_, err := repository.pool.Exec(ctx, `
		INSERT INTO of_consents (
			id, user_id, institution_id, external_consent_id, status, purpose, permissions_json, expiration_at,
			revoked_at, authorised_at, denied_at, redirect_uri, state, nonce, code_verifier_encrypted, created_at, updated_at
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17)
	`,
		consent.ID,
		consent.UserID,
		consent.InstitutionID,
		consent.ExternalConsentID,
		consent.Status,
		consent.Purpose,
		consent.PermissionsJSON,
		consent.ExpirationAt,
		consent.RevokedAt,
		consent.AuthorisedAt,
		consent.DeniedAt,
		consent.RedirectURI,
		consent.State,
		consent.Nonce,
		consent.CodeVerifierEncrypted,
		consent.CreatedAt,
		consent.UpdatedAt,
	)
	return err
}

func (repository *PostgresConsentRepository) GetByID(ctx context.Context, consentID string) (entity.Consent, error) {
	row := repository.pool.QueryRow(ctx, `
		SELECT id, user_id, institution_id, external_consent_id, status, purpose, permissions_json, expiration_at,
		       revoked_at, authorised_at, denied_at, redirect_uri, state, nonce, code_verifier_encrypted, created_at, updated_at
		FROM of_consents WHERE id = $1
	`, consentID)
	return scanConsent(row)
}

func (repository *PostgresConsentRepository) GetByState(ctx context.Context, state string) (entity.Consent, error) {
	row := repository.pool.QueryRow(ctx, `
		SELECT id, user_id, institution_id, external_consent_id, status, purpose, permissions_json, expiration_at,
		       revoked_at, authorised_at, denied_at, redirect_uri, state, nonce, code_verifier_encrypted, created_at, updated_at
		FROM of_consents WHERE state = $1
	`, state)
	return scanConsent(row)
}

func (repository *PostgresConsentRepository) Update(ctx context.Context, consent entity.Consent) error {
	_, err := repository.pool.Exec(ctx, `
		UPDATE of_consents SET
			external_consent_id = $2,
			status = $3,
			purpose = $4,
			permissions_json = $5,
			expiration_at = $6,
			revoked_at = $7,
			authorised_at = $8,
			denied_at = $9,
			redirect_uri = $10,
			state = $11,
			nonce = $12,
			code_verifier_encrypted = $13,
			updated_at = $14
		WHERE id = $1
	`,
		consent.ID,
		consent.ExternalConsentID,
		consent.Status,
		consent.Purpose,
		consent.PermissionsJSON,
		consent.ExpirationAt,
		consent.RevokedAt,
		consent.AuthorisedAt,
		consent.DeniedAt,
		consent.RedirectURI,
		consent.State,
		consent.Nonce,
		consent.CodeVerifierEncrypted,
		consent.UpdatedAt,
	)
	return err
}

type PostgresAuthorizationRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresAuthorizationRepository(pool *pgxpool.Pool) *PostgresAuthorizationRepository {
	return &PostgresAuthorizationRepository{pool: pool}
}

func (repository *PostgresAuthorizationRepository) Create(ctx context.Context, authorization entity.Authorization) error {
	_, err := repository.pool.Exec(ctx, `
		INSERT INTO of_authorizations (
			id, consent_id, authorization_code_hash, authorization_code_expires_at, pkce_method, redirect_received_at, created_at
		) VALUES ($1,$2,$3,$4,$5,$6,$7)
	`, authorization.ID, authorization.ConsentID, authorization.AuthorizationCodeHash, authorization.AuthorizationCodeExpiresAt, authorization.PKceMethod, authorization.RedirectReceivedAt, authorization.CreatedAt)
	return err
}

type PostgresTokenRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresTokenRepository(pool *pgxpool.Pool) *PostgresTokenRepository {
	return &PostgresTokenRepository{pool: pool}
}

func (repository *PostgresTokenRepository) UpsertByConsentID(ctx context.Context, token entity.Token) error {
	_, err := repository.pool.Exec(ctx, `
		INSERT INTO of_tokens (
			id, consent_id, institution_id, access_token_encrypted, refresh_token_encrypted, token_type, scope, expires_at,
			refresh_expires_at, last_refresh_at, revoked_at, created_at, updated_at
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13)
		ON CONFLICT (consent_id) DO UPDATE SET
			institution_id = EXCLUDED.institution_id,
			access_token_encrypted = EXCLUDED.access_token_encrypted,
			refresh_token_encrypted = EXCLUDED.refresh_token_encrypted,
			token_type = EXCLUDED.token_type,
			scope = EXCLUDED.scope,
			expires_at = EXCLUDED.expires_at,
			refresh_expires_at = EXCLUDED.refresh_expires_at,
			last_refresh_at = EXCLUDED.last_refresh_at,
			revoked_at = EXCLUDED.revoked_at,
			updated_at = EXCLUDED.updated_at
	`, token.ID, token.ConsentID, token.InstitutionID, token.AccessTokenEncrypted, token.RefreshTokenEncrypted, token.TokenType, token.Scope, token.ExpiresAt, token.RefreshExpiresAt, token.LastRefreshAt, token.RevokedAt, token.CreatedAt, token.UpdatedAt)
	return err
}

func (repository *PostgresTokenRepository) GetByConsentID(ctx context.Context, consentID string) (entity.Token, error) {
	row := repository.pool.QueryRow(ctx, `
		SELECT id, consent_id, institution_id, access_token_encrypted, refresh_token_encrypted, token_type, scope, expires_at,
		       refresh_expires_at, last_refresh_at, revoked_at, created_at, updated_at
		FROM of_tokens WHERE consent_id = $1
	`, consentID)
	return scanToken(row)
}

func (repository *PostgresTokenRepository) RevokeByConsentID(ctx context.Context, consentID string) error {
	_, err := repository.pool.Exec(ctx, `UPDATE of_tokens SET revoked_at = NOW(), updated_at = NOW() WHERE consent_id = $1`, consentID)
	return err
}

type PostgresConnectionRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresConnectionRepository(pool *pgxpool.Pool) *PostgresConnectionRepository {
	return &PostgresConnectionRepository{pool: pool}
}

func (repository *PostgresConnectionRepository) CreateOrUpdate(ctx context.Context, connection entity.Connection) (entity.Connection, error) {
	row := repository.pool.QueryRow(ctx, `
		INSERT INTO of_connections (
			id, user_id, institution_id, consent_id, status, first_sync_at, last_sync_at, last_successful_sync_at,
			last_error_code, last_error_message_redacted, created_at, updated_at
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)
		ON CONFLICT (consent_id) DO UPDATE SET
			status = EXCLUDED.status,
			first_sync_at = COALESCE(of_connections.first_sync_at, EXCLUDED.first_sync_at),
			last_sync_at = EXCLUDED.last_sync_at,
			last_successful_sync_at = EXCLUDED.last_successful_sync_at,
			last_error_code = EXCLUDED.last_error_code,
			last_error_message_redacted = EXCLUDED.last_error_message_redacted,
			updated_at = EXCLUDED.updated_at
		RETURNING id, user_id, institution_id, consent_id, status, first_sync_at, last_sync_at, last_successful_sync_at,
		          last_error_code, last_error_message_redacted, created_at, updated_at
	`,
		connection.ID, connection.UserID, connection.InstitutionID, connection.ConsentID, connection.Status, connection.FirstSyncAt, connection.LastSyncAt,
		connection.LastSuccessfulSyncAt, connection.LastErrorCode, connection.LastErrorMessageRedacted, connection.CreatedAt, connection.UpdatedAt,
	)
	return scanConnection(row)
}

func (repository *PostgresConnectionRepository) ListByUserID(ctx context.Context, userID string) ([]entity.Connection, error) {
	rows, err := repository.pool.Query(ctx, `
		SELECT id, user_id, institution_id, consent_id, status, first_sync_at, last_sync_at, last_successful_sync_at,
		       last_error_code, last_error_message_redacted, created_at, updated_at
		FROM of_connections WHERE user_id = $1 ORDER BY created_at DESC
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	connections := make([]entity.Connection, 0)
	for rows.Next() {
		connection, scanErr := scanConnection(rows)
		if scanErr != nil {
			return nil, scanErr
		}

		connections = append(connections, connection)
	}

	return connections, rows.Err()
}

func (repository *PostgresConnectionRepository) ListActive(ctx context.Context, limit int) ([]entity.Connection, error) {
	rows, err := repository.pool.Query(ctx, `
		SELECT id, user_id, institution_id, consent_id, status, first_sync_at, last_sync_at, last_successful_sync_at,
		       last_error_code, last_error_message_redacted, created_at, updated_at
		FROM of_connections
		WHERE status IN ($1, $2)
		ORDER BY COALESCE(last_sync_at, created_at) ASC
		LIMIT $3
	`, entity.ConnectionStatusActive, entity.ConnectionStatusSyncError, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	connections := make([]entity.Connection, 0)
	for rows.Next() {
		connection, scanErr := scanConnection(rows)
		if scanErr != nil {
			return nil, scanErr
		}

		connections = append(connections, connection)
	}

	return connections, rows.Err()
}

func (repository *PostgresConnectionRepository) GetByID(ctx context.Context, connectionID string) (entity.Connection, error) {
	row := repository.pool.QueryRow(ctx, `
		SELECT id, user_id, institution_id, consent_id, status, first_sync_at, last_sync_at, last_successful_sync_at,
		       last_error_code, last_error_message_redacted, created_at, updated_at
		FROM of_connections WHERE id = $1
	`, connectionID)
	return scanConnection(row)
}

func (repository *PostgresConnectionRepository) GetByConsentID(ctx context.Context, consentID string) (entity.Connection, error) {
	row := repository.pool.QueryRow(ctx, `
		SELECT id, user_id, institution_id, consent_id, status, first_sync_at, last_sync_at, last_successful_sync_at,
		       last_error_code, last_error_message_redacted, created_at, updated_at
		FROM of_connections WHERE consent_id = $1
	`, consentID)
	return scanConnection(row)
}

func (repository *PostgresConnectionRepository) Update(ctx context.Context, connection entity.Connection) error {
	_, err := repository.pool.Exec(ctx, `
		UPDATE of_connections SET
			status = $2,
			first_sync_at = $3,
			last_sync_at = $4,
			last_successful_sync_at = $5,
			last_error_code = $6,
			last_error_message_redacted = $7,
			updated_at = $8
		WHERE id = $1
	`, connection.ID, connection.Status, connection.FirstSyncAt, connection.LastSyncAt, connection.LastSuccessfulSyncAt, connection.LastErrorCode, connection.LastErrorMessageRedacted, connection.UpdatedAt)
	return err
}

type PostgresSyncJobRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresSyncJobRepository(pool *pgxpool.Pool) *PostgresSyncJobRepository {
	return &PostgresSyncJobRepository{pool: pool}
}

func (repository *PostgresSyncJobRepository) ReplaceForConnection(ctx context.Context, connectionID string, jobs []entity.SyncJob) error {
	tx, err := repository.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	if _, err := tx.Exec(ctx, `DELETE FROM of_sync_jobs WHERE connection_id = $1`, connectionID); err != nil {
		return err
	}

	for _, job := range jobs {
		_, err := tx.Exec(ctx, `
			INSERT INTO of_sync_jobs (
				id, connection_id, resource_type, status, cursor, window_start, window_end, attempt_count,
				scheduled_at, started_at, finished_at, error_code, error_message_redacted, created_at
			) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14)
		`, job.ID, job.ConnectionID, job.ResourceType, job.Status, job.Cursor, job.WindowStart, job.WindowEnd, job.AttemptCount, job.ScheduledAt, job.StartedAt, job.FinishedAt, job.ErrorCode, job.ErrorMessageRedacted, job.CreatedAt)
		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

func (repository *PostgresSyncJobRepository) ListByConnectionID(ctx context.Context, connectionID string) ([]entity.SyncJob, error) {
	rows, err := repository.pool.Query(ctx, `
		SELECT id, connection_id, resource_type, status, cursor, window_start, window_end, attempt_count,
		       scheduled_at, started_at, finished_at, error_code, error_message_redacted, created_at
		FROM of_sync_jobs WHERE connection_id = $1 ORDER BY created_at DESC
	`, connectionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	jobs := make([]entity.SyncJob, 0)
	for rows.Next() {
		job, scanErr := scanSyncJob(rows)
		if scanErr != nil {
			return nil, scanErr
		}

		jobs = append(jobs, job)
	}

	return jobs, rows.Err()
}

type scanner interface {
	Scan(dest ...any) error
}

func scanInstitution(row scanner) (entity.Institution, error) {
	var institution entity.Institution
	err := row.Scan(&institution.ID, &institution.DirectoryOrgID, &institution.BrandName, &institution.DisplayName, &institution.AuthorisationServerURL, &institution.ResourcesBaseURL, &institution.LogoURL, &institution.Status, &institution.SupportsDataSharing, &institution.SupportsPayments, &institution.LastDirectorySyncAt, &institution.CreatedAt, &institution.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.Institution{}, ofservice.ErrInstitutionNotFound
		}
		return entity.Institution{}, err
	}
	return institution, nil
}

func scanConsent(row scanner) (entity.Consent, error) {
	var consent entity.Consent
	err := row.Scan(&consent.ID, &consent.UserID, &consent.InstitutionID, &consent.ExternalConsentID, &consent.Status, &consent.Purpose, &consent.PermissionsJSON, &consent.ExpirationAt, &consent.RevokedAt, &consent.AuthorisedAt, &consent.DeniedAt, &consent.RedirectURI, &consent.State, &consent.Nonce, &consent.CodeVerifierEncrypted, &consent.CreatedAt, &consent.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.Consent{}, ofservice.ErrConsentNotFound
		}
		return entity.Consent{}, err
	}
	return consent, nil
}

func scanToken(row scanner) (entity.Token, error) {
	var token entity.Token
	err := row.Scan(&token.ID, &token.ConsentID, &token.InstitutionID, &token.AccessTokenEncrypted, &token.RefreshTokenEncrypted, &token.TokenType, &token.Scope, &token.ExpiresAt, &token.RefreshExpiresAt, &token.LastRefreshAt, &token.RevokedAt, &token.CreatedAt, &token.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.Token{}, ofservice.ErrTokenNotFound
		}
		return entity.Token{}, err
	}
	return token, nil
}

func scanConnection(row scanner) (entity.Connection, error) {
	var connection entity.Connection
	err := row.Scan(&connection.ID, &connection.UserID, &connection.InstitutionID, &connection.ConsentID, &connection.Status, &connection.FirstSyncAt, &connection.LastSyncAt, &connection.LastSuccessfulSyncAt, &connection.LastErrorCode, &connection.LastErrorMessageRedacted, &connection.CreatedAt, &connection.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.Connection{}, ofservice.ErrConnectionNotFound
		}
		return entity.Connection{}, err
	}
	return connection, nil
}

func scanSyncJob(row scanner) (entity.SyncJob, error) {
	var job entity.SyncJob
	err := row.Scan(&job.ID, &job.ConnectionID, &job.ResourceType, &job.Status, &job.Cursor, &job.WindowStart, &job.WindowEnd, &job.AttemptCount, &job.ScheduledAt, &job.StartedAt, &job.FinishedAt, &job.ErrorCode, &job.ErrorMessageRedacted, &job.CreatedAt)
	if err != nil {
		return entity.SyncJob{}, err
	}
	return job, nil
}
