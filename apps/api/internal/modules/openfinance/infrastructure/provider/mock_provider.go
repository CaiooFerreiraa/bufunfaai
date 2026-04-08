package provider

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/google/uuid"

	ofservice "github.com/bufunfaai/bufunfaai/apps/api/internal/modules/openfinance/application/service"
	"github.com/bufunfaai/bufunfaai/apps/api/internal/modules/openfinance/domain/entity"
)

type MockProvider struct{}

func NewMockProvider() *MockProvider {
	return &MockProvider{}
}

func (provider *MockProvider) ListInstitutions(_ context.Context) ([]entity.Institution, error) {
	now := time.Now().UTC()
	return []entity.Institution{
		{
			ID:                     "7a5acb89-2f24-49c5-b5e7-8fbd62af8f00",
			DirectoryOrgID:         "1",
			BrandName:              "Mock Bank",
			DisplayName:            "Mock Bank Sandbox",
			AuthorisationServerURL: "https://mock-bank.example/auth",
			ResourcesBaseURL:       "https://mock-bank.example/resources",
			LogoURL:                "https://mock-bank.example/logo.png",
			Status:                 "active",
			SupportsDataSharing:    true,
			SupportsPayments:       false,
			LastDirectorySyncAt:    &now,
			CreatedAt:              now,
			UpdatedAt:              now,
		},
	}, nil
}

func (provider *MockProvider) CreateConsent(_ context.Context, _ entity.Institution, _ entity.Consent, _ []string) (string, *time.Time, error) {
	expiresAt := time.Now().UTC().Add(90 * 24 * time.Hour)
	return "mock-consent-" + uuid.NewString(), &expiresAt, nil
}

func (provider *MockProvider) BuildAuthorizationURL(_ context.Context, _ entity.Institution, consent entity.Consent) (string, error) {
	callbackURL, err := url.Parse(consent.RedirectURI)
	if err != nil {
		return "", err
	}

	query := callbackURL.Query()
	query.Set("state", consent.State)
	query.Set("code", "mock-code-"+consent.ID)
	callbackURL.RawQuery = query.Encode()

	return callbackURL.String(), nil
}

func (provider *MockProvider) CreateConnectToken(_ context.Context, _ entity.Institution, consent entity.Consent) (ofservice.ProviderConnectToken, error) {
	return ofservice.ProviderConnectToken{
		ConnectToken:        "mock-connect-token-" + consent.ID,
		SelectedConnectorID: 1,
	}, nil
}

func (provider *MockProvider) GetItem(_ context.Context, itemID string) (ofservice.ProviderItem, error) {
	now := time.Now().UTC()
	return ofservice.ProviderItem{
		ID:            itemID,
		ConnectorID:   1,
		Status:        "UPDATED",
		LastUpdatedAt: &now,
	}, nil
}

func (provider *MockProvider) ExchangeCode(_ context.Context, _ entity.Institution, consent entity.Consent, code string) (ofservice.ProviderTokenSet, error) {
	if !strings.HasPrefix(code, "mock-code-"+consent.ID) {
		return ofservice.ProviderTokenSet{}, fmt.Errorf("invalid mock authorization code")
	}

	expiresAt := time.Now().UTC().Add(30 * time.Minute)
	refreshExpiresAt := time.Now().UTC().Add(180 * 24 * time.Hour)
	return ofservice.ProviderTokenSet{
		AccessToken:      "mock-access-token-" + consent.ID,
		RefreshToken:     "mock-refresh-token-" + consent.ID,
		TokenType:        "Bearer",
		Scope:            "openid accounts balances transactions",
		ExpiresAt:        expiresAt,
		RefreshExpiresAt: &refreshExpiresAt,
	}, nil
}

func (provider *MockProvider) RevokeConsent(_ context.Context, _ entity.Institution, _ entity.Consent) error {
	return nil
}

func (provider *MockProvider) SyncResources(_ context.Context, _ entity.Institution, _ entity.Consent, _ entity.Connection) ([]ofservice.ProviderSyncResult, error) {
	return []ofservice.ProviderSyncResult{
		{ResourceType: entity.ResourceAccounts, Status: entity.SyncJobStatusCompleted},
		{ResourceType: entity.ResourceBalances, Status: entity.SyncJobStatusCompleted},
		{ResourceType: entity.ResourceTransactions, Status: entity.SyncJobStatusCompleted},
	}, nil
}
