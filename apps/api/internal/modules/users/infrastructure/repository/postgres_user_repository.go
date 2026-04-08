package repository

import (
	"context"
	"errors"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/bufunfaai/bufunfaai/apps/api/internal/modules/users/application/service"
	"github.com/bufunfaai/bufunfaai/apps/api/internal/modules/users/domain/entity"
)

type PostgresUserRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresUserRepository(pool *pgxpool.Pool) *PostgresUserRepository {
	return &PostgresUserRepository{pool: pool}
}

func (repository *PostgresUserRepository) Create(ctx context.Context, user entity.User) (entity.User, error) {
	query := `
		INSERT INTO users (
			id, full_name, email, phone, password_hash, status, email_verified_at, last_login_at, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10
		)
		RETURNING id, full_name, email, phone, password_hash, status, email_verified_at, last_login_at, created_at, updated_at
	`

	row := repository.pool.QueryRow(
		ctx,
		query,
		user.ID,
		user.FullName,
		strings.ToLower(strings.TrimSpace(user.Email)),
		user.Phone,
		user.PasswordHash,
		user.Status,
		user.EmailVerifiedAt,
		user.LastLoginAt,
		user.CreatedAt,
		user.UpdatedAt,
	)

	return scanUser(row)
}

func (repository *PostgresUserRepository) GetByEmail(ctx context.Context, email string) (entity.User, error) {
	query := `
		SELECT id, full_name, email, phone, password_hash, status, email_verified_at, last_login_at, created_at, updated_at
		FROM users
		WHERE email = $1 AND deleted_at IS NULL
	`

	row := repository.pool.QueryRow(ctx, query, strings.ToLower(strings.TrimSpace(email)))
	return scanUser(row)
}

func (repository *PostgresUserRepository) GetByID(ctx context.Context, userID string) (entity.User, error) {
	query := `
		SELECT id, full_name, email, phone, password_hash, status, email_verified_at, last_login_at, created_at, updated_at
		FROM users
		WHERE id = $1 AND deleted_at IS NULL
	`

	row := repository.pool.QueryRow(ctx, query, userID)
	return scanUser(row)
}

func (repository *PostgresUserRepository) UpdateProfile(ctx context.Context, userID string, fullName string, phone string) (entity.User, error) {
	query := `
		UPDATE users
		SET full_name = $2, phone = $3, updated_at = NOW()
		WHERE id = $1 AND deleted_at IS NULL
		RETURNING id, full_name, email, phone, password_hash, status, email_verified_at, last_login_at, created_at, updated_at
	`

	row := repository.pool.QueryRow(ctx, query, userID, fullName, phone)
	return scanUser(row)
}

func (repository *PostgresUserRepository) UpdateLastLogin(ctx context.Context, userID string) error {
	commandTag, err := repository.pool.Exec(
		ctx,
		`UPDATE users SET last_login_at = NOW(), updated_at = NOW() WHERE id = $1 AND deleted_at IS NULL`,
		userID,
	)
	if err != nil {
		return err
	}

	if commandTag.RowsAffected() == 0 {
		return service.ErrUserNotFound
	}

	return nil
}

type scanner interface {
	Scan(dest ...any) error
}

func scanUser(row scanner) (entity.User, error) {
	var user entity.User
	err := row.Scan(
		&user.ID,
		&user.FullName,
		&user.Email,
		&user.Phone,
		&user.PasswordHash,
		&user.Status,
		&user.EmailVerifiedAt,
		&user.LastLoginAt,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.User{}, service.ErrUserNotFound
		}

		return entity.User{}, err
	}

	return user, nil
}
