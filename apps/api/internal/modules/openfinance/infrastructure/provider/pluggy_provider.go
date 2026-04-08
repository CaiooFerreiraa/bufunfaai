package provider

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"

	ofservice "github.com/bufunfaai/bufunfaai/apps/api/internal/modules/openfinance/application/service"
	"github.com/bufunfaai/bufunfaai/apps/api/internal/modules/openfinance/domain/entity"
	"github.com/bufunfaai/bufunfaai/apps/api/internal/platform/config"
)

const pluggyDefaultBaseURL string = "https://api.pluggy.ai"

type PluggyProvider struct {
	baseURL      string
	clientID     string
	clientSecret string
	staticAPIKey string
	httpClient   *http.Client

	mutex        sync.Mutex
	cachedAPIKey string
	cachedUntil  time.Time
}

type pluggyAuthResponse struct {
	APIKey      string `json:"apiKey"`
	AccessToken string `json:"accessToken"`
}

type pluggyConnectTokenResponse struct {
	AccessToken string `json:"accessToken"`
}

type pluggyConnectorListResponse struct {
	Results []pluggyConnector `json:"results"`
}

type pluggyConnector struct {
	ID             int64  `json:"id"`
	Name           string `json:"name"`
	InstitutionURL string `json:"institutionUrl"`
	ImageURL       string `json:"imageUrl"`
	Country        string `json:"country"`
	Type           string `json:"type"`
	IsOpenFinance  bool   `json:"isOpenFinance"`
}

type pluggyAccountListResponse struct {
	Results []pluggyAccount `json:"results"`
}

type pluggyAccount struct {
	ID            string            `json:"id"`
	ItemID        string            `json:"itemId"`
	Type          string            `json:"type"`
	Subtype       string            `json:"subtype"`
	Number        string            `json:"number"`
	Name          string            `json:"name"`
	MarketingName string            `json:"marketingName"`
	Balance       float64           `json:"balance"`
	CurrencyCode  string            `json:"currencyCode"`
	BankData      *pluggyBankData   `json:"bankData"`
	CreditData    *pluggyCreditData `json:"creditData"`
}

type pluggyBankData struct {
	TransferNumber string `json:"transferNumber"`
}

type pluggyCreditData struct {
	Brand                string  `json:"brand"`
	AvailableCreditLimit float64 `json:"availableCreditLimit"`
}

type pluggyTransactionListResponse struct {
	Total      int                 `json:"total"`
	TotalPages int                 `json:"totalPages"`
	Page       int                 `json:"page"`
	Results    []pluggyTransaction `json:"results"`
}

type pluggyTransaction struct {
	ID           string          `json:"id"`
	AccountID    string          `json:"accountId"`
	Description  string          `json:"description"`
	Amount       float64         `json:"amount"`
	Date         string          `json:"date"`
	CurrencyCode string          `json:"currencyCode"`
	Category     string          `json:"category"`
	Type         string          `json:"type"`
	Status       string          `json:"status"`
	Merchant     *pluggyMerchant `json:"merchant"`
}

type pluggyMerchant struct {
	Name string `json:"name"`
}

type pluggyItem struct {
	ID              string              `json:"id"`
	Status          string              `json:"status"`
	ExecutionStatus string              `json:"executionStatus"`
	ErrorCode       string              `json:"errorCode"`
	UpdatedAt       string              `json:"updatedAt"`
	Connector       pluggyItemConnector `json:"connector"`
}

type pluggyItemConnector struct {
	ID int64 `json:"id"`
}

func NewPluggyProvider(cfg config.Config) *PluggyProvider {
	baseURL := strings.TrimRight(strings.TrimSpace(cfg.OpenFinanceBaseURL), "/")
	if baseURL == "" {
		baseURL = pluggyDefaultBaseURL
	}

	return &PluggyProvider{
		baseURL:      baseURL,
		clientID:     strings.TrimSpace(cfg.OpenFinanceClientID),
		clientSecret: strings.TrimSpace(cfg.OpenFinanceSecret),
		staticAPIKey: strings.TrimSpace(cfg.OpenFinanceAPIKey),
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (provider *PluggyProvider) IsConfigured() bool {
	return (provider.clientID != "" && provider.clientSecret != "") || provider.staticAPIKey != ""
}

func (provider *PluggyProvider) ListInstitutions(ctx context.Context) ([]entity.Institution, error) {
	var response pluggyConnectorListResponse
	if err := provider.doJSON(ctx, http.MethodGet, "/connectors?pageSize=500", nil, &response); err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	institutions := make([]entity.Institution, 0, len(response.Results))
	for _, connector := range response.Results {
		if connector.Country != "BR" || !connector.IsOpenFinance {
			continue
		}

		institutions = append(institutions, entity.Institution{
			ID:                     uuid.NewSHA1(uuid.NameSpaceURL, []byte("pluggy:connector:"+strconv.FormatInt(connector.ID, 10))).String(),
			DirectoryOrgID:         strconv.FormatInt(connector.ID, 10),
			BrandName:              connector.Name,
			DisplayName:            connector.Name,
			AuthorisationServerURL: provider.baseURL,
			ResourcesBaseURL:       provider.baseURL,
			LogoURL:                connector.ImageURL,
			Status:                 "active",
			SupportsDataSharing:    true,
			SupportsPayments:       connector.Type == "PERSONAL_BANK" || connector.Type == "BUSINESS_BANK",
			LastDirectorySyncAt:    &now,
			CreatedAt:              now,
			UpdatedAt:              now,
		})
	}

	return institutions, nil
}

func (provider *PluggyProvider) CreateConsent(_ context.Context, _ entity.Institution, _ entity.Consent, _ []string) (string, *time.Time, error) {
	return "", nil, nil
}

func (provider *PluggyProvider) BuildAuthorizationURL(_ context.Context, _ entity.Institution, _ entity.Consent) (string, error) {
	return "", fmt.Errorf("authorization url not supported for pluggy provider")
}

func (provider *PluggyProvider) CreateConnectToken(ctx context.Context, institution entity.Institution, consent entity.Consent) (ofservice.ProviderConnectToken, error) {
	connectorID, err := strconv.ParseInt(strings.TrimSpace(institution.DirectoryOrgID), 10, 64)
	if err != nil {
		return ofservice.ProviderConnectToken{}, fmt.Errorf("invalid pluggy connector id: %w", err)
	}

	var response pluggyConnectTokenResponse
	body := map[string]any{
		"clientUserId":    consent.UserID,
		"avoidDuplicates": true,
	}
	if err := provider.doJSON(ctx, http.MethodPost, "/connect_token", body, &response); err != nil {
		return ofservice.ProviderConnectToken{}, err
	}

	return ofservice.ProviderConnectToken{
		ConnectToken:        response.AccessToken,
		SelectedConnectorID: connectorID,
	}, nil
}

func (provider *PluggyProvider) GetItem(ctx context.Context, itemID string) (ofservice.ProviderItem, error) {
	var response pluggyItem
	if err := provider.doJSON(ctx, http.MethodGet, "/items/"+itemID, nil, &response); err != nil {
		return ofservice.ProviderItem{}, err
	}

	return mapPluggyItem(response), nil
}

func (provider *PluggyProvider) ListAccounts(ctx context.Context, itemID string) ([]ofservice.ProviderAccount, error) {
	if strings.TrimSpace(itemID) == "" {
		return nil, fmt.Errorf("missing pluggy item id")
	}

	query := url.Values{}
	query.Set("itemId", itemID)
	query.Set("pageSize", "500")

	var response pluggyAccountListResponse
	if err := provider.doJSON(ctx, http.MethodGet, "/accounts?"+query.Encode(), nil, &response); err != nil {
		return nil, err
	}

	accounts := make([]ofservice.ProviderAccount, 0, len(response.Results))
	for _, account := range response.Results {
		mapped := ofservice.ProviderAccount{
			ID:            account.ID,
			ItemID:        account.ItemID,
			Type:          account.Type,
			Subtype:       account.Subtype,
			Number:        account.Number,
			Name:          account.Name,
			MarketingName: account.MarketingName,
			Balance:       account.Balance,
			CurrencyCode:  account.CurrencyCode,
		}

		if account.BankData != nil {
			mapped.BankTransferNumber = account.BankData.TransferNumber
		}
		if account.CreditData != nil {
			mapped.CreditBrand = account.CreditData.Brand
			mapped.AvailableCreditLimit = account.CreditData.AvailableCreditLimit
		}

		accounts = append(accounts, mapped)
	}

	return accounts, nil
}

func (provider *PluggyProvider) ListTransactions(ctx context.Context, accountID string, query ofservice.ProviderTransactionQuery) ([]ofservice.ProviderTransaction, error) {
	if strings.TrimSpace(accountID) == "" {
		return nil, fmt.Errorf("missing pluggy account id")
	}

	pageSize := query.PageSize
	if pageSize <= 0 || pageSize > 500 {
		pageSize = 500
	}

	page := 1
	transactions := make([]ofservice.ProviderTransaction, 0)
	for {
		values := url.Values{}
		values.Set("accountId", accountID)
		values.Set("page", strconv.Itoa(page))
		values.Set("pageSize", strconv.Itoa(pageSize))
		if query.From != nil {
			values.Set("from", query.From.UTC().Format("2006-01-02"))
		}
		if query.To != nil {
			values.Set("to", query.To.UTC().Format("2006-01-02"))
		}

		var response pluggyTransactionListResponse
		if err := provider.doJSON(ctx, http.MethodGet, "/transactions?"+values.Encode(), nil, &response); err != nil {
			return nil, err
		}

		for _, transaction := range response.Results {
			transactions = append(transactions, ofservice.ProviderTransaction{
				ID:           transaction.ID,
				AccountID:    transaction.AccountID,
				Description:  transaction.Description,
				Amount:       transaction.Amount,
				Date:         parsePluggyTimeValue(transaction.Date),
				CurrencyCode: transaction.CurrencyCode,
				Category:     transaction.Category,
				Type:         transaction.Type,
				Status:       transaction.Status,
				MerchantName: pluggyMerchantName(transaction.Merchant),
			})
		}

		if response.TotalPages <= page || len(response.Results) == 0 {
			break
		}
		page++
	}

	return transactions, nil
}

func (provider *PluggyProvider) ExchangeCode(_ context.Context, _ entity.Institution, _ entity.Consent, _ string) (ofservice.ProviderTokenSet, error) {
	return ofservice.ProviderTokenSet{}, fmt.Errorf("authorization code exchange not supported for pluggy provider")
}

func (provider *PluggyProvider) RevokeConsent(ctx context.Context, _ entity.Institution, consent entity.Consent) error {
	if strings.TrimSpace(consent.ExternalConsentID) == "" {
		return nil
	}

	return provider.doJSON(ctx, http.MethodDelete, "/items/"+consent.ExternalConsentID, nil, nil)
}

func (provider *PluggyProvider) SyncResources(ctx context.Context, _ entity.Institution, consent entity.Consent, _ entity.Connection) ([]ofservice.ProviderSyncResult, error) {
	itemID := strings.TrimSpace(consent.ExternalConsentID)
	if itemID == "" {
		return nil, fmt.Errorf("missing pluggy item id")
	}

	var refreshedItem pluggyItem
	err := provider.doJSON(ctx, http.MethodPatch, "/items/"+itemID, map[string]any{}, &refreshedItem)
	if err != nil {
		var currentItem pluggyItem
		if getErr := provider.doJSON(ctx, http.MethodGet, "/items/"+itemID, nil, &currentItem); getErr != nil {
			return nil, err
		}

		refreshedItem = currentItem
	}

	item := mapPluggyItem(refreshedItem)
	status := mapPluggyItemStatus(item)
	results := make([]ofservice.ProviderSyncResult, 0, 3)
	resourceTypes := []string{
		entity.ResourceAccounts,
		entity.ResourceBalances,
		entity.ResourceTransactions,
	}

	for _, resourceType := range resourceTypes {
		result := ofservice.ProviderSyncResult{
			ResourceType: resourceType,
			Status:       status,
		}

		if status == entity.SyncJobStatusError {
			result.ErrorCode = normalizePluggyErrorCode(item.ErrorCode)
			result.ErrorMessage = "Sua conexão com o banco ainda precisa de atenção."
		}

		results = append(results, result)
	}

	return results, nil
}

func (provider *PluggyProvider) doJSON(ctx context.Context, method string, path string, body any, output any) error {
	apiKey, err := provider.apiKey(ctx)
	if err != nil {
		return err
	}

	var payload []byte
	if body != nil {
		payload, err = json.Marshal(body)
		if err != nil {
			return err
		}
	}

	request, err := http.NewRequestWithContext(ctx, method, provider.baseURL+path, bytes.NewReader(payload))
	if err != nil {
		return err
	}

	request.Header.Set("X-API-KEY", apiKey)
	if body != nil {
		request.Header.Set("Content-Type", "application/json")
	}

	response, err := provider.httpClient.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode >= http.StatusBadRequest {
		return fmt.Errorf("pluggy request failed with status %d", response.StatusCode)
	}

	if output == nil {
		return nil
	}

	return json.NewDecoder(response.Body).Decode(output)
}

func (provider *PluggyProvider) apiKey(ctx context.Context) (string, error) {
	if provider.clientID == "" || provider.clientSecret == "" {
		if provider.staticAPIKey == "" {
			return "", fmt.Errorf("pluggy credentials not configured")
		}

		return provider.staticAPIKey, nil
	}

	provider.mutex.Lock()
	defer provider.mutex.Unlock()

	if provider.cachedAPIKey != "" && time.Now().UTC().Before(provider.cachedUntil) {
		return provider.cachedAPIKey, nil
	}

	payload, err := json.Marshal(map[string]string{
		"clientId":     provider.clientID,
		"clientSecret": provider.clientSecret,
	})
	if err != nil {
		if provider.staticAPIKey != "" {
			return provider.staticAPIKey, nil
		}

		return "", err
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, provider.baseURL+"/auth", bytes.NewReader(payload))
	if err != nil {
		if provider.staticAPIKey != "" {
			return provider.staticAPIKey, nil
		}

		return "", err
	}
	request.Header.Set("Content-Type", "application/json")

	response, err := provider.httpClient.Do(request)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	if response.StatusCode >= http.StatusBadRequest {
		if provider.staticAPIKey != "" {
			return provider.staticAPIKey, nil
		}

		return "", fmt.Errorf("pluggy auth failed with status %d", response.StatusCode)
	}

	var authResponse pluggyAuthResponse
	if err := json.NewDecoder(response.Body).Decode(&authResponse); err != nil {
		if provider.staticAPIKey != "" {
			return provider.staticAPIKey, nil
		}

		return "", err
	}

	provider.cachedAPIKey = strings.TrimSpace(authResponse.APIKey)
	if provider.cachedAPIKey == "" {
		provider.cachedAPIKey = strings.TrimSpace(authResponse.AccessToken)
	}
	if provider.cachedAPIKey == "" {
		if provider.staticAPIKey != "" {
			return provider.staticAPIKey, nil
		}

		return "", fmt.Errorf("pluggy auth response missing api key")
	}

	provider.cachedUntil = time.Now().UTC().Add(110 * time.Minute)
	return provider.cachedAPIKey, nil
}

func mapPluggyItem(item pluggyItem) ofservice.ProviderItem {
	return ofservice.ProviderItem{
		ID:              item.ID,
		ConnectorID:     item.Connector.ID,
		Status:          item.Status,
		ExecutionStatus: item.ExecutionStatus,
		LastUpdatedAt:   parsePluggyTime(item.UpdatedAt),
		ErrorCode:       item.ErrorCode,
	}
}

func pluggyMerchantName(merchant *pluggyMerchant) string {
	if merchant == nil {
		return ""
	}

	return strings.TrimSpace(merchant.Name)
}

func parsePluggyTime(raw string) *time.Time {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil
	}

	parsed, err := time.Parse(time.RFC3339, raw)
	if err != nil {
		return nil
	}

	return &parsed
}

func parsePluggyTimeValue(raw string) time.Time {
	parsed := parsePluggyTime(raw)
	if parsed == nil {
		return time.Time{}
	}

	return parsed.UTC()
}

func mapPluggyItemStatus(item ofservice.ProviderItem) string {
	switch item.Status {
	case "UPDATED":
		return entity.SyncJobStatusCompleted
	case "WAITING_USER_INPUT", "WAITING_USER_ACTION", "LOGIN_ERROR":
		return entity.SyncJobStatusError
	default:
		if strings.Contains(item.ExecutionStatus, "ERROR") {
			return entity.SyncJobStatusError
		}

		if strings.Contains(item.ExecutionStatus, "IN_PROGRESS") || item.Status == "UPDATING" || item.Status == "CREATED" {
			return entity.SyncJobStatusPending
		}

		return entity.SyncJobStatusCompleted
	}
}

func normalizePluggyErrorCode(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return "PLUGGY_ITEM_ERROR"
	}

	return value
}
