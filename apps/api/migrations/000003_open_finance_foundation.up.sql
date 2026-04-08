CREATE TABLE IF NOT EXISTS of_institutions (
    id UUID PRIMARY KEY,
    directory_org_id VARCHAR(120) NOT NULL,
    brand_name VARCHAR(160) NOT NULL,
    display_name VARCHAR(160) NOT NULL,
    authorisation_server_url TEXT NOT NULL,
    resources_base_url TEXT NOT NULL,
    logo_url TEXT NULL,
    status VARCHAR(30) NOT NULL,
    supports_data_sharing BOOLEAN NOT NULL DEFAULT TRUE,
    supports_payments BOOLEAN NOT NULL DEFAULT FALSE,
    last_directory_sync_at TIMESTAMPTZ NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS of_consents (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    institution_id UUID NOT NULL REFERENCES of_institutions(id) ON DELETE RESTRICT,
    external_consent_id VARCHAR(255) NOT NULL,
    status VARCHAR(30) NOT NULL,
    purpose VARCHAR(255) NOT NULL,
    permissions_json JSONB NOT NULL,
    expiration_at TIMESTAMPTZ NULL,
    revoked_at TIMESTAMPTZ NULL,
    authorised_at TIMESTAMPTZ NULL,
    denied_at TIMESTAMPTZ NULL,
    redirect_uri TEXT NOT NULL,
    state VARCHAR(255) NOT NULL UNIQUE,
    nonce VARCHAR(255) NOT NULL,
    code_verifier_encrypted TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_of_consents_user_id ON of_consents (user_id);
CREATE INDEX IF NOT EXISTS idx_of_consents_institution_id ON of_consents (institution_id);
CREATE INDEX IF NOT EXISTS idx_of_consents_status ON of_consents (status);

CREATE TABLE IF NOT EXISTS of_authorizations (
    id UUID PRIMARY KEY,
    consent_id UUID NOT NULL REFERENCES of_consents(id) ON DELETE CASCADE,
    authorization_code_hash CHAR(64) NOT NULL,
    authorization_code_expires_at TIMESTAMPTZ NULL,
    pkce_method VARCHAR(20) NOT NULL,
    redirect_received_at TIMESTAMPTZ NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_of_authorizations_consent_id ON of_authorizations (consent_id);

CREATE TABLE IF NOT EXISTS of_tokens (
    id UUID PRIMARY KEY,
    consent_id UUID NOT NULL UNIQUE REFERENCES of_consents(id) ON DELETE CASCADE,
    institution_id UUID NOT NULL REFERENCES of_institutions(id) ON DELETE RESTRICT,
    access_token_encrypted TEXT NOT NULL,
    refresh_token_encrypted TEXT NOT NULL,
    token_type VARCHAR(40) NOT NULL,
    scope TEXT NOT NULL,
    expires_at TIMESTAMPTZ NOT NULL,
    refresh_expires_at TIMESTAMPTZ NULL,
    last_refresh_at TIMESTAMPTZ NULL,
    revoked_at TIMESTAMPTZ NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS of_connections (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    institution_id UUID NOT NULL REFERENCES of_institutions(id) ON DELETE RESTRICT,
    consent_id UUID NOT NULL UNIQUE REFERENCES of_consents(id) ON DELETE CASCADE,
    status VARCHAR(40) NOT NULL,
    first_sync_at TIMESTAMPTZ NULL,
    last_sync_at TIMESTAMPTZ NULL,
    last_successful_sync_at TIMESTAMPTZ NULL,
    last_error_code VARCHAR(100) NULL,
    last_error_message_redacted TEXT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_of_connections_user_id ON of_connections (user_id);
CREATE INDEX IF NOT EXISTS idx_of_connections_status ON of_connections (status);

CREATE TABLE IF NOT EXISTS of_sync_jobs (
    id UUID PRIMARY KEY,
    connection_id UUID NOT NULL REFERENCES of_connections(id) ON DELETE CASCADE,
    resource_type VARCHAR(60) NOT NULL,
    status VARCHAR(30) NOT NULL,
    cursor TEXT NULL,
    window_start TIMESTAMPTZ NULL,
    window_end TIMESTAMPTZ NULL,
    attempt_count INT NOT NULL DEFAULT 0,
    scheduled_at TIMESTAMPTZ NULL,
    started_at TIMESTAMPTZ NULL,
    finished_at TIMESTAMPTZ NULL,
    error_code VARCHAR(100) NULL,
    error_message_redacted TEXT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_of_sync_jobs_connection_id ON of_sync_jobs (connection_id);

CREATE TABLE IF NOT EXISTS of_sync_checkpoints (
    id UUID PRIMARY KEY,
    connection_id UUID NOT NULL REFERENCES of_connections(id) ON DELETE CASCADE,
    resource_type VARCHAR(60) NOT NULL,
    cursor TEXT NULL,
    last_reference_datetime TIMESTAMPTZ NULL,
    etag TEXT NULL,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_of_sync_checkpoints_connection_resource ON of_sync_checkpoints (connection_id, resource_type);

CREATE TABLE IF NOT EXISTS of_webhook_events (
    id UUID PRIMARY KEY,
    institution_id UUID NULL REFERENCES of_institutions(id) ON DELETE SET NULL,
    event_type VARCHAR(100) NOT NULL,
    external_event_id VARCHAR(255) NULL,
    signature_valid BOOLEAN NOT NULL DEFAULT FALSE,
    payload_encrypted TEXT NOT NULL,
    received_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    processed_at TIMESTAMPTZ NULL,
    status VARCHAR(40) NOT NULL
);
