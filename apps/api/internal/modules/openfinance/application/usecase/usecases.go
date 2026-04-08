package usecase

import (
	"context"

	ofdto "github.com/bufunfaai/bufunfaai/apps/api/internal/modules/openfinance/application/dto"
	ofservice "github.com/bufunfaai/bufunfaai/apps/api/internal/modules/openfinance/application/service"
	"github.com/bufunfaai/bufunfaai/apps/api/internal/modules/openfinance/domain/entity"
	sharederrors "github.com/bufunfaai/bufunfaai/apps/api/internal/shared/errors"
)

type UseCases struct {
	service *ofservice.Service
}

func New(service *ofservice.Service) *UseCases {
	return &UseCases{service: service}
}

func (useCases *UseCases) ListInstitutions(ctx context.Context) ([]entity.Institution, *sharederrors.AppError) {
	return useCases.service.ListInstitutions(ctx)
}

func (useCases *UseCases) GetInstitution(ctx context.Context, institutionID string) (entity.Institution, *sharederrors.AppError) {
	return useCases.service.GetInstitution(ctx, institutionID)
}

func (useCases *UseCases) CreateConsent(ctx context.Context, userID string, request ofdto.CreateConsentRequest) (entity.Consent, *sharederrors.AppError) {
	return useCases.service.CreateConsent(ctx, userID, request)
}

func (useCases *UseCases) GetConsent(ctx context.Context, consentID string, userID string) (entity.Consent, *sharederrors.AppError) {
	return useCases.service.GetConsent(ctx, consentID, userID)
}

func (useCases *UseCases) AuthorizeConsent(ctx context.Context, consentID string, userID string) (string, *sharederrors.AppError) {
	return useCases.service.AuthorizeConsent(ctx, consentID, userID)
}

func (useCases *UseCases) CreateConnectToken(ctx context.Context, consentID string, userID string) (ofservice.ProviderConnectToken, *sharederrors.AppError) {
	return useCases.service.CreateConnectToken(ctx, consentID, userID)
}

func (useCases *UseCases) CompleteConsent(ctx context.Context, consentID string, userID string, itemID string) (entity.Consent, entity.Connection, *sharederrors.AppError) {
	return useCases.service.CompleteConsent(ctx, consentID, userID, itemID)
}

func (useCases *UseCases) HandleCallback(ctx context.Context, state string, code string) (entity.Consent, entity.Connection, *sharederrors.AppError) {
	return useCases.service.HandleCallback(ctx, state, code)
}

func (useCases *UseCases) RevokeConsent(ctx context.Context, consentID string, userID string) *sharederrors.AppError {
	return useCases.service.RevokeConsent(ctx, consentID, userID)
}

func (useCases *UseCases) ListConnections(ctx context.Context, userID string) ([]entity.Connection, *sharederrors.AppError) {
	return useCases.service.ListConnections(ctx, userID)
}

func (useCases *UseCases) GetConnection(ctx context.Context, connectionID string, userID string) (entity.Connection, *sharederrors.AppError) {
	return useCases.service.GetConnection(ctx, connectionID, userID)
}

func (useCases *UseCases) SyncConnection(ctx context.Context, connectionID string, userID string) ([]entity.SyncJob, *sharederrors.AppError) {
	return useCases.service.SyncConnection(ctx, connectionID, userID)
}

func (useCases *UseCases) SyncStatus(ctx context.Context, connectionID string, userID string) (entity.Connection, []entity.SyncJob, *sharederrors.AppError) {
	return useCases.service.SyncStatus(ctx, connectionID, userID)
}

func (useCases *UseCases) EnsureInstitutions(ctx context.Context) *sharederrors.AppError {
	return useCases.service.EnsureInstitutions(ctx)
}

func (useCases *UseCases) ReconcileConnections(ctx context.Context, limit int) (ofservice.ReconciliationResult, *sharederrors.AppError) {
	return useCases.service.ReconcileConnections(ctx, limit)
}
