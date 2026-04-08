package repository

import (
	"context"

	"github.com/bufunfaai/bufunfaai/apps/api/internal/modules/openfinance/domain/entity"
)

type InstitutionRepository interface {
	SaveMany(ctx context.Context, institutions []entity.Institution) error
	List(ctx context.Context) ([]entity.Institution, error)
	GetByID(ctx context.Context, institutionID string) (entity.Institution, error)
}

type ConsentRepository interface {
	Create(ctx context.Context, consent entity.Consent) error
	GetByID(ctx context.Context, consentID string) (entity.Consent, error)
	GetByState(ctx context.Context, state string) (entity.Consent, error)
	Update(ctx context.Context, consent entity.Consent) error
}

type AuthorizationRepository interface {
	Create(ctx context.Context, authorization entity.Authorization) error
}

type TokenRepository interface {
	UpsertByConsentID(ctx context.Context, token entity.Token) error
	GetByConsentID(ctx context.Context, consentID string) (entity.Token, error)
	RevokeByConsentID(ctx context.Context, consentID string) error
}

type ConnectionRepository interface {
	CreateOrUpdate(ctx context.Context, connection entity.Connection) (entity.Connection, error)
	ListByUserID(ctx context.Context, userID string) ([]entity.Connection, error)
	ListActive(ctx context.Context, limit int) ([]entity.Connection, error)
	GetByID(ctx context.Context, connectionID string) (entity.Connection, error)
	GetByConsentID(ctx context.Context, consentID string) (entity.Connection, error)
	Update(ctx context.Context, connection entity.Connection) error
}

type SyncJobRepository interface {
	ReplaceForConnection(ctx context.Context, connectionID string, jobs []entity.SyncJob) error
	ListByConnectionID(ctx context.Context, connectionID string) ([]entity.SyncJob, error)
}
