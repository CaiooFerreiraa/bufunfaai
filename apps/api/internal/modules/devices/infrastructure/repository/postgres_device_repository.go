package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/bufunfaai/bufunfaai/apps/api/internal/modules/devices/application/service"
	"github.com/bufunfaai/bufunfaai/apps/api/internal/modules/devices/domain/entity"
)

type PostgresDeviceRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresDeviceRepository(pool *pgxpool.Pool) *PostgresDeviceRepository {
	return &PostgresDeviceRepository{pool: pool}
}

func (repository *PostgresDeviceRepository) Upsert(ctx context.Context, device entity.Device) (entity.Device, error) {
	if device.FingerprintHash == "" {
		query := `
			INSERT INTO devices (
				id, user_id, device_name, platform, app_version, device_fingerprint_hash, last_seen_at, created_at
			) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
			RETURNING id, user_id, device_name, platform, app_version, device_fingerprint_hash, last_seen_at, created_at
		`

		row := repository.pool.QueryRow(
			ctx,
			query,
			device.ID,
			device.UserID,
			device.DeviceName,
			device.Platform,
			device.AppVersion,
			device.FingerprintHash,
			device.LastSeenAt,
			device.CreatedAt,
		)

		return scanDevice(row)
	}

	query := `
		INSERT INTO devices (
			id, user_id, device_name, platform, app_version, device_fingerprint_hash, last_seen_at, created_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		ON CONFLICT (user_id, device_fingerprint_hash)
		DO UPDATE SET
			device_name = EXCLUDED.device_name,
			platform = EXCLUDED.platform,
			app_version = EXCLUDED.app_version,
			last_seen_at = EXCLUDED.last_seen_at
		RETURNING id, user_id, device_name, platform, app_version, device_fingerprint_hash, last_seen_at, created_at
	`

	row := repository.pool.QueryRow(
		ctx,
		query,
		device.ID,
		device.UserID,
		device.DeviceName,
		device.Platform,
		device.AppVersion,
		device.FingerprintHash,
		device.LastSeenAt,
		device.CreatedAt,
	)

	return scanDevice(row)
}

func (repository *PostgresDeviceRepository) ListByUserID(ctx context.Context, userID string) ([]entity.Device, error) {
	rows, err := repository.pool.Query(
		ctx,
		`SELECT id, user_id, device_name, platform, app_version, device_fingerprint_hash, last_seen_at, created_at FROM devices WHERE user_id = $1 ORDER BY last_seen_at DESC`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	devices := make([]entity.Device, 0)
	for rows.Next() {
		device, err := scanDevice(rows)
		if err != nil {
			return nil, err
		}

		devices = append(devices, device)
	}

	return devices, rows.Err()
}

func (repository *PostgresDeviceRepository) Delete(ctx context.Context, userID string, deviceID string) error {
	commandTag, err := repository.pool.Exec(
		ctx,
		`DELETE FROM devices WHERE id = $1 AND user_id = $2`,
		deviceID,
		userID,
	)
	if err != nil {
		return err
	}

	if commandTag.RowsAffected() == 0 {
		return service.ErrDeviceNotFound
	}

	return nil
}

type deviceScanner interface {
	Scan(dest ...any) error
}

func scanDevice(row deviceScanner) (entity.Device, error) {
	var device entity.Device
	err := row.Scan(
		&device.ID,
		&device.UserID,
		&device.DeviceName,
		&device.Platform,
		&device.AppVersion,
		&device.FingerprintHash,
		&device.LastSeenAt,
		&device.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.Device{}, service.ErrDeviceNotFound
		}

		return entity.Device{}, err
	}

	return device, nil
}
