package repository

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	authservice "github.com/bufunfaai/bufunfaai/apps/api/internal/modules/auth/application/service"
	authentity "github.com/bufunfaai/bufunfaai/apps/api/internal/modules/auth/domain/entity"
)

type PostgresRefreshTokenRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresRefreshTokenRepository(pool *pgxpool.Pool) *PostgresRefreshTokenRepository {
	return &PostgresRefreshTokenRepository{pool: pool}
}

func (repository *PostgresRefreshTokenRepository) Create(ctx context.Context, token authentity.RefreshToken) error {
	_, err := repository.pool.Exec(
		ctx,
		`
			INSERT INTO refresh_tokens (
				id, user_id, token_hash, device_id, expires_at, revoked_at, created_at, last_used_at, ip_address, user_agent
			) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		`,
		token.ID,
		token.UserID,
		token.TokenHash,
		token.DeviceID,
		token.ExpiresAt,
		token.RevokedAt,
		token.CreatedAt,
		token.LastUsedAt,
		token.IPAddress,
		token.UserAgent,
	)
	return err
}

func (repository *PostgresRefreshTokenRepository) GetByTokenHash(ctx context.Context, tokenHash string) (authentity.RefreshToken, error) {
	row := repository.pool.QueryRow(
		ctx,
		`
			SELECT id, user_id, token_hash, device_id, expires_at, revoked_at, created_at, last_used_at, ip_address, user_agent
			FROM refresh_tokens
			WHERE token_hash = $1
		`,
		tokenHash,
	)

	var token authentity.RefreshToken
	err := row.Scan(
		&token.ID,
		&token.UserID,
		&token.TokenHash,
		&token.DeviceID,
		&token.ExpiresAt,
		&token.RevokedAt,
		&token.CreatedAt,
		&token.LastUsedAt,
		&token.IPAddress,
		&token.UserAgent,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return authentity.RefreshToken{}, authservice.ErrRefreshTokenNotFound
		}

		return authentity.RefreshToken{}, err
	}

	return token, nil
}

func (repository *PostgresRefreshTokenRepository) Revoke(ctx context.Context, tokenID string, revokedAt time.Time) error {
	_, err := repository.pool.Exec(ctx, `UPDATE refresh_tokens SET revoked_at = $2 WHERE id = $1`, tokenID, revokedAt)
	return err
}

func (repository *PostgresRefreshTokenRepository) RevokeByUserID(ctx context.Context, userID string, revokedAt time.Time) error {
	_, err := repository.pool.Exec(ctx, `UPDATE refresh_tokens SET revoked_at = $2 WHERE user_id = $1 AND revoked_at IS NULL`, userID, revokedAt)
	return err
}

func (repository *PostgresRefreshTokenRepository) RevokeByUserIDAndDeviceID(ctx context.Context, userID string, deviceID string, revokedAt time.Time) error {
	_, err := repository.pool.Exec(
		ctx,
		`UPDATE refresh_tokens SET revoked_at = $3 WHERE user_id = $1 AND device_id = $2 AND revoked_at IS NULL`,
		userID,
		deviceID,
		revokedAt,
	)
	return err
}
