package presenter

import (
	"encoding/json"
	"time"

	ofdto "github.com/bufunfaai/bufunfaai/apps/api/internal/modules/openfinance/application/dto"
	"github.com/bufunfaai/bufunfaai/apps/api/internal/modules/openfinance/domain/entity"
)

func InstitutionOutput(institution entity.Institution) ofdto.InstitutionOutput {
	return ofdto.InstitutionOutput{
		ID:                     institution.ID,
		DirectoryOrgID:         institution.DirectoryOrgID,
		BrandName:              institution.BrandName,
		DisplayName:            institution.DisplayName,
		AuthorisationServerURL: institution.AuthorisationServerURL,
		ResourcesBaseURL:       institution.ResourcesBaseURL,
		LogoURL:                institution.LogoURL,
		Status:                 institution.Status,
		SupportsDataSharing:    institution.SupportsDataSharing,
		SupportsPayments:       institution.SupportsPayments,
	}
}

func ConsentOutput(consent entity.Consent) ofdto.ConsentOutput {
	return ofdto.ConsentOutput{
		ID:                consent.ID,
		UserID:            consent.UserID,
		InstitutionID:     consent.InstitutionID,
		ExternalConsentID: consent.ExternalConsentID,
		Status:            consent.Status,
		Purpose:           consent.Purpose,
		Permissions:       parsePermissions(consent.PermissionsJSON),
		ExpirationAt:      formatTimePointer(consent.ExpirationAt),
		RevokedAt:         formatTimePointer(consent.RevokedAt),
		AuthorisedAt:      formatTimePointer(consent.AuthorisedAt),
		RedirectURI:       consent.RedirectURI,
		CreatedAt:         consent.CreatedAt.Format(time.RFC3339),
		UpdatedAt:         consent.UpdatedAt.Format(time.RFC3339),
	}
}

func ConnectionOutput(connection entity.Connection) ofdto.ConnectionOutput {
	return ofdto.ConnectionOutput{
		ID:                   connection.ID,
		UserID:               connection.UserID,
		InstitutionID:        connection.InstitutionID,
		ConsentID:            connection.ConsentID,
		Status:               connection.Status,
		FirstSyncAt:          formatTimePointer(connection.FirstSyncAt),
		LastSyncAt:           formatTimePointer(connection.LastSyncAt),
		LastSuccessfulSyncAt: formatTimePointer(connection.LastSuccessfulSyncAt),
		LastErrorCode:        connection.LastErrorCode,
		LastErrorMessage:     connection.LastErrorMessageRedacted,
		CreatedAt:            connection.CreatedAt.Format(time.RFC3339),
		UpdatedAt:            connection.UpdatedAt.Format(time.RFC3339),
	}
}

func SyncJobOutput(job entity.SyncJob) ofdto.SyncJobOutput {
	return ofdto.SyncJobOutput{
		ID:                   job.ID,
		ConnectionID:         job.ConnectionID,
		ResourceType:         job.ResourceType,
		Status:               job.Status,
		AttemptCount:         job.AttemptCount,
		ScheduledAt:          formatTimePointer(job.ScheduledAt),
		StartedAt:            formatTimePointer(job.StartedAt),
		FinishedAt:           formatTimePointer(job.FinishedAt),
		ErrorCode:            job.ErrorCode,
		ErrorMessageRedacted: job.ErrorMessageRedacted,
	}
}

func parsePermissions(rawPermissions string) []string {
	if rawPermissions == "" {
		return []string{}
	}

	permissions := make([]string, 0)
	if err := json.Unmarshal([]byte(rawPermissions), &permissions); err != nil {
		return []string{}
	}

	return permissions
}

func formatTimePointer(value *time.Time) string {
	if value == nil {
		return ""
	}

	return value.Format(time.RFC3339)
}
